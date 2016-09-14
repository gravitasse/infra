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
	"infra/platformd/pluginManager/pluginCommon"
	"models/events"
	"sync"
	"time"
	"utils/dbutils"
	"utils/eventUtils"
	"utils/logging"
	"utils/ringBuffer"
)

type FanSensorConfig struct {
	AdminState             string
	HigherAlarmThreshold   int32
	HigherWarningThreshold int32
	LowerWarningThreshold  int32
	LowerAlarmThreshold    int32
}

type TemperatureSensorConfig struct {
	AdminState             string
	HigherAlarmThreshold   float64
	HigherWarningThreshold float64
	LowerWarningThreshold  float64
	LowerAlarmThreshold    float64
}

type VoltageSensorConfig struct {
	AdminState             string
	HigherAlarmThreshold   float64
	HigherWarningThreshold float64
	LowerWarningThreshold  float64
	LowerAlarmThreshold    float64
}

type PowerConverterSensorConfig struct {
	AdminState             string
	HigherAlarmThreshold   float64
	HigherWarningThreshold float64
	LowerWarningThreshold  float64
	LowerAlarmThreshold    float64
}

type FanSensorEventData struct {
	Value int32
}

type TempSensorEventData struct {
	Value float64
}

type VoltageSensorEventData struct {
	Value float64
}

type PowerConverterSensorEventData struct {
	Value float64
}

const (
	classAInterval time.Duration = time.Duration(10) * time.Second // Polling Interval 10 sec
	classABufSize  int           = 6 * 60 * 24                     //Storage for 24 hrs
	classBInterval time.Duration = time.Duration(15) * time.Minute // Polling Interval 15 mins
	classBBufSize  int           = 4 * 24                          // Storage for 24 hrs
	classCInterval time.Duration = time.Duration(24) * time.Hour   // Polling Interval 24 Hrs
	classCBufSize  int           = 365                             // Storage for 365 days
)

type EventStatus struct {
	SentHigherAlarm bool
	SentHigherWarn  bool
	SentLowerWarn   bool
	SentLowerAlarm  bool
}

type SensorManager struct {
	logger                       logging.LoggerIntf
	plugin                       PluginIntf
	eventDbHdl                   dbutils.DBIntf
	classAPMTimer                *time.Timer
	classBPMTimer                *time.Timer
	classCPMTimer                *time.Timer
	fanSensorList                []string
	fanConfigMutex               sync.RWMutex
	fanConfigDB                  map[string]FanSensorConfig
	fanMsgStatus                 map[string]EventStatus
	fanClassAPMMutex             sync.RWMutex
	fanSensorClassAPM            map[string]*ringBuffer.RingBuffer
	fanClassBPMMutex             sync.RWMutex
	fanSensorClassBPM            map[string]*ringBuffer.RingBuffer
	fanClassCPMMutex             sync.RWMutex
	fanSensorClassCPM            map[string]*ringBuffer.RingBuffer
	tempSensorList               []string
	tempConfigMutex              sync.RWMutex
	tempConfigDB                 map[string]TemperatureSensorConfig
	tempMsgStatus                map[string]EventStatus
	tempClassAPMMutex            sync.RWMutex
	tempSensorClassAPM           map[string]*ringBuffer.RingBuffer
	tempClassBPMMutex            sync.RWMutex
	tempSensorClassBPM           map[string]*ringBuffer.RingBuffer
	tempClassCPMMutex            sync.RWMutex
	tempSensorClassCPM           map[string]*ringBuffer.RingBuffer
	voltageSensorList            []string
	voltageConfigMutex           sync.RWMutex
	voltageConfigDB              map[string]VoltageSensorConfig
	voltageMsgStatus             map[string]EventStatus
	voltageClassAPMMutex         sync.RWMutex
	voltageSensorClassAPM        map[string]*ringBuffer.RingBuffer
	voltageClassBPMMutex         sync.RWMutex
	voltageSensorClassBPM        map[string]*ringBuffer.RingBuffer
	voltageClassCPMMutex         sync.RWMutex
	voltageSensorClassCPM        map[string]*ringBuffer.RingBuffer
	powerConverterSensorList     []string
	powerConverterConfigMutex    sync.RWMutex
	powerConverterConfigDB       map[string]PowerConverterSensorConfig
	powerConverterMsgStatus      map[string]EventStatus
	powerConverterClassAPMMutex  sync.RWMutex
	powerConverterSensorClassAPM map[string]*ringBuffer.RingBuffer
	powerConverterClassBPMMutex  sync.RWMutex
	powerConverterSensorClassBPM map[string]*ringBuffer.RingBuffer
	powerConverterClassCPMMutex  sync.RWMutex
	powerConverterSensorClassCPM map[string]*ringBuffer.RingBuffer
	SensorStateMutex             sync.RWMutex
	SensorState                  *pluginCommon.SensorState
}

var SensorMgr SensorManager

func (sMgr *SensorManager) Init(logger logging.LoggerIntf, plugin PluginIntf, eventDbHdl dbutils.DBIntf) {
	var evtStatus EventStatus
	sMgr.logger = logger
	sMgr.eventDbHdl = eventDbHdl
	sMgr.logger.Info("sensor Manager Init( Start)")
	sMgr.plugin = plugin
	sMgr.SensorState = new(pluginCommon.SensorState)
	sMgr.SensorState.FanSensor = make(map[string]pluginCommon.FanSensorData)
	sMgr.SensorState.TemperatureSensor = make(map[string]pluginCommon.TemperatureSensorData)
	sMgr.SensorState.VoltageSensor = make(map[string]pluginCommon.VoltageSensorData)
	sMgr.SensorState.PowerConverterSensor = make(map[string]pluginCommon.PowerConverterSensorData)
	sMgr.fanSensorClassAPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.fanSensorClassBPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.fanSensorClassCPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.tempSensorClassAPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.tempSensorClassBPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.tempSensorClassCPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.voltageSensorClassAPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.voltageSensorClassBPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.voltageSensorClassCPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.powerConverterSensorClassAPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.powerConverterSensorClassBPM = make(map[string]*ringBuffer.RingBuffer)
	sMgr.powerConverterSensorClassCPM = make(map[string]*ringBuffer.RingBuffer)

	sMgr.SensorStateMutex.Lock()
	err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
	if err != nil {
		sMgr.logger.Info("Sensor Manager Init() Failed")
		return
	}
	sMgr.logger.Info("Sensor State:", sMgr.SensorState)
	sMgr.fanConfigDB = make(map[string]FanSensorConfig)
	sMgr.fanMsgStatus = make(map[string]EventStatus)
	for name, _ := range sMgr.SensorState.FanSensor {
		sMgr.logger.Info("Fan Sensor:", name)
		sMgr.fanConfigMutex.Lock()
		fanCfgEnt, _ := sMgr.fanConfigDB[name]
		// TODO: Read Json
		fanCfgEnt.AdminState = "Enabled"
		fanCfgEnt.HigherAlarmThreshold = 11000
		fanCfgEnt.HigherWarningThreshold = 11000
		fanCfgEnt.LowerAlarmThreshold = 1000
		fanCfgEnt.LowerWarningThreshold = 1000
		sMgr.fanConfigDB[name] = fanCfgEnt
		sMgr.fanMsgStatus[name] = evtStatus
		sMgr.fanConfigMutex.Unlock()
		sMgr.fanSensorList = append(sMgr.fanSensorList, name)
	}
	sMgr.tempConfigDB = make(map[string]TemperatureSensorConfig)
	sMgr.tempMsgStatus = make(map[string]EventStatus)
	for name, _ := range sMgr.SensorState.TemperatureSensor {
		sMgr.logger.Info("Temperature Sensor:", name)
		sMgr.tempConfigMutex.Lock()
		tempCfgEnt, _ := sMgr.tempConfigDB[name]
		// TODO: Read Json
		tempCfgEnt.AdminState = "Enabled"
		tempCfgEnt.HigherAlarmThreshold = 11000.0
		tempCfgEnt.HigherWarningThreshold = 11000.0
		tempCfgEnt.LowerAlarmThreshold = -1000.0
		tempCfgEnt.LowerWarningThreshold = -1000.0
		sMgr.tempConfigDB[name] = tempCfgEnt
		sMgr.tempMsgStatus[name] = evtStatus
		sMgr.tempConfigMutex.Unlock()
		sMgr.tempSensorList = append(sMgr.tempSensorList, name)
	}
	sMgr.voltageConfigDB = make(map[string]VoltageSensorConfig)
	sMgr.voltageMsgStatus = make(map[string]EventStatus)
	for name, _ := range sMgr.SensorState.VoltageSensor {
		sMgr.logger.Info("Voltage Sensor:", name)
		sMgr.voltageConfigMutex.Lock()
		voltageCfgEnt, _ := sMgr.voltageConfigDB[name]
		// TODO: Read Json
		voltageCfgEnt.AdminState = "Enabled"
		voltageCfgEnt.HigherAlarmThreshold = 11000
		voltageCfgEnt.HigherWarningThreshold = 11000
		voltageCfgEnt.LowerAlarmThreshold = 0
		voltageCfgEnt.LowerWarningThreshold = 0
		sMgr.voltageConfigDB[name] = voltageCfgEnt
		sMgr.voltageMsgStatus[name] = evtStatus
		sMgr.voltageConfigMutex.Unlock()
		sMgr.voltageSensorList = append(sMgr.voltageSensorList, name)
	}
	sMgr.powerConverterConfigDB = make(map[string]PowerConverterSensorConfig)
	sMgr.powerConverterMsgStatus = make(map[string]EventStatus)
	for name, _ := range sMgr.SensorState.PowerConverterSensor {
		sMgr.logger.Info("Power Sensor:", name)
		sMgr.powerConverterConfigMutex.Lock()
		powerConverterCfgEnt, _ := sMgr.powerConverterConfigDB[name]
		// TODO: Read Json
		powerConverterCfgEnt.AdminState = "Enabled"
		powerConverterCfgEnt.HigherAlarmThreshold = 11000
		powerConverterCfgEnt.HigherWarningThreshold = 11000
		powerConverterCfgEnt.LowerAlarmThreshold = 0
		powerConverterCfgEnt.LowerWarningThreshold = 0
		sMgr.powerConverterConfigDB[name] = powerConverterCfgEnt
		sMgr.powerConverterMsgStatus[name] = evtStatus
		sMgr.powerConverterConfigMutex.Unlock()
		sMgr.powerConverterSensorList = append(sMgr.powerConverterSensorList, name)
	}
	sMgr.SensorStateMutex.Unlock()
	sMgr.StartSensorPM()
	sMgr.logger.Info("sensor Manager Init( Done)")
}

