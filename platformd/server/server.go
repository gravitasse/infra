//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package server

import (
	"fmt"
	//"infra/platformd/objects"
	"infra/platformd/pluginManager"
	"infra/platformd/pluginManager/pluginCommon"
	"time"
	"utils/dbutils"
	"utils/logging"
)

type PlatformdServer struct {
	dmnName        string
	paramsDir      string
	pluginMgr      *pluginManager.PluginManager
	dbHdl          *dbutils.DBUtil
	Logger         logging.LoggerIntf
	InitCompleteCh chan bool
	ReqChan        chan *ServerRequest
	ReplyChan      chan interface{}
}

type InitParams struct {
	DmnName     string
	ParamsDir   string
	CfgFileName string
	DbHdl       *dbutils.DBUtil
	Logger      logging.LoggerIntf
}

func NewPlatformdServer(initParams *InitParams) *PlatformdServer {
	var svr PlatformdServer

	svr.dmnName = initParams.DmnName
	svr.paramsDir = initParams.ParamsDir
	svr.dbHdl = initParams.DbHdl
	svr.Logger = initParams.Logger
	svr.InitCompleteCh = make(chan bool)
	svr.ReqChan = make(chan *ServerRequest)
	svr.ReplyChan = make(chan interface{})

	CfgFileInfo, err := parseCfgFile(initParams.CfgFileName)
	if err != nil {
		svr.Logger.Err("Failed to parse platformd config file, using default values for all attributes")
	}
	pluginInitParams := &pluginCommon.PluginInitParams{
		Logger: svr.Logger,
	}
	svr.pluginMgr = pluginManager.NewPluginMgr(CfgFileInfo.PluginName, pluginInitParams)
	return &svr
}

func (svr *PlatformdServer) initServer() error {
	//Initialize plugin layer first
	err := svr.pluginMgr.Init()
	if err != nil {
		return err
	}

	return err
}

func (svr *PlatformdServer) handleRPCRequest(req *ServerRequest) {
	svr.Logger.Info(fmt.Sprintln("Calling handle RPC Request for:", *req))
	switch req.Op {
	case GET_FAN_STATE:
		var retObj GetFanStateOutArgs
		if val, ok := req.Data.(*GetFanStateInArgs); ok {
			retObj.Obj, retObj.Err = svr.getFanState(val.FanId)
		}
		svr.Logger.Info(fmt.Sprintln("Server GET_FAN_STATE request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_FAN_STATE:
		var retObj GetBulkFanStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkFanState(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case GET_FAN_CONFIG:
		var retObj GetFanConfigOutArgs
		if val, ok := req.Data.(*GetFanConfigInArgs); ok {
			retObj.Obj, retObj.Err = svr.getFanConfig(val.FanId)
		}
		svr.Logger.Info(fmt.Sprintln("Server GET_FAN_CONFIG request replying -", retObj))
		svr.ReplyChan <- interface{}(&retObj)
	case GET_BULK_FAN_CONFIG:
		var retObj GetBulkFanConfigOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = svr.getBulkFanConfig(val.FromIdx, val.Count)
		}
		svr.ReplyChan <- interface{}(&retObj)
	case UPDATE_FAN_CONFIG:
		var retObj UpdateConfigOutArgs
		if val, ok := req.Data.(*UpdateFanConfigInArgs); ok {
			retObj.RetVal, retObj.Err = svr.updateFanConfig(val.FanOldCfg, val.FanNewCfg, val.AttrSet)
		}
		svr.ReplyChan <- interface{}(&retObj)
	default:
		svr.Logger.Err(fmt.Sprintln("Error : Server recevied unrecognized request - ", req.Op))
	}
}

func (svr *PlatformdServer) Serve() {
	svr.Logger.Info("Server initialization started")
	err := svr.initServer()
	if err != nil {
		panic(err)
	}
	svr.InitCompleteCh <- true
	svr.Logger.Info("Server initialization complete, starting cfg/state listerner")
	for {
		select {
		case req := <-svr.ReqChan:
			svr.Logger.Info(fmt.Sprintln("Server request received - ", *req))
			svr.handleRPCRequest(req)

		}
		time.Sleep(10)
	}
}
