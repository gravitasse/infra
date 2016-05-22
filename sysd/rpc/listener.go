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
	"models"
	"strings"
	"sysd"
	"utils/logging"
)

const (
	SPECIAL_HOSTNAME_CHARS = "_"
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

func (h *SYSDHandler) UpdateSystemLogging(origConf *sysd.SystemLogging, newConf *sysd.SystemLogging, attrset []bool, op string) (bool, error) {
	h.logger.Info(fmt.Sprintln("Original global config attrs:", origConf))
	if newConf == nil {
		err := errors.New("Invalid global Configuration")
		return false, err
	}
	h.logger.Info(fmt.Sprintln("Update global config attrs:", newConf))
	return h.SendGlobalLoggingConfig(newConf), nil
}

func (h *SYSDHandler) UpdateComponentLogging(origConf *sysd.ComponentLogging, newConf *sysd.ComponentLogging, attrset []bool, op string) (bool, error) {
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
}

func (h *SYSDHandler) UpdateIpTableAcl(origConf *sysd.IpTableAcl, newConf *sysd.IpTableAcl,
	attrset []bool, op string) (bool, error) {
	err := errors.New("Not supported")
	return false, err
}

func (h *SYSDHandler) DeleteIpTableAcl(ipaclConfig *sysd.IpTableAcl) (bool, error) {
	h.logger.Info("Delete Ip Table rule " + ipaclConfig.Name)
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

func convertSystemParamThriftToModel(cfg *sysd.SystemParam) models.SystemParam {
	confg := models.SystemParam{
		Description: cfg.Description,
		Version:     cfg.Version,
		MgmtIp:      cfg.MgmtIp,
		RouterId:    cfg.RouterId,
		Hostname:    cfg.Hostname,
		SwitchMac:   cfg.SwitchMac,
		Vrf:         cfg.Vrf,
	}
	return confg
}

func (h *SYSDHandler) CreateSystemParam(cfg *sysd.SystemParam) (bool, error) {
	if h.server.SystemInfoCreated() {
		return false, errors.New("System Params Info is already created for Default VRF, please do update to modify params")
	}
	if strings.ContainsAny(cfg.Hostname, SPECIAL_HOSTNAME_CHARS) {
		return false, errors.New("Hostname should not have Special Characters " +
			SPECIAL_HOSTNAME_CHARS)
	}
	h.logger.Info(fmt.Sprintln("Configuring Global Object", cfg))
	confg := convertSystemParamThriftToModel(cfg)
	h.server.SystemParamConfig <- confg
	return true, nil
}

func (h *SYSDHandler) validatUpdateSystemParam(newCfg *sysd.SystemParam, attrset []bool) ([]string, error) {
	var updatedInfo []string
	/*
		1 : string Vrf
		2 : string MgmtIp
		3 : string Hostname
		4 : string RouterId
		5 : string Version
		6 : string SwitchMac
		7 : string Description
	*/
	for idx, _ := range attrset {
		if attrset[idx] == false {
			continue
		}
		switch idx {
		case 0:
			return updatedInfo, errors.New("VRF update is not supported")
		case 1:
			return updatedInfo, errors.New("MgmtIp update is not supported")
		case 2:
			if strings.ContainsAny(newCfg.Hostname, SPECIAL_HOSTNAME_CHARS) {
				return updatedInfo, errors.New("Hostname should not have Special Characters " +
					SPECIAL_HOSTNAME_CHARS)
			}
			updatedInfo = append(updatedInfo, "Hostname")
		case 3:
			return updatedInfo, errors.New("Router ID update is not supported")
		case 4:
			return updatedInfo, errors.New("Version update is not supported")
		case 5:
			return updatedInfo, errors.New("Switch Mac Address update is not supported")
		case 6:
			updatedInfo = append(updatedInfo, "Description")
		}
	}
	return updatedInfo, nil
}

func (h *SYSDHandler) UpdateSystemParam(org *sysd.SystemParam, new *sysd.SystemParam, attrset []bool,
	op string) (bool, error) {
	h.logger.Info(fmt.Sprintln("Received update for system param information", org, new, attrset))
	if org == nil || new == nil {
		return false, errors.New("Invalid information provided to server")
	}
	entriesUpdated, err := h.validatUpdateSystemParam(new, attrset)
	if err != nil {
		return false, err
	}
	cfg := convertSystemParamThriftToModel(new)
	updInfo := server.SystemParamUpdate{
		EntriesUpdated: entriesUpdated,
		NewCfg:         &cfg,
	}
	h.server.SysUpdCh <- &updInfo
	return true, nil
}

func (h *SYSDHandler) DeleteSystemParam(cfg *sysd.SystemParam) (bool, error) {
	return false, errors.New("Delete of system params for default vrf is not supported")
}

func convertSystemParamStateToThrift(info models.SystemParamState, entry *sysd.SystemParamState) {
	entry.Vrf = string(info.Vrf)
	entry.SwitchMac = string(info.SwitchMac)
	entry.RouterId = string(info.RouterId)
	entry.MgmtIp = string(info.MgmtIp)
	entry.Version = string(info.Version)
	entry.Description = string(info.Description)
	entry.Hostname = string(info.Hostname)
}

func (h *SYSDHandler) GetSystemParamState(name string) (*sysd.SystemParamState, error) {
	h.logger.Info("Get system params info for " + name)
	sysParamsResp := sysd.NewSystemParamState()
	sysInfo := h.server.GetSystemParam(name)
	if sysInfo == nil {
		return nil, errors.New("No Matching entry for " + name + "found")
	}
	convertSystemParamStateToThrift(*sysInfo, sysParamsResp)
	h.logger.Info(fmt.Sprintln("Returing System Info:", sysParamsResp))
	return sysParamsResp, nil
}

func (h *SYSDHandler) GetBulkSystemParamState(fromIdx sysd.Int, count sysd.Int) (*sysd.SystemParamStateGetInfo, error) {
	//@TODO: when we support vrf change get bulk... today only one system info is present
	sysParamsResp, err := h.GetSystemParamState("default")
	if err != nil {
		return nil, err
	}
	systemGetInfoResp := sysd.NewSystemParamStateGetInfo()
	systemGetInfoResp.Count = 1
	systemGetInfoResp.StartIdx = 0
	systemGetInfoResp.EndIdx = 1
	systemGetInfoResp.More = false
	respList := make([]*sysd.SystemParamState, 1)
	//respList = append(respList, sysParamsResp)
	respList[0] = sysParamsResp
	systemGetInfoResp.SystemParamStateList = respList
	return systemGetInfoResp, nil
}

/*********************** SYSD INTERNAL API **************************/

func (h *SYSDHandler) PeriodicKeepAlive(name string) error {
	h.server.KaRecvCh <- name
	return nil
}
