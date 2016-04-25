package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	nanomsg "github.com/op/go-nanomsg"
	"infra/sysd/iptables"
	"infra/sysd/sysdCommonDefs"
	"os"
	"os/signal"
	"syscall"
	"sysd"
	"utils/logging"
)

const (
	SYSD_TOTAL_DB_USERS = 2
)

type GlobalLoggingConfig struct {
	Enable bool
}

type ComponentLoggingConfig struct {
	Component string
	Level     sysdCommonDefs.SRDebugLevel
}

type SYSDServer struct {
	logger                   *logging.Writer
	paramsDir                string
	GlobalLoggingConfigCh    chan GlobalLoggingConfig
	ComponentLoggingConfigCh chan ComponentLoggingConfig
	sysdPubSocket            *nanomsg.PubSocket
	sysdIpTableMgr           *ipTable.SysdIpTableHandler
	notificationCh           chan []byte
	IptableAddCh             chan *sysd.IpTableAcl
	IptableDelCh             chan *sysd.IpTableAcl
	dbUserCh                 chan int
	KaRecvCh                 chan string
	KaRecvMap                map[string]*WDInfo
}

func NewSYSDServer(logger *logging.Writer) *SYSDServer {
	sysdServer := &SYSDServer{}
	sysdServer.sysdIpTableMgr = ipTable.SysdNewSysdIpTableHandler(logger)
	sysdServer.logger = logger
	sysdServer.GlobalLoggingConfigCh = make(chan GlobalLoggingConfig)
	sysdServer.ComponentLoggingConfigCh = make(chan ComponentLoggingConfig)
	sysdServer.notificationCh = make(chan []byte)
	sysdServer.IptableAddCh = make(chan *sysd.IpTableAcl)
	sysdServer.IptableDelCh = make(chan *sysd.IpTableAcl)
	sysdServer.dbUserCh = make(chan int, 1)
	return sysdServer
}

func (server *SYSDServer) SigHandler() {
	server.logger.Info(fmt.Sprintln("Starting SigHandler"))
	sigChan := make(chan os.Signal, 1)
	signalList := []os.Signal{syscall.SIGHUP}
	signal.Notify(sigChan, signalList...)

	for {
		select {
		case signal := <-sigChan:
			switch signal {
			case syscall.SIGHUP:
				server.logger.Info("Received SIGHUP signal. Exiting")
				os.Exit(0)
			default:
				server.logger.Info(fmt.Sprintln("Unhandled signal : ", signal))
			}
		}
	}
}

func (server *SYSDServer) InitServer(paramsDir string) {
	server.logger.Info(fmt.Sprintln("Initializing Sysd Server"))
	server.paramsDir = paramsDir
}

func (server *SYSDServer) InitPublisher(pub_str string) (pub *nanomsg.PubSocket) {
	server.logger.Info(fmt.Sprintln("Setting up ", pub_str, "publisher"))
	pub, err := nanomsg.NewPubSocket()
	if err != nil {
		server.logger.Info(fmt.Sprintln("Failed to open pub socket"))
		return nil
	}
	ep, err := pub.Bind(pub_str)
	if err != nil {
		server.logger.Info(fmt.Sprintln("Failed to bind pub socket - ", ep))
		return nil
	}
	err = pub.SetSendBuffer(1024)
	if err != nil {
		server.logger.Info(fmt.Sprintln("Failed to set send buffer size"))
		return nil
	}
	return pub
}

func (server *SYSDServer) PublishSysdNotifications() {
	server.sysdPubSocket = server.InitPublisher(sysdCommonDefs.PUB_SOCKET_ADDR)
	for {
		select {
		case event := <-server.notificationCh:
			server.logger.Info(fmt.Sprintln("Received call to notify ", event))
			_, err := server.sysdPubSocket.Send(event, nanomsg.DontWait)
			if err == syscall.EAGAIN {
				server.logger.Err(fmt.Sprintln("Failed to publish event"))
			}
		}
	}
}

