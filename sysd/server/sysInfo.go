//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
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
	"encoding/json"
	"fmt"
	"infra/sysd/sysdCommonDefs"
	"models/objects"
	"os/exec"
)

type SystemParamUpdate struct {
	EntriesUpdated []string
	NewCfg         *objects.SystemParam
}

func (svr *SYSDServer) ReadSystemInfoFromDB() error {
	svr.logger.Info("Reading System Information From Db")
	dbHdl := svr.dbHdl
	if dbHdl != nil {
		var dbObj objects.SystemParam
		objList, err := dbHdl.GetAllObjFromDb(dbObj)
		if err != nil {
			svr.logger.Err("DB query failed for System Info")
			return err
		}
		svr.logger.Info(fmt.Sprintln("Total System Entries are", len(objList)))
		for idx := 0; idx < len(objList); idx++ {
			dbObject := objList[idx].(objects.SystemParam)
			svr.SysInfo.SwitchMac = dbObject.SwitchMac
			svr.SysInfo.MgmtIp = dbObject.MgmtIp
			svr.SysInfo.Version = dbObject.Version
			svr.SysInfo.Description = dbObject.Description
			svr.SysInfo.Hostname = dbObject.Hostname
			break
		}
	}
	svr.logger.Info("reading system info from db done")
	return nil
}

// Func to send update nanomsg update notification to all the dameons on the system
func (svr *SYSDServer) SendSystemUpdate() ([]byte, error) {
	msgBuf, err := json.Marshal(svr.SysInfo)
	if err != nil {
		return nil, err
	}
	notification := sysdCommonDefs.Notification{
		Type:    uint8(sysdCommonDefs.SYSTEM_Info),
		Payload: msgBuf,
	}
	notificationBuf, err := json.Marshal(notification)
	if err != nil {
		return nil, err
	}
	return notificationBuf, nil
}

func (svr *SYSDServer) copyAndSendSystemParam(cfg objects.SystemParam) {
	sysInfo := svr.SysInfo
	sysInfo.SwitchMac = cfg.SwitchMac
	sysInfo.MgmtIp = cfg.MgmtIp
	sysInfo.Version = cfg.Version
	sysInfo.Description = cfg.Description
	sysInfo.Hostname = cfg.Hostname
	sysInfo.Vrf = cfg.Vrf
	svr.SendSystemUpdate()
}

// Initialize system information using json file...or whatever other means are
func (svr *SYSDServer) InitSystemInfo(cfg objects.SystemParam) {
	svr.SysInfo = &objects.SystemParam{}
	svr.copyAndSendSystemParam(cfg)
}

// Helper function for NB listener to determine whether a global object is created or not
func (svr *SYSDServer) SystemInfoCreated() bool {
	return (svr.SysInfo != nil)
}

// During Get calls we will use below api to read from run-time information
func (svr *SYSDServer) GetSystemParam(name string) *objects.SystemParamState {
	if svr.SysInfo == nil || svr.SysInfo.Vrf != name {
		return nil
	}
	sysParamsInfo := new(objects.SystemParamState)

	sysParamsInfo.Vrf = svr.SysInfo.Vrf
	sysParamsInfo.SwitchMac = svr.SysInfo.SwitchMac
	sysParamsInfo.MgmtIp = svr.SysInfo.MgmtIp
	sysParamsInfo.Version = svr.SysInfo.Version
	sysParamsInfo.Description = svr.SysInfo.Description
	sysParamsInfo.Hostname = svr.SysInfo.Hostname
	return sysParamsInfo
}

// Update runtime system param info and send a notification
func (svr *SYSDServer) UpdateSystemInfo(updateInfo *SystemParamUpdate) {
	svr.copyAndSendSystemParam(*updateInfo.NewCfg)
	for _, entry := range updateInfo.EntriesUpdated {
		switch entry {
		case "Hostname":
			binary, lookErr := exec.LookPath("hostname")
			if lookErr != nil {
				svr.logger.Err(fmt.Sprintln("Error searching path for hostname", lookErr))
				continue
			}
			cmd := exec.Command(binary, updateInfo.NewCfg.Hostname)
			err := cmd.Run()
			if err != nil {
				svr.logger.Err(fmt.Sprintln("Updating hostname in linux failed", err))
			}
		}
	}
}
