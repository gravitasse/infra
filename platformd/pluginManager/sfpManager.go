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

package pluginManager

import (
	"infra/platformd/objects"
	"utils/logging"
)

type SfpManager struct {
	logger logging.LoggerIntf
	plugin PluginIntf
}

/*
type SfpConfig struct {
	SfpId      int32
	AdminState string
}

type SfpState struct {
	SfpId   int32
	Status  string
	SfpLOS  string
	SfpType string
	EEPROM  string
}
*/

var SfpMgr SfpManager

func (sfpMgr *SfpManager) Init(logger logging.LoggerIntf, plugin PluginIntf) {
	sfpMgr.logger = logger
	sfpMgr.plugin = plugin
	sfpMgr.logger.Info("SFP Manager Init()")
}

func (sfpMgr *SfpManager) Deinit() {
	sfpMgr.logger.Info("SFP Manager Deinit()")
}

func (sfpMgr *SfpManager) GetSfpState(spfID int32) (*objects.SfpState, error) {
	var obj objects.SfpState
	return &obj, nil
}

func (sfpMgr *SfpManager) GetBulkSfpState(fromIdx, count int) (*objects.SfpStateGetInfo, error) {
	var obj objects.SfpStateGetInfo
	return &obj, nil
}

func (sfpMgr *SfpManager) GetSfpConfig(spfID int32) (*objects.SfpConfig, error) {
	var obj objects.SfpConfig
	return &obj, nil
}

func (sfpMgr *SfpManager) GetBulkSfpConfig(fromIdx, count int) (*objects.SfpConfigGetInfo, error) {
	var obj objects.SfpConfigGetInfo
	return &obj, nil
}

func (sfpMgr *SfpManager) UpdateSfpConfig(oldCfg *objects.SfpConfig, newCfg *objects.SfpConfig, attrset []bool) (bool, error) {
	return false, nil
}
