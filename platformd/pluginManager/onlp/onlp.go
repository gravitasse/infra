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

package onlp

import (
	"errors"
	"fmt"
	"infra/platformd/objects"
	"infra/platformd/pluginManager/pluginCommon"
	"utils/logging"
)

/*
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include "onlp.h"
#include "pluginCommon.h"
*/
import "C"

type onlpDriver struct {
	logger logging.LoggerIntf
}

var driver onlpDriver

func NewONLPPlugin(params *pluginCommon.PluginInitParams) (*onlpDriver, error) {
	driver.logger = params.Logger
	return &driver, nil
}

func (driver *onlpDriver) Init() error {
	driver.logger.Info("Initializing onlp driver")
	rv := int(C.Init())
	if rv < 0 {
		return errors.New("Error initializing Onlp Driver")
	}
	return nil
}

func (driver *onlpDriver) DeInit() error {
	driver.logger.Info("DeInitializing onlp driver")
	C.DeInit()
	return nil
}

func (driver *onlpDriver) GetFanState(fanId int32) (pluginCommon.FanState, error) {
	var retObj pluginCommon.FanState
	var fanInfo C.fan_info_t

	retVal := int(C.GetFanState(&fanInfo, C.int(fanId)))
	if retVal < 0 {
		return retObj, errors.New(fmt.Sprintln("Unable to fetch Fan State of", fanId))
	}
	retObj.FanId = int32(fanInfo.FanId)
	switch int(fanInfo.Mode) {
	case pluginCommon.FAN_MODE_OFF:
		retObj.OperMode = pluginCommon.FAN_MODE_OFF_STR
	case pluginCommon.FAN_MODE_ON:
		retObj.OperMode = pluginCommon.FAN_MODE_ON_STR
	}
	retObj.OperSpeed = int32(fanInfo.Speed)
	//states[idx].OperDirection = fanInfo[idx].Direction
	switch int(fanInfo.Direction) {
	case pluginCommon.FAN_DIR_B2F:
		retObj.OperDirection = pluginCommon.FAN_DIR_B2F_STR
	case pluginCommon.FAN_DIR_F2B:
		retObj.OperDirection = pluginCommon.FAN_DIR_F2B_STR
	case pluginCommon.FAN_DIR_INVALID:
		retObj.OperDirection = pluginCommon.FAN_DIR_INVALID_STR
	}
	//states[idx].Status = fanInfo[idx].Status
	switch int(fanInfo.Status) {
	case pluginCommon.FAN_STATUS_PRESENT:
		retObj.Status = pluginCommon.FAN_STATUS_PRESENT_STR
	case pluginCommon.FAN_STATUS_MISSING:
		retObj.Status = pluginCommon.FAN_STATUS_MISSING_STR
	case pluginCommon.FAN_STATUS_FAILED:
		retObj.Status = pluginCommon.FAN_STATUS_FAILED_STR
	case pluginCommon.FAN_STATUS_NORMAL:
		retObj.Status = pluginCommon.FAN_STATUS_NORMAL_STR
	}
	retObj.Model = C.GoString(&fanInfo.Model[0])
	//states[idx].Model = ""
	retObj.SerialNum = C.GoString(&fanInfo.SerialNum[0])
	return retObj, nil
}

func (driver *onlpDriver) GetFanConfig(fanId int32) (retObj *objects.FanConfig, err error) {
	return retObj, err
}

func (driver *onlpDriver) UpdateFanConfig(cfg *objects.FanConfig) (bool, error) {
	driver.logger.Info("Updating Onlp Fan Config")
	return true, nil
}

func (driver *onlpDriver) GetMaxNumOfFans() int {
	return int(C.GetMaxNumOfFans())
}

