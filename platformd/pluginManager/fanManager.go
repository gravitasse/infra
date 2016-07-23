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
	"errors"
	"fmt"
	"infra/platformd/objects"
	"utils/logging"
)

type FanManager struct {
	logger logging.LoggerIntf
	plugin PluginIntf
}

var FanMgr FanManager

func (fMgr *FanManager) Init(logger logging.LoggerIntf, plugin PluginIntf) {
	fMgr.logger = logger
	fMgr.plugin = plugin
	fMgr.logger.Info("Fan Manager Init()")
}

func (fMgr *FanManager) Deinit() {
	fMgr.logger.Info("Fan Manager Deinit()")
}

func (fMgr *FanManager) GetFanState(fanId int32) (*objects.FanState, error) {
	if fMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	return fMgr.plugin.GetFanState(fanId)
}

func (fMgr *FanManager) GetBulkFanState(fromIdx int, count int) (*objects.FanStateGetInfo, error) {
	var retObj objects.FanStateGetInfo
	var fanId int32
	fanId = 1
	if fMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	obj, err := fMgr.plugin.GetFanState(fanId)
	if err != nil {
		fMgr.logger.Err(fmt.Sprintln("Error getting the fan state for fanId:", fanId))
	}
	retObj.List = append(retObj.List, obj)
	retObj.More = false
	retObj.EndIdx = 1
	retObj.Count = 1
	return &retObj, nil
}

func (fMgr *FanManager) GetFanConfig(fanId int32) (*objects.FanConfig, error) {
	if fMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	return fMgr.plugin.GetFanConfig(fanId)
}

func (fMgr *FanManager) GetBulkFanConfig(fromIdx int, count int) (*objects.FanConfigGetInfo, error) {
	var retObj objects.FanConfigGetInfo
	var fanId int32
	fanId = 1
	if fMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	obj, err := fMgr.plugin.GetFanConfig(fanId)
	if err != nil {
		fMgr.logger.Err(fmt.Sprintln("Error getting the fan config for fanId:", fanId))
	}
	retObj.List = append(retObj.List, obj)
	retObj.More = false
	retObj.EndIdx = 1
	retObj.Count = 1
	return &retObj, nil
}

func (fMgr *FanManager) UpdateFanConfig(oldCfg *objects.FanConfig, newCfg *objects.FanConfig, attrset []bool) (bool, error) {
	if fMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	ret, err := fMgr.plugin.UpdateFanConfig(newCfg)
	return ret, err
}
