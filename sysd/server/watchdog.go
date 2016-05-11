package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"infra/sysd/sysdCommonDefs"
	"models"
	"os"
	"os/exec"
	"strings"
	"sysd"
	"time"
)

const (
	KA_TIMEOUT_COUNT_MIN = 0
	KA_TIMEOUT_COUNT     = 5 // After 5 KA missed from a daemon, sysd will assume the daemon as non-responsive. Restart it.
	WD_MAX_NUM_RESTARTS  = 5 // After 5 restarts, if this daemon is still not responsive then stop it.
)

const (
	REASON_NONE           = "None"
	REASON_KA_FAIL        = "Failed to receive keepalive messages"
	REASON_USER_RESTART   = "Restarted by user"
	REASON_DAEMON_STOPPED = "Stopped by user"
	REASON_COMING_UP      = "Started by user"
)

type DaemonInfo struct {
	Enable        bool
	State         sysdCommonDefs.SRDaemonStatus
	Reason        string
	RecvedKACount int32
	NumRestarts   int32
	RestartTime   string
	RestartReason string
}

func (daemonInfo *DaemonInfo) Initialize() error {
	daemonInfo.Enable = true
	daemonInfo.State = sysdCommonDefs.STARTING
	daemonInfo.Reason = REASON_COMING_UP
	daemonInfo.RecvedKACount = 0
	daemonInfo.NumRestarts = 0
	daemonInfo.RestartTime = ""
	daemonInfo.RestartReason = ""
	return nil
}

func (server *SYSDServer) StartWDRoutine() error {
	server.DaemonMap = make(map[string]*DaemonInfo)
	server.KaRecvCh = make(chan string, sysdCommonDefs.SYSD_TOTAL_KA_DAEMONS)
	server.DaemonConfigCh = make(chan DaemonConfig, sysdCommonDefs.SYSD_TOTAL_KA_DAEMONS)
	server.UpdateInfoInDbCh = make(chan string, sysdCommonDefs.SYSD_TOTAL_KA_DAEMONS)
	server.DaemonRestartCh = make(chan string, sysdCommonDefs.SYSD_TOTAL_KA_DAEMONS)
	go server.WDTimer()
	for {
		select {
		case kaDaemon := <-server.KaRecvCh:
			daemonInfo, exist := server.DaemonMap[kaDaemon]
			if !exist {
				daemonInfo = &DaemonInfo{}
				daemonInfo.Initialize()
				server.DaemonMap[kaDaemon] = daemonInfo
			}
			daemonInfo.RecvedKACount++
			if daemonInfo.State != sysdCommonDefs.UP {
				daemonInfo.State = sysdCommonDefs.UP
				daemonInfo.Reason = REASON_NONE
				server.PublishDaemonKANotification(kaDaemon, daemonInfo.State)
			}
			server.UpdateDaemonStateInDb(kaDaemon)
		case daemonConfig := <-server.DaemonConfigCh:
			server.logger.Info(fmt.Sprintln("Received daemon config for: ", daemonConfig.Name, " Enable ", daemonConfig.Enable))
			daemon := daemonConfig.Name
			enable := daemonConfig.Enable
			daemonInfo, exist := server.DaemonMap[daemon]
			daemonUpdated := false
			if enable {
				if !exist {
					daemonInfo = &DaemonInfo{}
					daemonInfo.Initialize()
					server.DaemonMap[daemon] = daemonInfo
				}
				if daemonInfo.State == sysdCommonDefs.STOPPED {
					go server.ToggleFlexswitchDaemon(daemon, true)
					daemonInfo.State = sysdCommonDefs.STARTING
					daemonInfo.Reason = REASON_COMING_UP
				}
				daemonInfo.Enable = true
				daemonUpdated = true
			} else {
				if exist {
					if daemonInfo.State != sysdCommonDefs.STOPPED {
						go server.ToggleFlexswitchDaemon(daemon, false)
						daemonInfo.State = sysdCommonDefs.STOPPED
						daemonInfo.Reason = REASON_DAEMON_STOPPED
						server.PublishDaemonKANotification(daemon, daemonInfo.State)
					}
					daemonInfo.Enable = false
					daemonUpdated = true
				} else {
					server.logger.Info(fmt.Sprintln("Received call to stop unknown daemon ", daemon))
				}
			}
			if daemonUpdated {
				server.UpdateDaemonStateInDb(daemon)
			}
		case daemon := <-server.DaemonRestartCh:
			if server.DaemonMap[daemon].State != sysdCommonDefs.STOPPED {
				go server.RestartFlexswitchDaemon(daemon)
			}
		case daemon := <-server.UpdateInfoInDbCh:
			server.UpdateDaemonStateInDb(daemon)
		}
	}
	return nil
}

func (server *SYSDServer) PublishDaemonKANotification(name string, status sysdCommonDefs.SRDaemonStatus) error {
	msg := sysdCommonDefs.DaemonStatus{
		Name:   name,
		Status: status,
	}
	msgBuf, err := json.Marshal(msg)
	if err != nil {
		server.logger.Err("Failed to marshal daemon status")
		return err
	}
	notification := sysdCommonDefs.Notification{
		Type:    uint8(sysdCommonDefs.KA_DAEMON),
		Payload: msgBuf,
	}
	notificationBuf, err := json.Marshal(notification)
	if err != nil {
		server.logger.Err("Failed to marshal daemon status message")
		return err
	}
	server.notificationCh <- notificationBuf
	return nil
}

