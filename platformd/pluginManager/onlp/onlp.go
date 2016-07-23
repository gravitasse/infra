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
	"infra/platformd/objects"
	"infra/platformd/pluginManager/pluginCommon"
	"utils/logging"
)

type onlpDriver struct {
	logger logging.LoggerIntf
}

var driver onlpDriver

func NewONLPPlugin(params *pluginCommon.PluginInitParams) *onlpDriver {
	driver.logger = params.Logger
	return &driver
}

func (driver *onlpDriver) Init() error {
	driver.logger.Info("Initializing onlp driver")
	return nil
}

func (driver *onlpDriver) DeInit() error {
	driver.logger.Info("DeInitializing onlp driver")
	return nil
}

func (driver *onlpDriver) GetFanState(fanId int32) (*objects.FanState, error) {
	var retObj objects.FanState
	retObj.FanId = fanId
	retObj.OperMode = "ON"
	retObj.OperSpeed = 10000
	retObj.OperDirection = "B2F"
	retObj.Status = "PRESENT"
	retObj.Model = "ONLP"
	retObj.SerialNum = "AABBCC112233"
	return &retObj, nil
}

func (driver *onlpDriver) GetFanConfig(fanId int32) (*objects.FanConfig, error) {
	var retObj objects.FanConfig
	retObj.FanId = fanId
	retObj.AdminSpeed = 10000
	retObj.AdminDirection = "B2F"
	return &retObj, nil
}

func (driver *onlpDriver) UpdateFanConfig(cfg *objects.FanConfig) (bool, error) {
	driver.logger.Info("Updating Onlp Fan Config")
	return true, nil
}
