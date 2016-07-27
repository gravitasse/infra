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

package openBMC

import (
	"infra/platformd/objects"
	"infra/platformd/pluginManager/pluginCommon"
	"utils/logging"
)

type openBMCDriver struct {
	logger logging.LoggerIntf
}

var driver openBMCDriver

func NewOpenBMCPlugin(params *pluginCommon.PluginInitParams) *openBMCDriver {
	driver.logger = params.Logger
	return &driver
}

func (driver *openBMCDriver) Init() error {
	driver.logger.Info("Initializing openBMC driver")
	return nil
}

func (driver *openBMCDriver) DeInit() error {
	driver.logger.Info("DeInitializing openBMC driver")
	return nil
}

func (driver *openBMCDriver) GetFanState(fanId int32) (pluginCommon.FanState, error) {
	var retObj pluginCommon.FanState
	retObj.FanId = fanId
	retObj.OperMode = "ON"
	retObj.OperSpeed = 10000
	retObj.OperDirection = "B2F"
	retObj.Status = "PRESENT"
	retObj.Model = "OPENBMC"
	retObj.SerialNum = "AABBCC112233"
	retObj.Valid = true
	return retObj, nil
}

func (driver *openBMCDriver) GetFanConfig(fanId int32) (*objects.FanConfig, error) {
	var retObj objects.FanConfig
	retObj.FanId = fanId
	retObj.AdminSpeed = 10000
	retObj.AdminDirection = "B2F"
	return &retObj, nil
}

func (driver *openBMCDriver) UpdateFanConfig(cfg *objects.FanConfig) (bool, error) {
	driver.logger.Info("Updating OpenBMC Fan Config")
	return true, nil
}

func (driver *openBMCDriver) GetMaxNumOfFans() int {
	driver.logger.Info("Inside OpenBMC: GetMaxNumOfFans()")
	return 0
}

func (driver *openBMCDriver) GetAllFanState(state []pluginCommon.FanState, cnt int) error {
	return nil
}
