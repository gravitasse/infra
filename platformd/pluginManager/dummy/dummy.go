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

package dummy

import (
	"infra/platformd/objects"
	"infra/platformd/pluginManager/pluginCommon"
	"utils/logging"
)

type dummyDriver struct {
	logger logging.LoggerIntf
}

var driver dummyDriver

func NewDummyPlugin(params *pluginCommon.PluginInitParams) (*dummyDriver, error) {
	driver.logger = params.Logger
	return &driver, nil
}

func (driver *dummyDriver) Init() error {
	driver.logger.Info("Initializing Dummy driver")
	return nil
}

func (driver *dummyDriver) DeInit() error {
	driver.logger.Info("DeInitializing Dummy driver")
	return nil
}

func (driver *dummyDriver) GetFanState(fanId int32) (retObj pluginCommon.FanState, err error) {
	return retObj, nil
}

func (driver *dummyDriver) GetFanConfig(fanId int32) (retObj *objects.FanConfig, err error) {
	return retObj, nil
}

func (driver *dummyDriver) UpdateFanConfig(cfg *objects.FanConfig) (bool, error) {
	driver.logger.Info("Updating Dummy Fan Config")
	return true, nil
}

func (driver *dummyDriver) GetMaxNumOfFans() int {
	driver.logger.Info("Inside Dummy: GetMaxNumOfFans()")
	return 0
}

func (driver *dummyDriver) GetAllFanState(state []pluginCommon.FanState, cnt int) error {
	return nil
}

func (driver *dummyDriver) GetSfpState(sfpId int32) (retObj pluginCommon.SfpState, err error) {
	// TODO
	retObj.SfpId = sfpId
	return retObj, nil
}

func (driver *dummyDriver) GetSfpConfig(sfpId int32) (retObj *objects.SfpConfig, err error) {
	// TODO
	retObj.SfpId = sfpId
	return retObj, nil
}

func (driver *dummyDriver) UpdateSfpConfig(cfg *objects.SfpConfig) (bool, error) {
	driver.logger.Info("Updating Onlp SFP Config")
	return true, nil
}

func (driver *dummyDriver) GetAllSfpState(states []pluginCommon.SfpState, cnt int) error {
	driver.logger.Info("GetAllSfpState")
	return nil
}

func (driver *dummyDriver) GetSfpCnt() int {
	driver.logger.Info("GetSfpCnt")
	return 0
}

func (driver *dummyDriver) GetPlatformState() (pluginCommon.PlatformState, error) {
	var retObj pluginCommon.PlatformState

	return retObj, nil
}
func (driver *dummyDriver) GetThermalState(thermalId int32) (retObj pluginCommon.ThermalState, err error) {
	return retObj, nil
}

func (driver *dummyDriver) GetMaxNumOfThermal() int {
	driver.logger.Info("Inside Dummy: GetMaxNumOfThermal()")
	return 0
}