func (driver *onlpDriver) GetAllFanState(states []pluginCommon.FanState, cnt int) error {
	var fanInfo []C.fan_info_t

	fanInfo = make([]C.fan_info_t, cnt)
	retVal := int(C.GetAllFanState(&fanInfo[0], C.int(cnt)))
	if retVal < 0 {
		return errors.New(fmt.Sprintln("Unable to fetch the fan State:"))
	}
	for idx := 0; idx < cnt; idx++ {
		if int(fanInfo[idx].valid) == 0 {
			states[idx].Valid = false
			continue
		}
		states[idx].Valid = true
		states[idx].FanId = int32(fanInfo[idx].FanId)
		switch int(fanInfo[idx].Mode) {
		case pluginCommon.FAN_MODE_OFF:
			states[idx].OperMode = pluginCommon.FAN_MODE_OFF_STR
		case pluginCommon.FAN_MODE_ON:
			states[idx].OperMode = pluginCommon.FAN_MODE_ON_STR
		}
		states[idx].OperSpeed = int32(fanInfo[idx].Speed)
		//states[idx].OperDirection = fanInfo[idx].Direction
		switch int(fanInfo[idx].Direction) {
		case pluginCommon.FAN_DIR_B2F:
			states[idx].OperDirection = pluginCommon.FAN_DIR_B2F_STR
		case pluginCommon.FAN_DIR_F2B:
			states[idx].OperDirection = pluginCommon.FAN_DIR_F2B_STR
		case pluginCommon.FAN_DIR_INVALID:
			states[idx].OperDirection = pluginCommon.FAN_DIR_INVALID_STR
		}
		//states[idx].Status = fanInfo[idx].Status
		switch int(fanInfo[idx].Status) {
		case pluginCommon.FAN_STATUS_PRESENT:
			states[idx].Status = pluginCommon.FAN_STATUS_PRESENT_STR
		case pluginCommon.FAN_STATUS_MISSING:
			states[idx].Status = pluginCommon.FAN_STATUS_MISSING_STR
		case pluginCommon.FAN_STATUS_FAILED:
			states[idx].Status = pluginCommon.FAN_STATUS_FAILED_STR
		case pluginCommon.FAN_STATUS_NORMAL:
			states[idx].Status = pluginCommon.FAN_STATUS_NORMAL_STR
		}
		states[idx].Model = C.GoString(&fanInfo[idx].Model[0])
		//states[idx].Model = ""
		states[idx].SerialNum = C.GoString(&fanInfo[idx].SerialNum[0])
	}
	return nil
}

//TODO
func (driver *onlpDriver) GetSfpState(sfpId int32) (pluginCommon.SfpState, error) {
	var retObj pluginCommon.SfpState

	// TODO
	retObj.SfpId = sfpId
	return retObj, nil
}

func (driver *onlpDriver) GetSfpConfig(sfpId int32) (*objects.SfpConfig, error) {
	var retObj objects.SfpConfig

	// TODO
	retObj.SfpId = sfpId
	return &retObj, nil
}

func (driver *onlpDriver) UpdateSfpConfig(cfg *objects.SfpConfig) (bool, error) {
	driver.logger.Info("Updating Onlp SFP Config")
	return true, nil
}

func (driver *onlpDriver) GetAllSfpState(states []pluginCommon.SfpState, cnt int) error {
	driver.logger.Info("GetAllSfpState")
	return nil
}

func (driver *onlpDriver) GetPlatformState() (pluginCommon.PlatformState, error) {
	var retObj pluginCommon.PlatformState
	var sysInfo C.sys_info_t

	rt := int(C.GetPlatformState(&sysInfo))

	if rt < 0 {
		return retObj, errors.New(fmt.Sprintln("Unable to fetch System info"))
	}

	retObj.ProductName = C.GoString(&sysInfo.product_name[0])
	retObj.Vendor = C.GoString(&sysInfo.vendor[0])
	retObj.SerialNum = C.GoString(&sysInfo.serial_number[0])
	retObj.Manufacturer = C.GoString(&sysInfo.manufacturer[0])
	retObj.Release = C.GoString(&sysInfo.label_revision[0])
	retObj.PlatformName = C.GoString(&sysInfo.platform_name[0])
	retObj.Version = C.GoString(&sysInfo.onie_version[0])

	return retObj, nil
}

func (driver *onlpDriver) GetMaxNumOfThermal() int {
	return 0
}

func (driver *onlpDriver) GetThermalState(thermalId int32) (tState pluginCommon.ThermalState, err error) {
	return tState, err
}
