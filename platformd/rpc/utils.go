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

func convertToRPCFmtPlatformState(obj *objects.PlatformState) *platformd.PlatformState {
	return &platformd.PlatformState{
		ObjName:      "Platform",
		ProductName:  obj.ProductName,
		SerialNum:    obj.SerialNum,
		Manufacturer: obj.Manufacturer,
		Vendor:       obj.Vendor,
		Release:      obj.Release,
		PlatformName: obj.PlatformName,
		Version:      obj.Version,
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
		LedId:         obj.LedId,
	}
}

func convertToRPCFmtFanConfig(obj *objects.FanConfig) *platformd.Fan {
	return &platformd.Fan{
		FanId:      obj.FanId,
		AdminSpeed: obj.AdminSpeed,
		AdminState: obj.AdminState,
	}
}

func convertRPCToObjFmtFanConfig(rpcObj *platformd.Fan) *objects.FanConfig {
	return &objects.FanConfig{
		FanId:      rpcObj.FanId,
		AdminSpeed: rpcObj.AdminSpeed,
		AdminState: rpcObj.AdminState,
	}
}

func convertToRPCFmtSfpConfig(obj *objects.SfpConfig) *platformd.Sfp {
	return &platformd.Sfp{
		SfpId: obj.SfpId,
	}
}

func convertToRPCFmtThermalState(obj *objects.ThermalState) *platformd.ThermalState {
	return &platformd.ThermalState{
		ThermalId:                 obj.ThermalId,
		Location:                  obj.Location,
		Temperature:               obj.Temperature,
		LowerWatermarkTemperature: obj.LowerWatermarkTemperature,
		UpperWatermarkTemperature: obj.UpperWatermarkTemperature,
		ShutdownTemperature:       obj.ShutdownTemperature,
	}
}
