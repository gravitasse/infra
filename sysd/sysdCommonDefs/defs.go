package sysdCommonDefs

import ()

const (
	PUB_SOCKET_ADDR = "ipc:///tmp/sysd.ipc"
)

const (
	G_LOG     uint8 = 1 // Global logging configuration
	C_LOG     uint8 = 2 // Component level logging configuration
	KA_DAEMON uint8 = 3 // Daemon keepalive notification
)

type Notification struct {
	Type    uint8
	Payload []byte
}

//Logging levels
type SRDebugLevel uint8

const (
	CRIT   SRDebugLevel = 0
	ERR    SRDebugLevel = 1
	WARN   SRDebugLevel = 2
	ALERT  SRDebugLevel = 3
	EMERG  SRDebugLevel = 4
	NOTICE SRDebugLevel = 5
	INFO   SRDebugLevel = 6
	DEBUG  SRDebugLevel = 7
	TRACE  SRDebugLevel = 8
)

type GlobalLogging struct {
	Enable bool
}

type ComponentLogging struct {
	Name  string
	Level SRDebugLevel
}

type SRDaemonStatus uint8

const (
	SYSD_TOTAL_KA_DAEMONS = 32
)

const (
	KA_UP   SRDaemonStatus = 0
	KA_DOWN SRDaemonStatus = 1
)

type DaemonStatus struct {
	Name   string
	Status SRDaemonStatus
}
