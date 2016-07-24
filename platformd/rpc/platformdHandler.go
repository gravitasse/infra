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

package rpc

import (
	//	"errors"
	"fmt"
	"infra/platformd/api"
	"platformd"
)

func (rpcHdl *rpcServiceHandler) GetPlatformSystemState(objName string) (*platformd.PlatformSystemState, error) {
	var rpcObj *platformd.PlatformSystemState

	obj, err := api.GetPlatformSystemState(objName)
	if err == nil {
		rpcObj = convertToRPCFmtPlatformSystemState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkPlatformSystemState(fromIdx, count platformd.Int) (*platformd.PlatformSystemStateGetInfo, error) {
	var getBulkObj platformd.PlatformSystemStateGetInfo

	info, err := api.GetBulkPlatformSystemState(int(fromIdx), int(count))
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.PlatformSystemStateList = append(getBulkObj.PlatformSystemStateList, convertToRPCFmtPlatformSystemState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetFanState(fanId int32) (*platformd.FanState, error) {
	var rpcObj *platformd.FanState
	var err error

	obj, err := api.GetFanState(fanId)
	if err == nil {
		rpcObj = convertToRPCFmtFanState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkFanState(fromIdx, count platformd.Int) (*platformd.FanStateGetInfo, error) {
	var getBulkObj platformd.FanStateGetInfo
	var err error

	fmt.Println("=====Inside GetBulkFanState====")
	info, err := api.GetBulkFanState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.FanStateList = append(getBulkObj.FanStateList, convertToRPCFmtFanState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) CreateFan(config *platformd.Fan) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) DeleteFan(config *platformd.Fan) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) UpdateFan(oldConfig *platformd.Fan, newConfig *platformd.Fan, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtFanConfig(oldConfig)
	newCfg := convertRPCToObjFmtFanConfig(newConfig)
	rv, err := api.UpdateFan(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetFan(fanId int32) (*platformd.Fan, error) {
	var rpcObj *platformd.Fan
	var err error

	obj, err := api.GetFanConfig(fanId)
	if err == nil {
		rpcObj = convertToRPCFmtFanConfig(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkFan(fromIdx, count platformd.Int) (*platformd.FanGetInfo, error) {
	var getBulkObj platformd.FanGetInfo
	var err error

	fmt.Println("=====Inside GetBulkFanConfig====")
	info, err := api.GetBulkFanConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.FanList = append(getBulkObj.FanList, convertToRPCFmtFanConfig(info.List[idx]))
	}
	return &getBulkObj, err
}
