package server

import (
	"encoding/json"
	"fmt"
	"infra/sysd/sysdCommonDefs"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	KA_TIMEOUT_COUNT_MIN  = 0
	KA_TIMEOUT_COUNT      = 5 // After 5 KA missed from a daemon, sysd will assume the daemon as non-responsive. Restart it.
	WD_MAX_NUM_RESTARTS   = 5 // After 5 restarts, if this daemon is still not responsive then stop it.
	SYSD_TOTAL_KA_DAEMONS = 32
)

type WDInfo struct {
	Active        bool
	RecvedKACount int32
	NumRestarts   int32
}

func (server *SYSDServer) StartWDRoutine() error {
	server.KaRecvCh = make(chan string, sysdCommonDefs.SYSD_TOTAL_KA_DAEMONS)
	server.KaRecvMap = make(map[string]*WDInfo)
	go server.WDTimer()
	for {
		select {
		case kaDaemon := <-server.KaRecvCh:
			if server.KaRecvMap[kaDaemon] == nil {
				wdInfo := &WDInfo{}
				server.KaRecvMap[kaDaemon] = wdInfo
			}
			server.KaRecvMap[kaDaemon].RecvedKACount++
			if !server.KaRecvMap[kaDaemon].Active {
				server.KaRecvMap[kaDaemon].Active = true
			}
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
	cmdArgs := []string{"-n", daemon, "-o", op, "-d", cmdDir}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		server.logger.Info(fmt.Sprintln(os.Stderr, "There was an error to ", op, " flexswitch daemon ", daemon, " : ", err))
		return err
	}
	out := string(cmdOut)
	server.logger.Info(fmt.Sprintln("Flexswitch daemon ", daemon, op, " returned ", out))

	if up {
		server.PublishDaemonKANotification(daemon, sysdCommonDefs.KA_UP)
	} else {
		server.PublishDaemonKANotification(daemon, sysdCommonDefs.KA_DOWN)
	}
	return nil
}

func (server *SYSDServer) RestartFlexswitchDaemon(daemon string) error {
	server.ToggleFlexswitchDaemon(daemon, false)
	server.ToggleFlexswitchDaemon(daemon, true)
	return nil
}

func (server *SYSDServer) WDTimer() error {
	server.logger.Info("Starting system WD")
	wdTimer := time.NewTicker(time.Second * KA_TIMEOUT_COUNT)
	for t := range wdTimer.C {
		_ = t
		for daemon, wd := range server.KaRecvMap {
			if wd.RecvedKACount < KA_TIMEOUT_COUNT && wd.RecvedKACount > KA_TIMEOUT_COUNT_MIN {
				server.logger.Info(fmt.Sprintln("Daemon ", daemon, " is slowing down. Monitoring it."))
			}
			if wd.RecvedKACount == KA_TIMEOUT_COUNT_MIN {
				if wd.Active {
					server.logger.Info(fmt.Sprintln("Daemon ", daemon, " is not responsive. Restarting it."))
					go server.RestartFlexswitchDaemon(daemon)
					wd.Active = false
				}
			}
			wd.RecvedKACount = 0
		}
	}
	return nil
}
