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
	"errors"
	//"fmt"
	"infra/platformd/api"
	"platformd"
)

func (rpcHdl *rpcServiceHandler) GetPlatformState(objName string) (*platformd.PlatformState, error) {
	var rpcObj *platformd.PlatformState

	obj, err := api.GetPlatformState(objName)
	if err == nil {
		rpcObj = convertToRPCFmtPlatformState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkPlatformState(fromIdx, count platformd.Int) (*platformd.PlatformStateGetInfo, error) {
	var getBulkObj platformd.PlatformStateGetInfo

	info, err := api.GetBulkPlatformState(int(fromIdx), int(count))
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.PlatformStateList = append(getBulkObj.PlatformStateList, convertToRPCFmtPlatformState(info.List[idx]))
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

func (rpcHdl *rpcServiceHandler) GetBulkSfp(fromIdx, count platformd.Int) (*platformd.SfpGetInfo, error) {
	var getBulkObj platformd.SfpGetInfo
	var err error

	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) CreateSfp(config *platformd.Sfp) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) DeleteSfp(config *platformd.Sfp) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) UpdateSfp(oldConfig *platformd.Sfp, newConfig *platformd.Sfp, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	return false, nil
}

func (rpcHdl *rpcServiceHandler) GetSfp(sfpID int32) (*platformd.Sfp, error) {
	var rpcObj platformd.Sfp

	return &rpcObj, nil
}

func (rpcHdl *rpcServiceHandler) GetSfpState(sfpId int32) (*platformd.SfpState, error) {
	var rpcObj *platformd.SfpState

	obj, err := api.GetSfpState(sfpId)
	if err == nil {
		rpcObj = convertToRPCFmtSfpState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkSfpState(fromIdx, count platformd.Int) (*platformd.SfpStateGetInfo, error) {
	var getBulkObj platformd.SfpStateGetInfo
	var err error

	info, err := api.GetBulkSfpState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.SfpStateList = append(getBulkObj.SfpStateList, convertToRPCFmtSfpState(info.List[idx]))
	}
	return &getBulkObj, err
}

// TODO
func (rpcHdl *rpcServiceHandler) CreateLed(config *platformd.Led) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) DeleteLed(config *platformd.Led) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) UpdateLed(oldConfig *platformd.Led, newConfig *platformd.Led, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	return false, nil
}

