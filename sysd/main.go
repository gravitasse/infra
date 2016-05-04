package main

import (
	"flag"
	"fmt"
	"infra/sysd/rpc"
	"infra/sysd/server"
	"utils/dbutils"
	"utils/logging"
)

/*
#cgo LDFLAGS: -L../../../../external/src/github.com/netfilter/libiptables/lib
*/
import "C"

func main() {
	fmt.Println("Starting system daemon")
	paramsDir := flag.String("params", "./params", "Params directory")
	flag.Parse()
	fileName := *paramsDir
	if fileName[len(fileName)-1] != '/' {
		fileName = fileName + "/"
	}

	logger, err := logging.NewLogger("sysd", "SYSTEM", false)
	if err != nil {
		fmt.Println("Failed to start the logger. Nothing will be logged...")
	}
	logger.Info("Started the logger successfully.")

	/*
		dbHdl, err := redis.Dial("tcp", ":6379")
		if err != nil {
			logger.Err("Failed to dial out to Redis server")
			return
		}
	*/
	dbHdl := dbutils.NewDBUtil(logger)
	if err := dbHdl.Connect(); err != nil {
		return
	}

	clientsFileName := fileName + "clients.json"

	logger.Info(fmt.Sprintln("Starting Sysd Server..."))
	sysdServer := server.NewSYSDServer(logger, dbHdl, fileName)
	// Initialize sysd server
	sysdServer.InitServer()
	// Start signal handler first
	go sysdServer.SigHandler(dbHdl)
	// Start sysd server
	go sysdServer.StartServer()
	<-sysdServer.ServerStartedCh

	// Read IpTableAclConfig during restart case
	sysdServer.ReadIpAclConfigFromDB(dbHdl)
	logger.Info(fmt.Sprintln("Starting Sysd Config listener..."))
	confIface := rpc.NewSYSDHandler(logger, sysdServer)
	rpc.StartServer(logger, confIface, clientsFileName)
}