func (sMgr *SensorManager) Deinit() {
	sMgr.logger.Info("Sensor Manager Deinit()")
}

func (sMgr *SensorManager) GetFanSensorState(Name string) (*objects.FanSensorState, error) {
	var fanSensorObj objects.FanSensorState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	_, exist := sMgr.fanConfigDB[Name]
	if !exist {
		return nil, errors.New("Invalid Fan Sensor Name")
	}

	sMgr.SensorStateMutex.Lock()
	err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
	if err != nil {
		sMgr.SensorStateMutex.Unlock()
		return nil, err
	}
	fanSensorState, _ := sMgr.SensorState.FanSensor[Name]
	fanSensorObj.Name = Name
	fanSensorObj.CurrentSpeed = fanSensorState.Value
	sMgr.SensorStateMutex.Unlock()
	return &fanSensorObj, err
}

func (sMgr *SensorManager) GetBulkFanSensorState(fromIdx int, cnt int) (*objects.FanSensorStateGetInfo, error) {
	var retObj objects.FanSensorStateGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.fanSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.fanSensorList) {
		retObj.EndIdx = len(sMgr.fanSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.fanSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		fanSensorName := sMgr.fanSensorList[idx]
		obj, err := sMgr.GetFanSensorState(fanSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the fan state for fan Sensor:", fanSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (sMgr *SensorManager) GetFanSensorConfig(Name string) (*objects.FanSensorConfig, error) {
	var fanSensorObj objects.FanSensorConfig
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	sMgr.fanConfigMutex.RLock()
	fanSensorCfgEnt, exist := sMgr.fanConfigDB[Name]
	if !exist {
		sMgr.fanConfigMutex.RUnlock()
		return nil, errors.New("Invalid Fan Sensor Name")
	}
	fanSensorObj.Name = Name
	fanSensorObj.AdminState = fanSensorCfgEnt.AdminState
	fanSensorObj.HigherAlarmThreshold = fanSensorCfgEnt.HigherAlarmThreshold
	fanSensorObj.HigherWarningThreshold = fanSensorCfgEnt.HigherWarningThreshold
	fanSensorObj.LowerAlarmThreshold = fanSensorCfgEnt.LowerAlarmThreshold
	fanSensorObj.LowerWarningThreshold = fanSensorCfgEnt.LowerWarningThreshold
	sMgr.fanConfigMutex.RUnlock()
	return &fanSensorObj, nil
}

func (sMgr *SensorManager) GetBulkFanSensorConfig(fromIdx int, cnt int) (*objects.FanSensorConfigGetInfo, error) {
	var retObj objects.FanSensorConfigGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.fanSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.fanSensorList) {
		retObj.EndIdx = len(sMgr.fanSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.fanSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		fanSensorName := sMgr.fanSensorList[idx]
		obj, err := sMgr.GetFanSensorConfig(fanSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the fan state for fan sensor:", fanSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func genFanSensorUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0
	for idx, val := range attrset {
		if true == val {
			switch idx {
			case 0:
				//ObjKey Fan Name
			case 1:
				mask |= objects.FAN_SENSOR_UPDATE_ADMIN_STATE
			case 2:
				mask |= objects.FAN_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD
			case 3:
				mask |= objects.FAN_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD
			case 4:
				mask |= objects.FAN_SENSOR_UPDATE_LOWER_WARN_THRESHOLD
			case 5:
				mask |= objects.FAN_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD
			}
		}
	}
	return mask
}

func (sMgr *SensorManager) UpdateFanSensorConfig(oldCfg *objects.FanSensorConfig, newCfg *objects.FanSensorConfig, attrset []bool) (bool, error) {
	if sMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	sMgr.fanConfigMutex.Lock()
	fanSensorCfgEnt, exist := sMgr.fanConfigDB[newCfg.Name]
	if !exist {
		sMgr.fanConfigMutex.Unlock()
		return false, errors.New("Invalid FanSensor Name")
	}
	mask := genFanSensorUpdateMask(attrset)
	if mask&objects.FAN_SENSOR_UPDATE_ADMIN_STATE == objects.FAN_SENSOR_UPDATE_ADMIN_STATE {
		fanSensorCfgEnt.AdminState = newCfg.AdminState
	}
	if mask&objects.FAN_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD == objects.FAN_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD {
		fanSensorCfgEnt.HigherAlarmThreshold = newCfg.HigherAlarmThreshold
	}
	if mask&objects.FAN_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD == objects.FAN_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD {
		fanSensorCfgEnt.HigherWarningThreshold = newCfg.HigherWarningThreshold
	}
	if mask&objects.FAN_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD == objects.FAN_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD {
		fanSensorCfgEnt.LowerAlarmThreshold = newCfg.LowerAlarmThreshold
	}
	if mask&objects.FAN_SENSOR_UPDATE_LOWER_WARN_THRESHOLD == objects.FAN_SENSOR_UPDATE_LOWER_WARN_THRESHOLD {
		fanSensorCfgEnt.LowerWarningThreshold = newCfg.LowerWarningThreshold
	}

	if !(fanSensorCfgEnt.HigherAlarmThreshold >= fanSensorCfgEnt.HigherWarningThreshold &&
		fanSensorCfgEnt.HigherWarningThreshold > fanSensorCfgEnt.LowerWarningThreshold &&
		fanSensorCfgEnt.LowerWarningThreshold >= fanSensorCfgEnt.LowerAlarmThreshold) {
		sMgr.fanConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, Please verify the thresholds")
	}
	sMgr.fanConfigDB[newCfg.Name] = fanSensorCfgEnt
	sMgr.fanConfigMutex.Unlock()
	return true, nil
}

func (sMgr *SensorManager) GetTemperatureSensorState(Name string) (*objects.TemperatureSensorState, error) {
	var tempSensorObj objects.TemperatureSensorState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	_, exist := sMgr.tempConfigDB[Name]
	if !exist {
		return nil, errors.New("Invalid Temperature Sensor Name")
	}

	sMgr.SensorStateMutex.Lock()
	err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
	if err != nil {
		sMgr.SensorStateMutex.Unlock()
		return nil, err
	}
	tempSensorState, _ := sMgr.SensorState.TemperatureSensor[Name]
	tempSensorObj.Name = Name
	tempSensorObj.CurrentTemperature = tempSensorState.Value
	sMgr.SensorStateMutex.Unlock()
	return &tempSensorObj, err
}

func (sMgr *SensorManager) GetBulkTemperatureSensorState(fromIdx int, cnt int) (*objects.TemperatureSensorStateGetInfo, error) {
	var retObj objects.TemperatureSensorStateGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.tempSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.tempSensorList) {
		retObj.EndIdx = len(sMgr.tempSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.tempSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		tempSensorName := sMgr.tempSensorList[idx]
		obj, err := sMgr.GetTemperatureSensorState(tempSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the temp state for temp Sensor:", tempSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (sMgr *SensorManager) GetTemperatureSensorConfig(Name string) (*objects.TemperatureSensorConfig, error) {
	var tempSensorObj objects.TemperatureSensorConfig
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	sMgr.tempConfigMutex.RLock()
	tempSensorCfgEnt, exist := sMgr.tempConfigDB[Name]
	if !exist {
		sMgr.tempConfigMutex.RUnlock()
		return nil, errors.New("Invalid Temperature Sensor Name")
	}
	tempSensorObj.Name = Name
	tempSensorObj.AdminState = tempSensorCfgEnt.AdminState
	tempSensorObj.HigherAlarmThreshold = tempSensorCfgEnt.HigherAlarmThreshold
	tempSensorObj.HigherWarningThreshold = tempSensorCfgEnt.HigherWarningThreshold
	tempSensorObj.LowerAlarmThreshold = tempSensorCfgEnt.LowerAlarmThreshold
	tempSensorObj.LowerWarningThreshold = tempSensorCfgEnt.LowerWarningThreshold
	sMgr.tempConfigMutex.RUnlock()
	return &tempSensorObj, nil
}

func (sMgr *SensorManager) GetBulkTemperatureSensorConfig(fromIdx int, cnt int) (*objects.TemperatureSensorConfigGetInfo, error) {
	var retObj objects.TemperatureSensorConfigGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.tempSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.tempSensorList) {
		retObj.EndIdx = len(sMgr.tempSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.tempSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		tempSensorName := sMgr.tempSensorList[idx]
		obj, err := sMgr.GetTemperatureSensorConfig(tempSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the temp state for temp sensor:", tempSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func genTempSensorUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0
	for idx, val := range attrset {
		if true == val {
			switch idx {
			case 0:
				//ObjKey Temp Sensor Name
			case 1:
				mask |= objects.TEMP_SENSOR_UPDATE_ADMIN_STATE
			case 2:
				mask |= objects.TEMP_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD
			case 3:
				mask |= objects.TEMP_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD
			case 4:
				mask |= objects.TEMP_SENSOR_UPDATE_LOWER_WARN_THRESHOLD
			case 5:
				mask |= objects.TEMP_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD
			}
		}
	}
	return mask
}

func (sMgr *SensorManager) UpdateTemperatureSensorConfig(oldCfg *objects.TemperatureSensorConfig, newCfg *objects.TemperatureSensorConfig, attrset []bool) (bool, error) {
	if sMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	sMgr.tempConfigMutex.Lock()
	tempSensorCfgEnt, exist := sMgr.tempConfigDB[newCfg.Name]
	if !exist {
		sMgr.tempConfigMutex.Unlock()
		return false, errors.New("Invalid TemperatureSensor Name")
	}
	mask := genTempSensorUpdateMask(attrset)
	if mask&objects.TEMP_SENSOR_UPDATE_ADMIN_STATE == objects.TEMP_SENSOR_UPDATE_ADMIN_STATE {
		tempSensorCfgEnt.AdminState = newCfg.AdminState
	}
	if mask&objects.TEMP_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD == objects.TEMP_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD {
		tempSensorCfgEnt.HigherAlarmThreshold = newCfg.HigherAlarmThreshold
	}
	if mask&objects.TEMP_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD == objects.TEMP_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD {
		tempSensorCfgEnt.HigherWarningThreshold = newCfg.HigherWarningThreshold
	}
	if mask&objects.TEMP_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD == objects.TEMP_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD {
		tempSensorCfgEnt.LowerAlarmThreshold = newCfg.LowerAlarmThreshold
	}
	if mask&objects.TEMP_SENSOR_UPDATE_LOWER_WARN_THRESHOLD == objects.TEMP_SENSOR_UPDATE_LOWER_WARN_THRESHOLD {
		tempSensorCfgEnt.LowerWarningThreshold = newCfg.LowerWarningThreshold
	}
	if !(tempSensorCfgEnt.HigherAlarmThreshold >= tempSensorCfgEnt.HigherWarningThreshold &&
		tempSensorCfgEnt.HigherWarningThreshold > tempSensorCfgEnt.LowerWarningThreshold &&
		tempSensorCfgEnt.LowerWarningThreshold >= tempSensorCfgEnt.LowerAlarmThreshold) {
		sMgr.tempConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, Please verify the thresholds")
	}
	sMgr.tempConfigDB[newCfg.Name] = tempSensorCfgEnt
	sMgr.tempConfigMutex.Unlock()
	return true, nil
}

func (sMgr *SensorManager) GetVoltageSensorState(Name string) (*objects.VoltageSensorState, error) {
	var voltageSensorObj objects.VoltageSensorState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	_, exist := sMgr.voltageConfigDB[Name]
	if !exist {
		return nil, errors.New("Invalid Voltage Sensor Name")
	}

	sMgr.SensorStateMutex.Lock()
	err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
	if err != nil {
		sMgr.SensorStateMutex.Unlock()
		return nil, err
	}
	voltageSensorState, _ := sMgr.SensorState.VoltageSensor[Name]
	voltageSensorObj.Name = Name
	voltageSensorObj.CurrentVoltage = voltageSensorState.Value
	sMgr.SensorStateMutex.Unlock()
	return &voltageSensorObj, err
}

func (sMgr *SensorManager) GetBulkVoltageSensorState(fromIdx int, cnt int) (*objects.VoltageSensorStateGetInfo, error) {
	var retObj objects.VoltageSensorStateGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.voltageSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.voltageSensorList) {
		retObj.EndIdx = len(sMgr.voltageSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.voltageSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		voltageSensorName := sMgr.voltageSensorList[idx]
		obj, err := sMgr.GetVoltageSensorState(voltageSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the voltage state for voltage Sensor:", voltageSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (sMgr *SensorManager) GetVoltageSensorConfig(Name string) (*objects.VoltageSensorConfig, error) {
	var voltageSensorObj objects.VoltageSensorConfig
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	sMgr.voltageConfigMutex.RLock()
	voltageSensorCfgEnt, exist := sMgr.voltageConfigDB[Name]
	if !exist {
		sMgr.voltageConfigMutex.RUnlock()
		return nil, errors.New("Invalid Voltage Sensor Name")
	}
	voltageSensorObj.Name = Name
	voltageSensorObj.AdminState = voltageSensorCfgEnt.AdminState
	voltageSensorObj.HigherAlarmThreshold = voltageSensorCfgEnt.HigherAlarmThreshold
	voltageSensorObj.HigherWarningThreshold = voltageSensorCfgEnt.HigherWarningThreshold
	voltageSensorObj.LowerAlarmThreshold = voltageSensorCfgEnt.LowerAlarmThreshold
	voltageSensorObj.LowerWarningThreshold = voltageSensorCfgEnt.LowerWarningThreshold
	sMgr.voltageConfigMutex.RUnlock()
	return &voltageSensorObj, nil
}

func (sMgr *SensorManager) GetBulkVoltageSensorConfig(fromIdx int, cnt int) (*objects.VoltageSensorConfigGetInfo, error) {
	var retObj objects.VoltageSensorConfigGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.voltageSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.voltageSensorList) {
		retObj.EndIdx = len(sMgr.voltageSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.voltageSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		voltageSensorName := sMgr.voltageSensorList[idx]
		obj, err := sMgr.GetVoltageSensorConfig(voltageSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the voltage state for voltage sensor:", voltageSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func genVoltageSensorUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0
	for idx, val := range attrset {
		if true == val {
			switch idx {
			case 0:
				//ObjKey Voltage Sensor Name
			case 1:
				mask |= objects.VOLTAGE_SENSOR_UPDATE_ADMIN_STATE
			case 2:
				mask |= objects.VOLTAGE_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD
			case 3:
				mask |= objects.VOLTAGE_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD
			case 4:
				mask |= objects.VOLTAGE_SENSOR_UPDATE_LOWER_WARN_THRESHOLD
			case 5:
				mask |= objects.VOLTAGE_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD
			}
		}
	}
	return mask
}

func (sMgr *SensorManager) UpdateVoltageSensorConfig(oldCfg *objects.VoltageSensorConfig, newCfg *objects.VoltageSensorConfig, attrset []bool) (bool, error) {
	if sMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	sMgr.voltageConfigMutex.Lock()
	voltageSensorCfgEnt, exist := sMgr.voltageConfigDB[newCfg.Name]
	if !exist {
		sMgr.voltageConfigMutex.Unlock()
		return false, errors.New("Invalid VoltageSensor Name")
	}
	mask := genVoltageSensorUpdateMask(attrset)
	if mask&objects.VOLTAGE_SENSOR_UPDATE_ADMIN_STATE == objects.VOLTAGE_SENSOR_UPDATE_ADMIN_STATE {
		voltageSensorCfgEnt.AdminState = newCfg.AdminState
	}
	if mask&objects.VOLTAGE_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD == objects.VOLTAGE_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD {
		voltageSensorCfgEnt.HigherAlarmThreshold = newCfg.HigherAlarmThreshold
	}
	if mask&objects.VOLTAGE_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD == objects.VOLTAGE_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD {
		voltageSensorCfgEnt.HigherWarningThreshold = newCfg.HigherWarningThreshold
	}
	if mask&objects.VOLTAGE_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD == objects.VOLTAGE_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD {
		voltageSensorCfgEnt.LowerAlarmThreshold = newCfg.LowerAlarmThreshold
	}
	if mask&objects.VOLTAGE_SENSOR_UPDATE_LOWER_WARN_THRESHOLD == objects.VOLTAGE_SENSOR_UPDATE_LOWER_WARN_THRESHOLD {
		voltageSensorCfgEnt.LowerWarningThreshold = newCfg.LowerWarningThreshold
	}
	if !(voltageSensorCfgEnt.HigherAlarmThreshold >= voltageSensorCfgEnt.HigherWarningThreshold &&
		voltageSensorCfgEnt.HigherWarningThreshold > voltageSensorCfgEnt.LowerWarningThreshold &&
		voltageSensorCfgEnt.LowerWarningThreshold >= voltageSensorCfgEnt.LowerAlarmThreshold) {
		sMgr.voltageConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, Please verify the thresholds")
	}
	sMgr.voltageConfigDB[newCfg.Name] = voltageSensorCfgEnt
	sMgr.voltageConfigMutex.Unlock()

	return true, nil
}

func (sMgr *SensorManager) GetPowerConverterSensorState(Name string) (*objects.PowerConverterSensorState, error) {
	var powerConverterSensorObj objects.PowerConverterSensorState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	_, exist := sMgr.powerConverterConfigDB[Name]
	if !exist {
		return nil, errors.New("Invalid PowerConverter Sensor Name")
	}

	sMgr.SensorStateMutex.Lock()
	err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
	if err != nil {
		sMgr.SensorStateMutex.Unlock()
		return nil, err
	}
	powerConverterSensorState, _ := sMgr.SensorState.PowerConverterSensor[Name]
	powerConverterSensorObj.Name = Name
	powerConverterSensorObj.CurrentPower = powerConverterSensorState.Value
	sMgr.SensorStateMutex.Unlock()
	return &powerConverterSensorObj, err
}

func (sMgr *SensorManager) GetBulkPowerConverterSensorState(fromIdx int, cnt int) (*objects.PowerConverterSensorStateGetInfo, error) {
	var retObj objects.PowerConverterSensorStateGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.powerConverterSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.powerConverterSensorList) {
		retObj.EndIdx = len(sMgr.powerConverterSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.powerConverterSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		powerConverterSensorName := sMgr.powerConverterSensorList[idx]
		obj, err := sMgr.GetPowerConverterSensorState(powerConverterSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the powerConverter state for powerConverter Sensor:", powerConverterSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (sMgr *SensorManager) GetPowerConverterSensorConfig(Name string) (*objects.PowerConverterSensorConfig, error) {
	var powerConverterSensorObj objects.PowerConverterSensorConfig
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	sMgr.powerConverterConfigMutex.RLock()
	powerConverterSensorCfgEnt, exist := sMgr.powerConverterConfigDB[Name]
	if !exist {
		sMgr.powerConverterConfigMutex.RUnlock()
		return nil, errors.New("Invalid PowerConverter Sensor Name")
	}
	powerConverterSensorObj.Name = Name
	powerConverterSensorObj.AdminState = powerConverterSensorCfgEnt.AdminState
	powerConverterSensorObj.HigherAlarmThreshold = powerConverterSensorCfgEnt.HigherAlarmThreshold
	powerConverterSensorObj.HigherWarningThreshold = powerConverterSensorCfgEnt.HigherWarningThreshold
	powerConverterSensorObj.LowerAlarmThreshold = powerConverterSensorCfgEnt.LowerAlarmThreshold
	powerConverterSensorObj.LowerWarningThreshold = powerConverterSensorCfgEnt.LowerWarningThreshold
	sMgr.powerConverterConfigMutex.RUnlock()
	return &powerConverterSensorObj, nil
}

func (sMgr *SensorManager) GetBulkPowerConverterSensorConfig(fromIdx int, cnt int) (*objects.PowerConverterSensorConfigGetInfo, error) {
	var retObj objects.PowerConverterSensorConfigGetInfo
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(sMgr.powerConverterSensorList) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(sMgr.powerConverterSensorList) {
		retObj.EndIdx = len(sMgr.powerConverterSensorList)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(sMgr.powerConverterSensorList) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		powerConverterSensorName := sMgr.powerConverterSensorList[idx]
		obj, err := sMgr.GetPowerConverterSensorConfig(powerConverterSensorName)
		if err != nil {
			sMgr.logger.Err(fmt.Sprintln("Error getting the powerConverter state for powerConverter sensor:", powerConverterSensorName))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func genPowerConverterSensorUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0
	for idx, val := range attrset {
		if true == val {
			switch idx {
			case 0:
				//ObjKey PowerConverter Sensor Name
			case 1:
				mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_ADMIN_STATE
			case 2:
				mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD
			case 3:
				mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD
			case 4:
				mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_WARN_THRESHOLD
			case 5:
				mask |= objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD
			}
		}
	}
	return mask
}

func (sMgr *SensorManager) UpdatePowerConverterSensorConfig(oldCfg *objects.PowerConverterSensorConfig, newCfg *objects.PowerConverterSensorConfig, attrset []bool) (bool, error) {
	if sMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	sMgr.powerConverterConfigMutex.Lock()
	powerConverterSensorCfgEnt, exist := sMgr.powerConverterConfigDB[newCfg.Name]
	if !exist {
		sMgr.powerConverterConfigMutex.Unlock()
		return false, errors.New("Invalid PowerConverterSensor Name")
	}
	mask := genPowerConverterSensorUpdateMask(attrset)
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_ADMIN_STATE == objects.POWER_CONVERTER_SENSOR_UPDATE_ADMIN_STATE {
		powerConverterSensorCfgEnt.AdminState = newCfg.AdminState
	}
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD == objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_ALARM_THRESHOLD {
		powerConverterSensorCfgEnt.HigherAlarmThreshold = newCfg.HigherAlarmThreshold
	}
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD == objects.POWER_CONVERTER_SENSOR_UPDATE_HIGHER_WARN_THRESHOLD {
		powerConverterSensorCfgEnt.HigherWarningThreshold = newCfg.HigherWarningThreshold
	}
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD == objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_ALARM_THRESHOLD {
		powerConverterSensorCfgEnt.LowerAlarmThreshold = newCfg.LowerAlarmThreshold
	}
	if mask&objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_WARN_THRESHOLD == objects.POWER_CONVERTER_SENSOR_UPDATE_LOWER_WARN_THRESHOLD {
		powerConverterSensorCfgEnt.LowerWarningThreshold = newCfg.LowerWarningThreshold
	}
	if !(powerConverterSensorCfgEnt.HigherAlarmThreshold >= powerConverterSensorCfgEnt.HigherWarningThreshold &&
		powerConverterSensorCfgEnt.HigherWarningThreshold > powerConverterSensorCfgEnt.LowerWarningThreshold &&
		powerConverterSensorCfgEnt.LowerWarningThreshold >= powerConverterSensorCfgEnt.LowerAlarmThreshold) {
		sMgr.powerConverterConfigMutex.Unlock()
		return false, errors.New("Invalid configuration, Please verify the thresholds")
	}
	sMgr.powerConverterConfigDB[newCfg.Name] = powerConverterSensorCfgEnt
	sMgr.powerConverterConfigMutex.Unlock()

	return true, nil
}

func (sMgr *SensorManager) InitSensorPM() {
	for fanSensorName, _ := range sMgr.fanConfigDB {
		sMgr.fanSensorClassAPM[fanSensorName] = new(ringBuffer.RingBuffer)
		sMgr.fanSensorClassAPM[fanSensorName].SetRingBufferCapacity(classABufSize)
		sMgr.fanSensorClassBPM[fanSensorName] = new(ringBuffer.RingBuffer)
		sMgr.fanSensorClassBPM[fanSensorName].SetRingBufferCapacity(classBBufSize)
		sMgr.fanSensorClassCPM[fanSensorName] = new(ringBuffer.RingBuffer)
		sMgr.fanSensorClassCPM[fanSensorName].SetRingBufferCapacity(classCBufSize)
	}
	for tempSensorName, _ := range sMgr.tempConfigDB {
		sMgr.tempSensorClassAPM[tempSensorName] = new(ringBuffer.RingBuffer)
		sMgr.tempSensorClassAPM[tempSensorName].SetRingBufferCapacity(classABufSize)
		sMgr.tempSensorClassBPM[tempSensorName] = new(ringBuffer.RingBuffer)
		sMgr.tempSensorClassBPM[tempSensorName].SetRingBufferCapacity(classBBufSize)
		sMgr.tempSensorClassCPM[tempSensorName] = new(ringBuffer.RingBuffer)
		sMgr.tempSensorClassCPM[tempSensorName].SetRingBufferCapacity(classCBufSize)
	}

	for voltageSensorName, _ := range sMgr.voltageConfigDB {
		sMgr.voltageSensorClassAPM[voltageSensorName] = new(ringBuffer.RingBuffer)
		sMgr.voltageSensorClassAPM[voltageSensorName].SetRingBufferCapacity(classABufSize)
		sMgr.voltageSensorClassBPM[voltageSensorName] = new(ringBuffer.RingBuffer)
		sMgr.voltageSensorClassBPM[voltageSensorName].SetRingBufferCapacity(classBBufSize)
		sMgr.voltageSensorClassCPM[voltageSensorName] = new(ringBuffer.RingBuffer)
		sMgr.voltageSensorClassCPM[voltageSensorName].SetRingBufferCapacity(classCBufSize)
	}
	for powerConverterSensorName, _ := range sMgr.powerConverterConfigDB {
		sMgr.powerConverterSensorClassAPM[powerConverterSensorName] = new(ringBuffer.RingBuffer)
		sMgr.powerConverterSensorClassAPM[powerConverterSensorName].SetRingBufferCapacity(classABufSize)
		sMgr.powerConverterSensorClassBPM[powerConverterSensorName] = new(ringBuffer.RingBuffer)
		sMgr.powerConverterSensorClassBPM[powerConverterSensorName].SetRingBufferCapacity(classBBufSize)
		sMgr.powerConverterSensorClassCPM[powerConverterSensorName] = new(ringBuffer.RingBuffer)
		sMgr.powerConverterSensorClassCPM[powerConverterSensorName].SetRingBufferCapacity(classCBufSize)
	}
}

func (sMgr *SensorManager) StartSensorPM() {
	sMgr.InitSensorPM()
	sMgr.StartSensorPMClass("Class-A")
	sMgr.StartSensorPMClass("Class-B")
	sMgr.StartSensorPMClass("Class-C")
}

func (sMgr *SensorManager) ProcessFanSensorPM(sensorState *pluginCommon.SensorState, class string) {
	sMgr.fanConfigMutex.Lock()
	for fanSensorName, fanSensorCfgEnt := range sMgr.fanConfigDB {
		fanSensorState, _ := sMgr.SensorState.FanSensor[fanSensorName]
		fanSensorPMData := objects.FanSensorPMData{
			TimeStamp: time.Now().String(),
			Value:     fanSensorState.Value,
		}
		if fanSensorCfgEnt.AdminState == "Enabled" {
			eventKey := events.FanSensorKey{
				Name: fanSensorName,
			}
			eventData := FanSensorEventData{
				Value: fanSensorState.Value,
			}
			txEvent := eventUtils.TxEvent{
				Key:            eventKey,
				AdditionalInfo: "",
				AdditionalData: eventData,
			}
			var evts []events.EventId
			var curEvents EventStatus
			prevEventStatus := sMgr.fanMsgStatus[fanSensorName]
			if fanSensorState.Value >= fanSensorCfgEnt.HigherAlarmThreshold {
				curEvents.SentHigherAlarm = true
			}
			if fanSensorState.Value < fanSensorCfgEnt.HigherAlarmThreshold &&
				fanSensorState.Value >= fanSensorCfgEnt.HigherWarningThreshold {
				curEvents.SentHigherWarn = true
			}
			if fanSensorState.Value <= fanSensorCfgEnt.LowerWarningThreshold &&
				fanSensorState.Value > fanSensorCfgEnt.LowerAlarmThreshold {
				curEvents.SentLowerWarn = true
			}
			if fanSensorState.Value <= fanSensorCfgEnt.LowerAlarmThreshold {
				curEvents.SentLowerAlarm = true
			}
			if prevEventStatus.SentHigherAlarm != curEvents.SentHigherAlarm {
				if curEvents.SentHigherAlarm == true {
					evts = append(evts, events.FanHigherTCAAlarm)
				} else {
					evts = append(evts, events.FanHigherTCAAlarmClear)
				}
			}
			if prevEventStatus.SentHigherWarn != curEvents.SentHigherWarn {
				if curEvents.SentHigherWarn == true {
					evts = append(evts, events.FanHigherTCAWarn)
				} else {
					evts = append(evts, events.FanHigherTCAWarnClear)
				}
			}
			if prevEventStatus.SentLowerAlarm != curEvents.SentLowerAlarm {
				if curEvents.SentLowerAlarm == true {
					evts = append(evts, events.FanLowerTCAAlarm)
				} else {
					evts = append(evts, events.FanLowerTCAAlarmClear)
				}
			}
			if prevEventStatus.SentLowerWarn != curEvents.SentLowerWarn {
				if curEvents.SentLowerWarn == true {
					evts = append(evts, events.FanLowerTCAWarn)
				} else {
					evts = append(evts, events.FanLowerTCAWarnClear)
				}
			}
			sMgr.fanMsgStatus[fanSensorName] = curEvents
			for _, evt := range evts {
				txEvent.EventId = evt
				err := eventUtils.PublishEvents(&txEvent)
				if err != nil {
					sMgr.logger.Err("Error publish events")
				}
			}

		}
		switch class {
		case "Class-A":
			sMgr.fanClassAPMMutex.Lock()
			sMgr.fanSensorClassAPM[fanSensorName].InsertIntoRingBuffer(fanSensorPMData)
			sMgr.fanClassAPMMutex.Unlock()
		case "Class-B":
			sMgr.fanClassBPMMutex.Lock()
			sMgr.fanSensorClassBPM[fanSensorName].InsertIntoRingBuffer(fanSensorPMData)
			sMgr.fanClassBPMMutex.Unlock()
		case "Class-C":
			sMgr.fanClassCPMMutex.Lock()
			sMgr.fanSensorClassCPM[fanSensorName].InsertIntoRingBuffer(fanSensorPMData)
			sMgr.fanClassCPMMutex.Unlock()
		}
	}
	sMgr.fanConfigMutex.Unlock()
}

func (sMgr *SensorManager) ProcessTempSensorPM(sensorState *pluginCommon.SensorState, class string) {
	sMgr.tempConfigMutex.Lock()
	for tempSensorName, tempSensorCfgEnt := range sMgr.tempConfigDB {
		tempSensorState, _ := sMgr.SensorState.TemperatureSensor[tempSensorName]
		tempSensorPMData := objects.TemperatureSensorPMData{
			TimeStamp: time.Now().String(),
			Value:     tempSensorState.Value,
		}
		if tempSensorCfgEnt.AdminState == "Enabled" {
			eventKey := events.TemperatureSensorKey{
				Name: tempSensorName,
			}
			eventData := TempSensorEventData{
				Value: tempSensorState.Value,
			}
			txEvent := eventUtils.TxEvent{
				Key:            eventKey,
				AdditionalInfo: "",
				AdditionalData: eventData,
			}
			var evts []events.EventId
			var curEvents EventStatus
			prevEventStatus := sMgr.tempMsgStatus[tempSensorName]
			if tempSensorState.Value >= tempSensorCfgEnt.HigherAlarmThreshold {
				curEvents.SentHigherAlarm = true
			}
			if tempSensorState.Value < tempSensorCfgEnt.HigherAlarmThreshold &&
				tempSensorState.Value >= tempSensorCfgEnt.HigherWarningThreshold {
				curEvents.SentHigherWarn = true
			}
			if tempSensorState.Value <= tempSensorCfgEnt.LowerWarningThreshold &&
				tempSensorState.Value > tempSensorCfgEnt.LowerAlarmThreshold {
				curEvents.SentLowerWarn = true
			}
			if tempSensorState.Value <= tempSensorCfgEnt.LowerAlarmThreshold {
				curEvents.SentLowerAlarm = true
			}
			if prevEventStatus.SentHigherAlarm != curEvents.SentHigherAlarm {
				if curEvents.SentHigherAlarm == true {
					evts = append(evts, events.TemperatureHigherTCAAlarm)
				} else {
					evts = append(evts, events.TemperatureHigherTCAAlarmClear)
				}
			}
			if prevEventStatus.SentHigherWarn != curEvents.SentHigherWarn {
				if curEvents.SentHigherWarn == true {
					evts = append(evts, events.TemperatureHigherTCAWarn)
				} else {
					evts = append(evts, events.TemperatureHigherTCAWarnClear)
				}
			}
			if prevEventStatus.SentLowerAlarm != curEvents.SentLowerAlarm {
				if curEvents.SentLowerAlarm == true {
					evts = append(evts, events.TemperatureLowerTCAAlarm)
				} else {
					evts = append(evts, events.TemperatureLowerTCAAlarmClear)
				}
			}
			if prevEventStatus.SentLowerWarn != curEvents.SentLowerWarn {
				if curEvents.SentLowerWarn == true {
					evts = append(evts, events.TemperatureLowerTCAWarn)
				} else {
					evts = append(evts, events.TemperatureLowerTCAWarnClear)
				}
			}
			sMgr.tempMsgStatus[tempSensorName] = curEvents
			for _, evt := range evts {
				txEvent.EventId = evt
				err := eventUtils.PublishEvents(&txEvent)
				if err != nil {
					sMgr.logger.Err("Error publish events")
				}
			}

		}

		switch class {
		case "Class-A":
			sMgr.tempClassAPMMutex.Lock()
			sMgr.tempSensorClassAPM[tempSensorName].InsertIntoRingBuffer(tempSensorPMData)
			sMgr.tempClassAPMMutex.Unlock()
		case "Class-B":
			sMgr.tempClassBPMMutex.Lock()
			sMgr.tempSensorClassBPM[tempSensorName].InsertIntoRingBuffer(tempSensorPMData)
			sMgr.tempClassBPMMutex.Unlock()
		case "Class-C":
			sMgr.tempClassCPMMutex.Lock()
			sMgr.tempSensorClassCPM[tempSensorName].InsertIntoRingBuffer(tempSensorPMData)
			sMgr.tempClassCPMMutex.Unlock()
		}
	}
	sMgr.tempConfigMutex.Unlock()
}

func (sMgr *SensorManager) ProcessVoltageSensorPM(sensorState *pluginCommon.SensorState, class string) {
	sMgr.voltageConfigMutex.Lock()
	for voltageSensorName, voltageSensorCfgEnt := range sMgr.voltageConfigDB {
		voltageSensorState, _ := sMgr.SensorState.VoltageSensor[voltageSensorName]
		voltageSensorPMData := objects.VoltageSensorPMData{
			TimeStamp: time.Now().String(),
			Value:     voltageSensorState.Value,
		}
		if voltageSensorCfgEnt.AdminState == "Enabled" {
			eventKey := events.VoltageSensorKey{
				Name: voltageSensorName,
			}
			eventData := VoltageSensorEventData{
				Value: voltageSensorState.Value,
			}
			txEvent := eventUtils.TxEvent{
				Key:            eventKey,
				AdditionalInfo: "",
				AdditionalData: eventData,
			}
			var evts []events.EventId
			var curEvents EventStatus
			prevEventStatus := sMgr.voltageMsgStatus[voltageSensorName]
			if voltageSensorState.Value >= voltageSensorCfgEnt.HigherAlarmThreshold {
				curEvents.SentHigherAlarm = true
			}
			if voltageSensorState.Value < voltageSensorCfgEnt.HigherAlarmThreshold &&
				voltageSensorState.Value >= voltageSensorCfgEnt.HigherWarningThreshold {
				curEvents.SentHigherWarn = true
			}
			if voltageSensorState.Value <= voltageSensorCfgEnt.LowerWarningThreshold &&
				voltageSensorState.Value > voltageSensorCfgEnt.LowerAlarmThreshold {
				curEvents.SentLowerWarn = true
			}
			if voltageSensorState.Value <= voltageSensorCfgEnt.LowerAlarmThreshold {
				curEvents.SentLowerAlarm = true
			}
			if prevEventStatus.SentHigherAlarm != curEvents.SentHigherAlarm {
				if curEvents.SentHigherAlarm == true {
					evts = append(evts, events.VoltageHigherTCAAlarm)
				} else {
					evts = append(evts, events.VoltageHigherTCAAlarmClear)
				}
			}
			if prevEventStatus.SentHigherWarn != curEvents.SentHigherWarn {
				if curEvents.SentHigherWarn == true {
					evts = append(evts, events.VoltageHigherTCAWarn)
				} else {
					evts = append(evts, events.VoltageHigherTCAWarnClear)
				}
			}
			if prevEventStatus.SentLowerAlarm != curEvents.SentLowerAlarm {
				if curEvents.SentLowerAlarm == true {
					evts = append(evts, events.VoltageLowerTCAAlarm)
				} else {
					evts = append(evts, events.VoltageLowerTCAAlarmClear)
				}
			}
			if prevEventStatus.SentLowerWarn != curEvents.SentLowerWarn {
				if curEvents.SentLowerWarn == true {
					evts = append(evts, events.VoltageLowerTCAWarn)
				} else {
					evts = append(evts, events.VoltageLowerTCAWarnClear)
				}
			}
			sMgr.voltageMsgStatus[voltageSensorName] = curEvents
			for _, evt := range evts {
				txEvent.EventId = evt
				err := eventUtils.PublishEvents(&txEvent)
				if err != nil {
					sMgr.logger.Err("Error publish events")
				}
			}

		}

		switch class {
		case "Class-A":
			sMgr.voltageClassAPMMutex.Lock()
			sMgr.voltageSensorClassAPM[voltageSensorName].InsertIntoRingBuffer(voltageSensorPMData)
			sMgr.voltageClassAPMMutex.Unlock()
		case "Class-B":
			sMgr.voltageClassBPMMutex.Lock()
			sMgr.voltageSensorClassBPM[voltageSensorName].InsertIntoRingBuffer(voltageSensorPMData)
			sMgr.voltageClassBPMMutex.Unlock()
		case "Class-C":
			sMgr.voltageClassCPMMutex.Lock()
			sMgr.voltageSensorClassCPM[voltageSensorName].InsertIntoRingBuffer(voltageSensorPMData)
			sMgr.voltageClassCPMMutex.Unlock()
		}
	}
	sMgr.voltageConfigMutex.Unlock()
}

func (sMgr *SensorManager) ProcessPowerConverterSensorPM(sensorState *pluginCommon.SensorState, class string) {
	sMgr.powerConverterConfigMutex.Lock()
	for powerConverterSensorName, powerConverterSensorCfgEnt := range sMgr.powerConverterConfigDB {
		powerConverterSensorState, _ := sMgr.SensorState.PowerConverterSensor[powerConverterSensorName]
		powerConverterSensorPMData := objects.PowerConverterSensorPMData{
			TimeStamp: time.Now().String(),
			Value:     powerConverterSensorState.Value,
		}
		if powerConverterSensorCfgEnt.AdminState == "Enabled" {
			eventKey := events.PowerConverterSensorKey{
				Name: powerConverterSensorName,
			}
			eventData := PowerConverterSensorEventData{
				Value: powerConverterSensorState.Value,
			}
			txEvent := eventUtils.TxEvent{
				Key:            eventKey,
				AdditionalInfo: "",
				AdditionalData: eventData,
			}
			var evts []events.EventId
			var curEvents EventStatus
			prevEventStatus := sMgr.powerConverterMsgStatus[powerConverterSensorName]
			if powerConverterSensorState.Value >= powerConverterSensorCfgEnt.HigherAlarmThreshold {
				curEvents.SentHigherAlarm = true
			}
			if powerConverterSensorState.Value < powerConverterSensorCfgEnt.HigherAlarmThreshold &&
				powerConverterSensorState.Value >= powerConverterSensorCfgEnt.HigherWarningThreshold {
				curEvents.SentHigherWarn = true
			}
			if powerConverterSensorState.Value <= powerConverterSensorCfgEnt.LowerWarningThreshold &&
				powerConverterSensorState.Value > powerConverterSensorCfgEnt.LowerAlarmThreshold {
				curEvents.SentLowerWarn = true
			}
			if powerConverterSensorState.Value <= powerConverterSensorCfgEnt.LowerAlarmThreshold {
				curEvents.SentLowerAlarm = true
			}
			if prevEventStatus.SentHigherAlarm != curEvents.SentHigherAlarm {
				if curEvents.SentHigherAlarm == true {
					evts = append(evts, events.PowerConverterHigherTCAAlarm)
				} else {
					evts = append(evts, events.PowerConverterHigherTCAAlarmClear)
				}
			}
			if prevEventStatus.SentHigherWarn != curEvents.SentHigherWarn {
				if curEvents.SentHigherWarn == true {
					evts = append(evts, events.PowerConverterHigherTCAWarn)
				} else {
					evts = append(evts, events.PowerConverterHigherTCAWarnClear)
				}
			}
			if prevEventStatus.SentLowerAlarm != curEvents.SentLowerAlarm {
				if curEvents.SentLowerAlarm == true {
					evts = append(evts, events.PowerConverterLowerTCAAlarm)
				} else {
					evts = append(evts, events.PowerConverterLowerTCAAlarmClear)
				}
			}
			if prevEventStatus.SentLowerWarn != curEvents.SentLowerWarn {
				if curEvents.SentLowerWarn == true {
					evts = append(evts, events.PowerConverterLowerTCAWarn)
				} else {
					evts = append(evts, events.PowerConverterLowerTCAWarnClear)
				}
			}
			sMgr.powerConverterMsgStatus[powerConverterSensorName] = curEvents
			for _, evt := range evts {
				txEvent.EventId = evt
				err := eventUtils.PublishEvents(&txEvent)
				if err != nil {
					sMgr.logger.Err("Error publish events")
				}
			}

		}

		switch class {
		case "Class-A":
			sMgr.powerConverterClassAPMMutex.Lock()
			sMgr.powerConverterSensorClassAPM[powerConverterSensorName].InsertIntoRingBuffer(powerConverterSensorPMData)
			sMgr.powerConverterClassAPMMutex.Unlock()
		case "Class-B":
			sMgr.powerConverterClassBPMMutex.Lock()
			sMgr.powerConverterSensorClassBPM[powerConverterSensorName].InsertIntoRingBuffer(powerConverterSensorPMData)
			sMgr.powerConverterClassBPMMutex.Unlock()
		case "Class-C":
			sMgr.powerConverterClassCPMMutex.Lock()
			sMgr.powerConverterSensorClassCPM[powerConverterSensorName].InsertIntoRingBuffer(powerConverterSensorPMData)
			sMgr.powerConverterClassCPMMutex.Unlock()
		}
	}
	sMgr.powerConverterConfigMutex.Unlock()
}

func (sMgr *SensorManager) StartSensorPMClass(class string) {
	sMgr.SensorStateMutex.Lock()
	err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
	if err != nil {
		sMgr.logger.Err("Error getting sensor data during start of PM")
		return
	}
	sMgr.ProcessFanSensorPM(sMgr.SensorState, class)
	sMgr.ProcessTempSensorPM(sMgr.SensorState, class)
	sMgr.ProcessVoltageSensorPM(sMgr.SensorState, class)
	sMgr.ProcessPowerConverterSensorPM(sMgr.SensorState, class)
	sMgr.SensorStateMutex.Unlock()

	switch class {
	case "Class-A":
		classAPMFunc := func() {
			sMgr.SensorStateMutex.Lock()
			err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
			if err != nil {
				sMgr.logger.Err("Error getting sensor data during PM processing")
				sMgr.SensorStateMutex.Unlock()
				return
			}
			sMgr.ProcessFanSensorPM(sMgr.SensorState, class)
			sMgr.ProcessTempSensorPM(sMgr.SensorState, class)
			sMgr.ProcessVoltageSensorPM(sMgr.SensorState, class)
			sMgr.ProcessPowerConverterSensorPM(sMgr.SensorState, class)
			sMgr.SensorStateMutex.Unlock()
			sMgr.classAPMTimer.Reset(classAInterval)
		}
		sMgr.classAPMTimer = time.AfterFunc(classAInterval, classAPMFunc)
	case "Class-B":
		classBPMFunc := func() {
			sMgr.SensorStateMutex.Lock()
			err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
			if err != nil {
				sMgr.logger.Err("Error getting sensor data during PM processing")
				sMgr.SensorStateMutex.Unlock()
				return
			}
			sMgr.ProcessFanSensorPM(sMgr.SensorState, class)
			sMgr.ProcessTempSensorPM(sMgr.SensorState, class)
			sMgr.ProcessVoltageSensorPM(sMgr.SensorState, class)
			sMgr.ProcessPowerConverterSensorPM(sMgr.SensorState, class)
			sMgr.SensorStateMutex.Unlock()
			sMgr.classBPMTimer.Reset(classBInterval)
		}
		sMgr.classBPMTimer = time.AfterFunc(classBInterval, classBPMFunc)
	case "Class-C":
		classCPMFunc := func() {
			sMgr.SensorStateMutex.Lock()
			err := sMgr.plugin.GetAllSensorState(sMgr.SensorState)
			if err != nil {
				sMgr.logger.Err("Error getting sensor data during PM processing")
				sMgr.SensorStateMutex.Unlock()
				return
			}
			sMgr.ProcessFanSensorPM(sMgr.SensorState, class)
			sMgr.ProcessTempSensorPM(sMgr.SensorState, class)
			sMgr.ProcessVoltageSensorPM(sMgr.SensorState, class)
			sMgr.ProcessPowerConverterSensorPM(sMgr.SensorState, class)
			sMgr.SensorStateMutex.Unlock()
			sMgr.classCPMTimer.Reset(classCInterval)
		}
		sMgr.classCPMTimer = time.AfterFunc(classCInterval, classCPMFunc)
	}

}

func (sMgr *SensorManager) GetFanSensorPMState(Name string, Class string) (*objects.FanSensorPMState, error) {
	var fanSensorPMObj objects.FanSensorPMState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	_, exist := sMgr.fanConfigDB[Name]
	if !exist {
		return nil, errors.New("Invalid Fan Sensor Name")
	}

	fanSensorPMObj.Name = Name
	fanSensorPMObj.Class = Class
	switch Class {
	case "Class-A":
		sMgr.fanClassAPMMutex.RLock()
		fanSensorPMObj.Data = sMgr.fanSensorClassAPM[Name].GetListOfEntriesFromRingBuffer()
		sMgr.fanClassAPMMutex.RUnlock()
	case "Class-B":
		sMgr.fanClassBPMMutex.RLock()
		fanSensorPMObj.Data = sMgr.fanSensorClassBPM[Name].GetListOfEntriesFromRingBuffer()
		sMgr.fanClassBPMMutex.RUnlock()
	case "Class-C":
		sMgr.fanClassCPMMutex.RLock()
		fanSensorPMObj.Data = sMgr.fanSensorClassCPM[Name].GetListOfEntriesFromRingBuffer()
		sMgr.fanClassCPMMutex.RUnlock()
	}
	return &fanSensorPMObj, nil
}

func (sMgr *SensorManager) GetTempSensorPMState(Name string, Class string) (*objects.TemperatureSensorPMState, error) {
	var tempSensorPMObj objects.TemperatureSensorPMState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	_, exist := sMgr.tempConfigDB[Name]
	if !exist {
		return nil, errors.New("Invalid Temp Sensor Name")
	}

	tempSensorPMObj.Name = Name
	tempSensorPMObj.Class = Class
	switch Class {
	case "Class-A":
		sMgr.tempClassAPMMutex.RLock()
		tempSensorPMObj.Data = sMgr.tempSensorClassAPM[Name].GetListOfEntriesFromRingBuffer()
		sMgr.tempClassAPMMutex.RUnlock()
	case "Class-B":
		sMgr.tempClassBPMMutex.RLock()
		tempSensorPMObj.Data = sMgr.tempSensorClassBPM[Name].GetListOfEntriesFromRingBuffer()
		sMgr.tempClassBPMMutex.RUnlock()
	case "Class-C":
		sMgr.tempClassCPMMutex.RLock()
		tempSensorPMObj.Data = sMgr.tempSensorClassCPM[Name].GetListOfEntriesFromRingBuffer()
		sMgr.tempClassCPMMutex.RUnlock()
	}
	return &tempSensorPMObj, nil
}

func (sMgr *SensorManager) GetVoltageSensorPMState(Name string, Class string) (*objects.VoltageSensorPMState, error) {
	var voltageSensorPMObj objects.VoltageSensorPMState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	_, exist := sMgr.voltageConfigDB[Name]
	if !exist {
		return nil, errors.New("Invalid Voltage Sensor Name")
	}

	voltageSensorPMObj.Name = Name
	voltageSensorPMObj.Class = Class
	switch Class {
	case "Class-A":
		sMgr.voltageClassAPMMutex.RLock()
		voltageSensorPMObj.Data = sMgr.voltageSensorClassAPM[Name].GetListOfEntriesFromRingBuffer()
		sMgr.voltageClassAPMMutex.RUnlock()
	case "Class-B":
		sMgr.voltageClassBPMMutex.RLock()
		voltageSensorPMObj.Data = sMgr.voltageSensorClassBPM[Name].GetListOfEntriesFromRingBuffer()
		sMgr.voltageClassBPMMutex.RUnlock()
	case "Class-C":
		sMgr.voltageClassCPMMutex.RLock()
		voltageSensorPMObj.Data = sMgr.voltageSensorClassCPM[Name].GetListOfEntriesFromRingBuffer()
		sMgr.voltageClassCPMMutex.RUnlock()
	}
	return &voltageSensorPMObj, nil
}

func (sMgr *SensorManager) GetPowerConverterSensorPMState(Name string, Class string) (*objects.PowerConverterSensorPMState, error) {
	var powerConverterSensorPMObj objects.PowerConverterSensorPMState
	if sMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	_, exist := sMgr.powerConverterConfigDB[Name]
	if !exist {
		return nil, errors.New("Invalid PowerConverter Sensor Name")
	}

	powerConverterSensorPMObj.Name = Name
	powerConverterSensorPMObj.Class = Class
	switch Class {
	case "Class-A":
		sMgr.powerConverterClassAPMMutex.RLock()
		powerConverterSensorPMObj.Data = sMgr.powerConverterSensorClassAPM[Name].GetListOfEntriesFromRingBuffer()
		sMgr.powerConverterClassAPMMutex.RUnlock()
	case "Class-B":
		sMgr.powerConverterClassBPMMutex.RLock()
		powerConverterSensorPMObj.Data = sMgr.powerConverterSensorClassBPM[Name].GetListOfEntriesFromRingBuffer()
		sMgr.powerConverterClassBPMMutex.RUnlock()
	case "Class-C":
		sMgr.powerConverterClassCPMMutex.RLock()
		powerConverterSensorPMObj.Data = sMgr.powerConverterSensorClassCPM[Name].GetListOfEntriesFromRingBuffer()
		sMgr.powerConverterClassCPMMutex.RUnlock()
	}
	return &powerConverterSensorPMObj, nil
}