func (rpcHdl *rpcServiceHandler) GetBulkLed(fromIdx, count platformd.Int) (*platformd.LedGetInfo, error) {
	var getBulkObj platformd.LedGetInfo
	var err error

	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetLed(ledID int32) (*platformd.Led, error) {
	var rpcObj platformd.Led

	return &rpcObj, nil
}

func (rpcHdl *rpcServiceHandler) GetBulkLedState(fromIdx, count platformd.Int) (*platformd.LedStateGetInfo, error) {
	var obj platformd.LedStateGetInfo

	return &obj, nil
}

func (rpcHdl *rpcServiceHandler) GetLedState(sfpId int32) (*platformd.LedState, error) {
	var obj platformd.LedState

	return &obj, nil
}

func (rpcHdl *rpcServiceHandler) GetBulkThermalState(fromIdx, count platformd.Int) (*platformd.ThermalStateGetInfo, error) {
	var getBulkObj platformd.ThermalStateGetInfo
	var err error

	info, err := api.GetBulkThermalState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.ThermalStateList = append(getBulkObj.ThermalStateList, convertToRPCFmtThermalState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetThermalState(thermalId int32) (*platformd.ThermalState, error) {
	var rpcObj *platformd.ThermalState
	var err error

	obj, err := api.GetThermalState(thermalId)
	if err == nil {
		rpcObj = convertToRPCFmtThermalState(obj)
	}
	return rpcObj, err

}

func (rpcHdl *rpcServiceHandler) CreatePsu(config *platformd.Psu) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) DeletePsu(config *platformd.Psu) (bool, error) {
	return true, nil
}

func (rpcHdl *rpcServiceHandler) UpdatePsu(oldConfig *platformd.Psu, newConfig *platformd.Psu, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	return false, nil
}

func (rpcHdl *rpcServiceHandler) GetBulkPsu(fromIdx, count platformd.Int) (*platformd.PsuGetInfo, error) {
	var getBulkObj platformd.PsuGetInfo
	var err error

	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetPsu(psuID int32) (*platformd.Psu, error) {
	var rpcObj platformd.Psu

	return &rpcObj, nil
}

func (rpcHdl *rpcServiceHandler) GetBulkPsuState(fromIdx, count platformd.Int) (*platformd.PsuStateGetInfo, error) {
	var obj platformd.PsuStateGetInfo

	return &obj, nil
}

func (rpcHdl *rpcServiceHandler) GetPsuState(PsuId int32) (*platformd.PsuState, error) {
	var obj platformd.PsuState

	return &obj, nil
}

func (rpcHdl *rpcServiceHandler) CreateFanSensor(config *platformd.FanSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) DeleteFanSensor(config *platformd.FanSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) UpdateFanSensor(oldConfig *platformd.FanSensor, newConfig *platformd.FanSensor, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtFanSensorConfig(oldConfig)
	newCfg := convertRPCToObjFmtFanSensorConfig(newConfig)
	rv, err := api.UpdateFanSensor(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetFanSensor(Name string) (*platformd.FanSensor, error) {
	return nil, errors.New("Not Supported")
}

func (rpcHdl *rpcServiceHandler) GetBulkFanSensor(fromIdx, count platformd.Int) (*platformd.FanSensorGetInfo, error) {
	var getBulkObj platformd.FanSensorGetInfo
	var err error

	info, err := api.GetBulkFanSensorConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.FanSensorList = append(getBulkObj.FanSensorList, convertToRPCFmtFanSensorConfig(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetFanSensorState(Name string) (*platformd.FanSensorState, error) {
	var rpcObj *platformd.FanSensorState
	var err error

	obj, err := api.GetFanSensorState(Name)
	if err == nil {
		rpcObj = convertToRPCFmtFanSensorState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkFanSensorState(fromIdx, count platformd.Int) (*platformd.FanSensorStateGetInfo, error) {
	var getBulkObj platformd.FanSensorStateGetInfo
	var err error

	info, err := api.GetBulkFanSensorState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.FanSensorStateList = append(getBulkObj.FanSensorStateList, convertToRPCFmtFanSensorState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetFanSensorPMDataState(Name string, Class string) (*platformd.FanSensorPMDataState, error) {
	var rpcObj *platformd.FanSensorPMDataState
	var err error

	obj, err := api.GetFanSensorPMDataState(Name, Class)
	if err == nil {
		rpcObj = convertToRPCFmtFanSensorPMState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkFanSensorPMDataState(fromIdx, count platformd.Int) (*platformd.FanSensorPMDataStateGetInfo, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) CreateTemperatureSensor(config *platformd.TemperatureSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) DeleteTemperatureSensor(config *platformd.TemperatureSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) UpdateTemperatureSensor(oldConfig *platformd.TemperatureSensor, newConfig *platformd.TemperatureSensor, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtTemperatureSensorConfig(oldConfig)
	newCfg := convertRPCToObjFmtTemperatureSensorConfig(newConfig)
	rv, err := api.UpdateTemperatureSensor(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetTemperatureSensor(Name string) (*platformd.TemperatureSensor, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) GetBulkTemperatureSensor(fromIdx, count platformd.Int) (*platformd.TemperatureSensorGetInfo, error) {
	var getBulkObj platformd.TemperatureSensorGetInfo
	var err error

	info, err := api.GetBulkTemperatureSensorConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.TemperatureSensorList = append(getBulkObj.TemperatureSensorList, convertToRPCFmtTemperatureSensorConfig(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetTemperatureSensorState(Name string) (*platformd.TemperatureSensorState, error) {
	var rpcObj *platformd.TemperatureSensorState
	var err error

	obj, err := api.GetTemperatureSensorState(Name)
	if err == nil {
		rpcObj = convertToRPCFmtTemperatureSensorState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkTemperatureSensorState(fromIdx, count platformd.Int) (*platformd.TemperatureSensorStateGetInfo, error) {
	var getBulkObj platformd.TemperatureSensorStateGetInfo
	var err error

	info, err := api.GetBulkTemperatureSensorState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.TemperatureSensorStateList = append(getBulkObj.TemperatureSensorStateList, convertToRPCFmtTemperatureSensorState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetTemperatureSensorPMDataState(Name string, Class string) (*platformd.TemperatureSensorPMDataState, error) {
	var rpcObj *platformd.TemperatureSensorPMDataState
	var err error

	obj, err := api.GetTempSensorPMDataState(Name, Class)
	if err == nil {
		rpcObj = convertToRPCFmtTempSensorPMState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkTemperatureSensorPMDataState(fromIdx, count platformd.Int) (*platformd.TemperatureSensorPMDataStateGetInfo, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) CreateVoltageSensor(config *platformd.VoltageSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) DeleteVoltageSensor(config *platformd.VoltageSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) UpdateVoltageSensor(oldConfig *platformd.VoltageSensor, newConfig *platformd.VoltageSensor, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtVoltageSensorConfig(oldConfig)
	newCfg := convertRPCToObjFmtVoltageSensorConfig(newConfig)
	rv, err := api.UpdateVoltageSensor(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetVoltageSensor(Name string) (*platformd.VoltageSensor, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) GetBulkVoltageSensor(fromIdx, count platformd.Int) (*platformd.VoltageSensorGetInfo, error) {
	var getBulkObj platformd.VoltageSensorGetInfo
	var err error

	info, err := api.GetBulkVoltageSensorConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.VoltageSensorList = append(getBulkObj.VoltageSensorList, convertToRPCFmtVoltageSensorConfig(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetVoltageSensorState(Name string) (*platformd.VoltageSensorState, error) {
	var rpcObj *platformd.VoltageSensorState
	var err error

	obj, err := api.GetVoltageSensorState(Name)
	if err == nil {
		rpcObj = convertToRPCFmtVoltageSensorState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkVoltageSensorState(fromIdx, count platformd.Int) (*platformd.VoltageSensorStateGetInfo, error) {
	var getBulkObj platformd.VoltageSensorStateGetInfo
	var err error

	info, err := api.GetBulkVoltageSensorState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.VoltageSensorStateList = append(getBulkObj.VoltageSensorStateList, convertToRPCFmtVoltageSensorState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetVoltageSensorPMDataState(Name string, Class string) (*platformd.VoltageSensorPMDataState, error) {
	var rpcObj *platformd.VoltageSensorPMDataState
	var err error

	obj, err := api.GetVoltageSensorPMDataState(Name, Class)
	if err == nil {
		rpcObj = convertToRPCFmtVoltageSensorPMState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkVoltageSensorPMDataState(fromIdx, count platformd.Int) (*platformd.VoltageSensorPMDataStateGetInfo, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) CreatePowerConverterSensor(config *platformd.PowerConverterSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) DeletePowerConverterSensor(config *platformd.PowerConverterSensor) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) UpdatePowerConverterSensor(oldConfig *platformd.PowerConverterSensor, newConfig *platformd.PowerConverterSensor, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtPowerConverterSensorConfig(oldConfig)
	newCfg := convertRPCToObjFmtPowerConverterSensorConfig(newConfig)
	rv, err := api.UpdatePowerConverterSensor(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetPowerConverterSensor(Name string) (*platformd.PowerConverterSensor, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) GetBulkPowerConverterSensor(fromIdx, count platformd.Int) (*platformd.PowerConverterSensorGetInfo, error) {
	var getBulkObj platformd.PowerConverterSensorGetInfo
	var err error

	info, err := api.GetBulkPowerConverterSensorConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.PowerConverterSensorList = append(getBulkObj.PowerConverterSensorList, convertToRPCFmtPowerConverterSensorConfig(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetPowerConverterSensorState(Name string) (*platformd.PowerConverterSensorState, error) {
	var rpcObj *platformd.PowerConverterSensorState
	var err error

	obj, err := api.GetPowerConverterSensorState(Name)
	if err == nil {
		rpcObj = convertToRPCFmtPowerConverterSensorState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkPowerConverterSensorState(fromIdx, count platformd.Int) (*platformd.PowerConverterSensorStateGetInfo, error) {
	var getBulkObj platformd.PowerConverterSensorStateGetInfo
	var err error

	info, err := api.GetBulkPowerConverterSensorState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.PowerConverterSensorStateList = append(getBulkObj.PowerConverterSensorStateList, convertToRPCFmtPowerConverterSensorState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetPowerConverterSensorPMDataState(Name string, Class string) (*platformd.PowerConverterSensorPMDataState, error) {
	var rpcObj *platformd.PowerConverterSensorPMDataState
	var err error

	obj, err := api.GetPowerConverterSensorPMDataState(Name, Class)
	if err == nil {
		rpcObj = convertToRPCFmtPowerConverterSensorPMState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkPowerConverterSensorPMDataState(fromIdx, count platformd.Int) (*platformd.PowerConverterSensorPMDataStateGetInfo, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) CreateQsfp(config *platformd.Qsfp) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) DeleteQsfp(config *platformd.Qsfp) (bool, error) {
	return false, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) UpdateQsfp(oldConfig *platformd.Qsfp, newConfig *platformd.Qsfp, attrset []bool, op []*platformd.PatchOpInfo) (bool, error) {
	oldCfg := convertRPCToObjFmtQsfpConfig(oldConfig)
	newCfg := convertRPCToObjFmtQsfpConfig(newConfig)
	rv, err := api.UpdateQsfp(oldCfg, newCfg, attrset)
	return rv, err
}

func (rpcHdl *rpcServiceHandler) GetQsfp(QsfpId int32) (*platformd.Qsfp, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) GetBulkQsfp(fromIdx, count platformd.Int) (*platformd.QsfpGetInfo, error) {
	var getBulkObj platformd.QsfpGetInfo
	var err error

	info, err := api.GetBulkQsfpConfig(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.QsfpList = append(getBulkObj.QsfpList, convertToRPCFmtQsfpConfig(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetQsfpState(QsfpId int32) (*platformd.QsfpState, error) {
	var rpcObj *platformd.QsfpState
	var err error

	obj, err := api.GetQsfpState(QsfpId)
	if err == nil {
		rpcObj = convertToRPCFmtQsfpState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkQsfpState(fromIdx, count platformd.Int) (*platformd.QsfpStateGetInfo, error) {
	var getBulkObj platformd.QsfpStateGetInfo
	var err error

	info, err := api.GetBulkQsfpState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.QsfpStateList = append(getBulkObj.QsfpStateList, convertToRPCFmtQsfpState(info.List[idx]))
	}
	return &getBulkObj, err
}

func (rpcHdl *rpcServiceHandler) GetQsfpPMDataState(QspfId int32, Resource string, Class string) (*platformd.QsfpPMDataState, error) {
	return nil, nil
}

func (rpcHdl *rpcServiceHandler) GetBulkQsfpPMDataState(fromIdx, count platformd.Int) (*platformd.QsfpPMDataStateGetInfo, error) {
	return nil, errors.New("Not supported")
}

func (rpcHdl *rpcServiceHandler) GetPlatformMgmtDeviceState(DeviceName string) (*platformd.PlatformMgmtDeviceState, error) {
	var rpcObj *platformd.PlatformMgmtDeviceState
	var err error

	obj, err := api.GetPlatformMgmtDeviceState(DeviceName)
	if err == nil {
		rpcObj = convertToRPCFmtPlatformMgmtDeviceState(obj)
	}
	return rpcObj, err
}

func (rpcHdl *rpcServiceHandler) GetBulkPlatformMgmtDeviceState(fromIdx, count platformd.Int) (*platformd.PlatformMgmtDeviceStateGetInfo, error) {
	var getBulkObj platformd.PlatformMgmtDeviceStateGetInfo
	var err error

	info, err := api.GetBulkPlatformMgmtDeviceState(int(fromIdx), int(count))
	if err != nil {
		return nil, err
	}
	getBulkObj.StartIdx = fromIdx
	getBulkObj.EndIdx = platformd.Int(info.EndIdx)
	getBulkObj.More = info.More
	getBulkObj.Count = platformd.Int(len(info.List))
	for idx := 0; idx < len(info.List); idx++ {
		getBulkObj.PlatformMgmtDeviceStateList = append(getBulkObj.PlatformMgmtDeviceStateList, convertToRPCFmtPlatformMgmtDeviceState(info.List[idx]))
	}
	return &getBulkObj, err
}