func (server *SYSDServer) ToggleFlexswitchDaemon(daemon string, up bool) error {
	var (
		cmdOut []byte
		err    error
		op     string
	)
	cmdDir := strings.TrimSuffix(server.paramsDir, "params/")
	cmdName := cmdDir + "flexswitch"
	if up {
		op = "start"
	} else {
		op = "stop"
	}
	//cmdArgs := []string{"-n", daemon, "-o", op, "-d", cmdDir}
	cmdArgs := []string{"-n", daemon, "-o", op}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		server.logger.Info(fmt.Sprintln(os.Stderr, "There was an error to ", op, " flexswitch daemon ", daemon, " : ", err))
		return err
	}
	out := string(cmdOut)
	server.logger.Info(fmt.Sprintln("Flexswitch daemon ", daemon, op, " returned ", out))

	return nil
}

func (server *SYSDServer) RestartFlexswitchDaemon(daemon string) error {
	server.PublishDaemonKANotification(daemon, sysdCommonDefs.RESTARTING)
	server.ToggleFlexswitchDaemon(daemon, false)
	server.ToggleFlexswitchDaemon(daemon, true)
	return nil
}

func (server *SYSDServer) WDTimer() error {
	server.logger.Info("Starting system WD")
	wdTimer := time.NewTicker(time.Second * KA_TIMEOUT_COUNT)
	for t := range wdTimer.C {
		_ = t
		for daemon, daemonInfo := range server.DaemonMap {
			if daemonInfo.State == sysdCommonDefs.UP {
				if daemonInfo.RecvedKACount < KA_TIMEOUT_COUNT && daemonInfo.RecvedKACount > KA_TIMEOUT_COUNT_MIN {
					server.logger.Info(fmt.Sprintln("Daemon ", daemon, " is slowing down. Monitoring it."))
				}
				if daemonInfo.RecvedKACount == KA_TIMEOUT_COUNT_MIN {
					server.logger.Info(fmt.Sprintln("Daemon ", daemon, " is not responsive. Restarting it."))
					server.DaemonRestartCh <- daemon
					daemonInfo.State = sysdCommonDefs.RESTARTING
					daemonInfo.NumRestarts++
					daemonInfo.RestartTime = time.Now().String()
					daemonInfo.RestartReason = REASON_KA_FAIL
					server.UpdateInfoInDbCh <- daemon
				}
			}
			daemonInfo.RecvedKACount = 0
		}
	}
	return nil
}

func (server *SYSDServer) ConvertDaemonStateToThrift(ent DaemonState) *sysd.DaemonState {
	dState := sysd.NewDaemonState()
	dState.Name = string(ent.Name)
	dState.Enable = ent.Enable
	dState.State = string(sysdCommonDefs.ConvertDaemonStateCodeToString(ent.State))
	dState.Reason = string(ent.Reason)
	kaStr := fmt.Sprintf("Received %d keepalives", ent.RecvedKACount)
	dState.KeepAlive = string(kaStr)
	dState.RestartCount = int32(ent.NumRestarts)
	dState.RestartTime = string(ent.RestartTime)
	dState.RestartReason = string(ent.RestartReason)
	return dState
}

func (server *SYSDServer) ConvertDaemonStateToObj(ent DaemonState) models.DaemonState {
	kaStr := fmt.Sprintf("Received %d keepalives", ent.RecvedKACount)
	dState := models.DaemonState{
		Name:          ent.Name,
		Enable:        ent.Enable,
		State:         sysdCommonDefs.ConvertDaemonStateCodeToString(ent.State),
		Reason:        ent.Reason,
		KeepAlive:     kaStr,
		RestartCount:  ent.NumRestarts,
		RestartTime:   ent.RestartTime,
		RestartReason: ent.RestartReason,
	}
	return dState
}

func (server *SYSDServer) GetDaemonState(name string) *DaemonState {
	daemonState := new(DaemonState)
	daemonInfo, found := server.DaemonMap[name]
	if found {
		daemonState.Name = name
		daemonState.Enable = daemonInfo.Enable
		daemonState.State = daemonInfo.State
		daemonState.Reason = daemonInfo.Reason
		daemonState.RecvedKACount = daemonInfo.RecvedKACount
		daemonState.NumRestarts = daemonInfo.NumRestarts
		daemonState.RestartTime = daemonInfo.RestartTime
		daemonState.RestartReason = daemonInfo.RestartReason
	}
	return daemonState
}

func (server *SYSDServer) GetBulkDaemonStates(idx int, cnt int) (int, int, []DaemonState) {
	var nextIdx int
	var count int
	server.logger.Info(fmt.Sprintln("GetBulk DaemonStates"))
	result := make([]DaemonState, cnt)
	i := 0
	for daemon, daemonInfo := range server.DaemonMap {
		result[i].Name = daemon
		result[i].Enable = daemonInfo.Enable
		result[i].State = daemonInfo.State
		result[i].Reason = daemonInfo.Reason
		result[i].RecvedKACount = daemonInfo.RecvedKACount
		result[i].NumRestarts = daemonInfo.NumRestarts
		result[i].RestartTime = daemonInfo.RestartTime
		result[i].RestartReason = daemonInfo.RestartReason
		i++
	}
	count = i
	nextIdx = 0
	return nextIdx, count, result
}

func (server *SYSDServer) UpdateDaemonStateInDb(name string) error {
	var err error
	daemonState := server.GetDaemonState(name)
	if daemonState != nil {
		obj := server.ConvertDaemonStateToObj(*daemonState)
		server.dbHdl.StoreObjectInDb(obj)
	} else {
		errStr := "Failed to get daemon " + name
		server.logger.Info(errStr)
		err = errors.New(errStr)
	}
	return err
}
