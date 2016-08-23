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

package api

import (
	"errors"
	"infra/platformd/objects"
	"infra/platformd/server"
)

var svr *server.PlatformdServer

//var logger *logging.LogginIntf

//Initialize server handle
func InitApiLayer(server *server.PlatformdServer) {
	svr = server
	svr.Logger.Info("Initializing API layer")
}

func GetPlatformState(objName string) (*objects.PlatformState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_PLATFORM_STATE,
		Data: interface{}(&server.GetPlatformStateInArgs{
			ObjName: objName,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetPlatformStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetPlatformState")
	}
}

func GetBulkPlatformState(fromIdx, count int) (*objects.PlatformStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_PLATFORM_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkPlatformStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkPlatformState")
	}
}

func GetFanState(fanId int32) (*objects.FanState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_FAN_STATE,
		Data: interface{}(&server.GetFanStateInArgs{
			FanId: fanId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetFanStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetFanState")
	}
}

func GetBulkFanState(fromIdx, count int) (*objects.FanStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_FAN_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkFanStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkFanState")
	}
}

func UpdateFan(oldCfg *objects.FanConfig, newCfg *objects.FanConfig, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_FAN_CONFIG,
		Data: interface{}(&server.UpdateFanConfigInArgs{
			FanOldCfg: oldCfg,
			FanNewCfg: newCfg,
			AttrSet:   attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdateFan")
}

func GetFanConfig(fanId int32) (*objects.FanConfig, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_FAN_CONFIG,
		Data: interface{}(&server.GetFanConfigInArgs{
			FanId: fanId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetFanConfigOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetFanConfig")
	}
}

func GetBulkFanConfig(fromIdx, count int) (*objects.FanConfigGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_FAN_CONFIG,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkFanConfigOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkFanConfig")
	}
}

func UpdateSfp(oldCfg *objects.SfpConfig, newCfg *objects.SfpConfig, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_SFP_CONFIG,
		Data: interface{}(&server.UpdateSfpConfigInArgs{
			SfpOldCfg: oldCfg,
			SfpNewCfg: newCfg,
			AttrSet:   attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdateFan")
}

func GetSfpConfig(sfpId int32) (*objects.SfpConfig, error) {
	var obj objects.SfpConfig

	return &obj, nil
}

func GetBulkSfpConfig(fromIdx, count int) (*objects.SfpConfigGetInfo, error) {
	var obj objects.SfpConfigGetInfo

	return &obj, nil
}

func GetSfpState(sfpId int32) (*objects.SfpState, error) {
	var obj objects.SfpState

	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_SFP_STATE,
		Data: interface{}(&server.GetSfpStateInArgs{
			SfpId: sfpId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetSfpStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during SfpStateGet")
	}

	return &obj, nil
}

func GetBulkSfpState(fromIdx, count int) (*objects.SfpStateGetInfo, error) {
	var obj objects.SfpStateGetInfo

	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_SFP_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkSfpStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkSfpState")
	}

	return &obj, nil
}

func GetThermalState(thermalId int32) (*objects.ThermalState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_THERMAL_STATE,
		Data: interface{}(&server.GetThermalStateInArgs{
			ThermalId: thermalId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetThermalStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetThermalState")
	}
}

func GetBulkThermalState(fromIdx, count int) (*objects.ThermalStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_THERMAL_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkThermalStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkThermalState")
	}
}
