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
	"infra/platformd/objects"
	"platformd"
)

func convertToRPCFmtPlatformSystemState(obj *objects.PlatformSystemState) *platformd.PlatformSystemState {
	return &platformd.PlatformSystemState{
		ObjName:   "PlatformSystemState",
		SerialNum: obj.SerialNum,
	}
}

func convertToRPCFmtFanState(obj *objects.FanState) *platformd.FanState {
	return &platformd.FanState{
		FanId:         obj.FanId,
		OperMode:      obj.OperMode,
		OperSpeed:     obj.OperSpeed,
		OperDirection: obj.OperDirection,
		Status:        obj.Status,
		Model:         obj.Model,
		SerialNum:     obj.SerialNum,
	}
}

func convertToRPCFmtFanConfig(obj *objects.FanConfig) *platformd.Fan {
	return &platformd.Fan{
		FanId:          obj.FanId,
		AdminSpeed:     obj.AdminSpeed,
		AdminDirection: obj.AdminDirection,
	}
}

func convertRPCToObjFmtFanConfig(rpcObj *platformd.Fan) *objects.FanConfig {
	return &objects.FanConfig{
		FanId:          rpcObj.FanId,
		AdminSpeed:     rpcObj.AdminSpeed,
		AdminDirection: rpcObj.AdminDirection,
	}
}
