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

package rpc

import (
	"errors"
	"fmt"
	"infra/sysd/server"
	"sysd"
	"utils/logging"
)

type SYSDHandler struct {
	server *server.SYSDServer
	logger *logging.Writer
}

func NewSYSDHandler(logger *logging.Writer, server *server.SYSDServer) *SYSDHandler {
	h := new(SYSDHandler)
	h.server = server
	h.logger = logger
	return h
}

func (h *SYSDHandler) SendGlobalLoggingConfig(gLoggingConfig *sysd.SystemLogging) bool {
	var Logging bool
	if gLoggingConfig.Logging == "on" {
		Logging = true
	} else {
		Logging = false
	}
	gConf := server.GlobalLoggingConfig{
		Enable: Logging,
	}
	h.server.GlobalLoggingConfigCh <- gConf
	return true
}

func (h *SYSDHandler) SendComponentLoggingConfig(cLoggingConfig *sysd.ComponentLogging) bool {
	cConf := server.ComponentLoggingConfig{
		Component: cLoggingConfig.Module,
		Level:     logging.ConvertLevelStrToVal(cLoggingConfig.Level),
	}
	h.server.ComponentLoggingConfigCh <- cConf
	return true
}

func (h *SYSDHandler) UpdateSystemLogging(origConf *sysd.SystemLogging, newConf *sysd.SystemLogging, attrset []bool) (bool, error) {
	h.logger.Info(fmt.Sprintln("Original global config attrs:", origConf))
	if newConf == nil {
		err := errors.New("Invalid global Configuration")
		return false, err
	}
	h.logger.Info(fmt.Sprintln("Update global config attrs:", newConf))
	return h.SendGlobalLoggingConfig(newConf), nil
}

func (h *SYSDHandler) UpdateComponentLogging(origConf *sysd.ComponentLogging, newConf *sysd.ComponentLogging, attrset []bool) (bool, error) {
	h.logger.Info(fmt.Sprintln("Original component config attrs:", origConf))
	if newConf == nil {
		err := errors.New("Invalid component Configuration")
		return false, err
	}
	h.logger.Info(fmt.Sprintln("Update component config attrs:", newConf))
	return h.SendComponentLoggingConfig(newConf), nil
}

func (h *SYSDHandler) CreateSystemLogging(gLoggingConf *sysd.SystemLogging) (bool, error) {
	h.logger.Info(fmt.Sprintln("Create global config attrs:", gLoggingConf))
	return h.SendGlobalLoggingConfig(gLoggingConf), nil
}

func (h *SYSDHandler) CreateComponentLogging(cLoggingConf *sysd.ComponentLogging) (bool, error) {
	h.logger.Info(fmt.Sprintln("Create component config attrs:", cLoggingConf))
	return h.SendComponentLoggingConfig(cLoggingConf), nil
}

func (h *SYSDHandler) DeleteSystemLogging(gLoggingConf *sysd.SystemLogging) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete global config attrs:", gLoggingConf))
	err := errors.New("Not supported")
	return false, err
}

func (h *SYSDHandler) DeleteComponentLogging(cLoggingConf *sysd.ComponentLogging) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete component config attrs:", cLoggingConf))
	err := errors.New("Not supported")
	return false, err
}

func (h *SYSDHandler) CreateIpTableAcl(ipaclConfig *sysd.IpTableAcl) (bool, error) {
	h.logger.Info("Create Ip Table rule " + ipaclConfig.Name)
	h.server.IptableAddCh <- ipaclConfig
	return true, nil
	//return (h.server.AddIpTableRule(ipaclConfig, false /* non - restart*/))
}

func (h *SYSDHandler) UpdateIpTableAcl(origConf *sysd.IpTableAcl,
	newConf *sysd.IpTableAcl, attrset []bool) (bool, error) {
	err := errors.New("Not supported")
	return false, err
}

func (h *SYSDHandler) DeleteIpTableAcl(ipaclConfig *sysd.IpTableAcl) (bool, error) {
	h.logger.Info("Delete Ip Table rule " + ipaclConfig.Name)
	//return (h.server.DelIpTableRule(ipaclConfig))
	h.server.IptableDelCh <- ipaclConfig
	return true, nil
}

func (h *SYSDHandler) ExecuteActionDaemon(daemonConfig *sysd.Daemon) (bool, error) {
	h.logger.Info(fmt.Sprintln("Daemon action attrs: ", daemonConfig))
	dConf := server.DaemonConfig{
		Name:     daemonConfig.Name,
		Enable:   daemonConfig.Enable,
		WatchDog: daemonConfig.WatchDog,
	}
	h.server.DaemonConfigCh <- dConf
	return true, nil
}

func (h *SYSDHandler) GetDaemonState(name string) (*sysd.DaemonState, error) {
	h.logger.Info(fmt.Sprintln("Get Daemon state ", name))
	daemonStateResponse := sysd.NewDaemonState()
	dState := h.server.GetDaemonState(name)
	daemonState := h.server.ConvertDaemonStateToThrift(*dState)
	daemonStateResponse = daemonState
	return daemonStateResponse, nil
}

func (h *SYSDHandler) GetBulkDaemonState(fromIdx sysd.Int, count sysd.Int) (*sysd.DaemonStateGetInfo, error) {
	h.logger.Info(fmt.Sprintln("Get Daemon states "))
	nextIdx, currCount, daemonStates := h.server.GetBulkDaemonStates(int(fromIdx), int(count))
	if daemonStates == nil {
		err := errors.New("System server is busy")
		return nil, err
	}
	daemonStatesResponse := make([]*sysd.DaemonState, len(daemonStates))
	for idx, item := range daemonStates {
		daemonStatesResponse[idx] = h.server.ConvertDaemonStateToThrift(item)
	}
	daemonStateGetInfo := sysd.NewDaemonStateGetInfo()
	daemonStateGetInfo.Count = sysd.Int(currCount)
	daemonStateGetInfo.StartIdx = sysd.Int(fromIdx)
	daemonStateGetInfo.EndIdx = sysd.Int(nextIdx)
	daemonStateGetInfo.More = (nextIdx != 0)
	daemonStateGetInfo.DaemonStateList = daemonStatesResponse
	return daemonStateGetInfo, nil
}

func (h *SYSDHandler) PeriodicKeepAlive(name string) error {
	h.server.KaRecvCh <- name
	return nil
}