func (server *SYSDServer) ReadGlobalLoggingConfigFromDB(dbHdl *sql.DB) error {
	return nil
}

func (server *SYSDServer) ReadComponentLoggingConfigFromDB(dbHdl *sql.DB) error {
	return nil
}

func (server *SYSDServer) ReadConfigFromDB(dbHdl *sql.DB) error {
	var err error
	// BfdGlobalConfig
	err = server.ReadGlobalLoggingConfigFromDB(dbHdl)
	if err != nil {
		server.dbUserCh <- 1
		return err
	}
	// BfdIntfConfig
	err = server.ReadComponentLoggingConfigFromDB(dbHdl)
	if err != nil {
		server.dbUserCh <- 1
		return err
	}
	server.dbUserCh <- 1
	return nil
}

func (server *SYSDServer) ProcessGlobalLoggingConfig(gLogConf GlobalLoggingConfig) error {
	server.logger.SetGlobal(gLogConf.Enable)
	msg := sysdCommonDefs.GlobalLogging{
		Enable: gLogConf.Enable,
	}
	msgBuf, err := json.Marshal(msg)
	if err != nil {
		server.logger.Err("Failed to marshal Global logging message")
		return err
	}
	notification := sysdCommonDefs.Notification{
		Type:    uint8(sysdCommonDefs.G_LOG),
		Payload: msgBuf,
	}
	notificationBuf, err := json.Marshal(notification)
	if err != nil {
		server.logger.Err("Failed to marshal Global logging message")
		return err
	}
	server.notificationCh <- notificationBuf
	return nil
}

func (server *SYSDServer) ProcessComponentLoggingConfig(cLogConf ComponentLoggingConfig) error {
	if cLogConf.Component == server.logger.MyComponentName {
		server.logger.SetLevel(cLogConf.Level)
	} else {
		msg := sysdCommonDefs.ComponentLogging{
			Name:  cLogConf.Component,
			Level: cLogConf.Level,
		}
		msgBuf, err := json.Marshal(msg)
		if err != nil {
			server.logger.Err("Failed to marshal Global logging message")
			return err
		}
		notification := sysdCommonDefs.Notification{
			Type:    uint8(sysdCommonDefs.C_LOG),
			Payload: msgBuf,
		}
		notificationBuf, err := json.Marshal(notification)
		if err != nil {
			server.logger.Err("Failed to marshal Global logging message")
			return err
		}
		server.notificationCh <- notificationBuf
	}
	return nil
}

func (server *SYSDServer) StartServer(paramFile string, dbHdl *sql.DB) {
	// Start signal handler first
	go server.SigHandler()
	// Start notification publish thread
	go server.PublishSysdNotifications()
	// Read configurations already present in DB
	go server.ReadConfigFromDB(dbHdl)
	// Read IpTableAclConfig during restart case
	go server.ReadIpAclConfigFromDB(dbHdl)
	// Start watchdog routine
	go server.StartWDRoutine()
	users := 0
	// Now, wait on below channels to process
	for {
		select {
		case gLogConf := <-server.GlobalLoggingConfigCh:
			server.logger.Info(fmt.Sprintln("Received call for performing Global logging Configuration",
				gLogConf))
			server.ProcessGlobalLoggingConfig(gLogConf)
		case compLogConf := <-server.ComponentLoggingConfigCh:
			server.logger.Info(fmt.Sprintln("Received call for performing Component logging Configuration",
				compLogConf))
			server.ProcessComponentLoggingConfig(compLogConf)
		case addConfig := <-server.IptableAddCh:
			server.sysdIpTableMgr.AddIpRule(addConfig, false /*non-restart*/)
		case delConfig := <-server.IptableDelCh:
			server.sysdIpTableMgr.DelIpRule(delConfig)
		case totalUsers := <-server.dbUserCh:
			users = totalUsers + users
			if users == SYSD_TOTAL_DB_USERS {
				server.logger.Info("Closing DB as all the db users are done")
				dbHdl.Close()
			}
		}
	}
}
