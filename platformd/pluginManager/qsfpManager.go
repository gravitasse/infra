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
	//"infra/platformd/pluginManager/pluginCommon"
	"utils/logging"
)

type QsfpState struct {
	Present            bool
	VendorName         string
	VendorOUI          string
	VendorPartNumber   string
	VendorRevision     string
	VendorSerialNumber string
	DataCode           string
	Temperature        float64
	Voltage            float64
	RX1Power           float64
	RX2Power           float64
	RX3Power           float64
	RX4Power           float64
	TX1Power           float64
	TX2Power           float64
	TX3Power           float64
	TX4Power           float64
	TX1Bias            float64
	TX2Bias            float64
	TX3Bias            float64
	TX4Bias            float64
}

type QsfpConfig struct {
	AdminState               string
	HigherAlarmTemperature   float64
	HigherAlarmVoltage       float64
	HigherAlarmRXPower       float64
	HigherAlarmTXPower       float64
	HigherAlarmTXBias        float64
	HigherWarningTemperature float64
	HigherWarningVoltage     float64
	HigherWarningRXPower     float64
	HigherWarningTXPower     float64
	HigherWarningTXBias      float64
	LowerAlarmTemperature    float64
	LowerAlarmVoltage        float64
	LowerAlarmRXPower        float64
	LowerAlarmTXPower        float64
	LowerAlarmTXBias         float64
	LowerWarningTemperature  float64
	LowerWarningVoltage      float64
	LowerWarningRXPower      float64
	LowerWarningTXPower      float64
	LowerWarningTXBias       float64
}

type QsfpManager struct {
	logger   logging.LoggerIntf
	plugin   PluginIntf
	configDB []QsfpConfig
}

var QsfpMgr QsfpManager

func (qMgr *QsfpManager) Init(logger logging.LoggerIntf, plugin PluginIntf) {
	qMgr.logger = logger
	qMgr.plugin = plugin
	numOfQsfps := qMgr.plugin.GetMaxNumOfQsfp()
	qMgr.configDB = make([]QsfpConfig, numOfQsfps)
	for id := 0; id < numOfQsfps; id++ {
		qsfpCfgEnt := qMgr.configDB[id]
		qsfpCfgEnt.AdminState = "Enabled"
		qsfpCfgEnt.HigherAlarmTemperature = 100.0
		qsfpCfgEnt.HigherAlarmVoltage = 10.0
		qsfpCfgEnt.HigherAlarmRXPower = 100.0
		qsfpCfgEnt.HigherAlarmTXPower = 100.0
		qsfpCfgEnt.HigherAlarmTXBias = 100.0
		qsfpCfgEnt.HigherWarningTemperature = 100.0
		qsfpCfgEnt.HigherWarningVoltage = 10.0
		qsfpCfgEnt.HigherWarningRXPower = 100.0
		qsfpCfgEnt.HigherWarningTXPower = 100.0
		qsfpCfgEnt.HigherWarningTXBias = 100.0
		qsfpCfgEnt.LowerAlarmTemperature = -100.0
		qsfpCfgEnt.LowerAlarmVoltage = -10.0
		qsfpCfgEnt.LowerAlarmRXPower = -100.0
		qsfpCfgEnt.LowerAlarmTXPower = -100.0
		qsfpCfgEnt.LowerAlarmTXBias = -100.0
		qsfpCfgEnt.LowerWarningTemperature = -100.0
		qsfpCfgEnt.LowerWarningVoltage = -10.0
		qsfpCfgEnt.LowerWarningRXPower = -100.0
		qsfpCfgEnt.LowerWarningTXPower = -100.0
		qsfpCfgEnt.LowerWarningTXBias = -100.0
		qMgr.configDB[id] = qsfpCfgEnt
	}
	qMgr.logger.Info("Qsfp Manager Init()")
}

func (qMgr *QsfpManager) Deinit() {
	qMgr.logger.Info("Fan Manager Deinit()")
}

func (qMgr *QsfpManager) GetQsfpState(id int32) (*objects.QsfpState, error) {
	var qsfpObj objects.QsfpState
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	numOfQsfp := int32(len(qMgr.configDB))
	if (id > numOfQsfp) && (id < 1) {
		return nil, errors.New("Invalid QsfpId")
	}

	qsfpState, err := qMgr.plugin.GetQsfpState(id)
	if err != nil {
		qsfpObj.QsfpId = id
		qsfpObj.Present = false
	} else {
		qsfpObj.QsfpId = id
		qsfpObj.Present = true
		qsfpObj.VendorName = qsfpState.VendorName
		qsfpObj.VendorOUI = qsfpState.VendorOUI
		qsfpObj.VendorPartNumber = qsfpState.VendorPartNumber
		qsfpObj.VendorRevision = qsfpState.VendorRevision
		qsfpObj.VendorSerialNumber = qsfpState.VendorSerialNumber
		qsfpObj.DataCode = qsfpState.DataCode
		qsfpObj.Temperature = qsfpState.Temperature
		qsfpObj.Voltage = qsfpState.Voltage
		qsfpObj.RX1Power = qsfpState.RX1Power
		qsfpObj.RX2Power = qsfpState.RX2Power
		qsfpObj.RX3Power = qsfpState.RX3Power
		qsfpObj.RX4Power = qsfpState.RX4Power
		qsfpObj.TX1Power = qsfpState.TX1Power
		qsfpObj.TX2Power = qsfpState.TX2Power
		qsfpObj.TX3Power = qsfpState.TX3Power
		qsfpObj.TX4Power = qsfpState.TX4Power
		qsfpObj.TX1Bias = qsfpState.TX1Bias
		qsfpObj.TX2Bias = qsfpState.TX2Bias
		qsfpObj.TX3Bias = qsfpState.TX3Bias
		qsfpObj.TX4Bias = qsfpState.TX4Bias
	}
	return &qsfpObj, err
}

func (qMgr *QsfpManager) GetBulkQsfpState(fromIdx int, cnt int) (*objects.QsfpStateGetInfo, error) {
	var retObj objects.QsfpStateGetInfo
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	if fromIdx >= len(qMgr.configDB) {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > len(qMgr.configDB) {
		retObj.EndIdx = len(qMgr.configDB)
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = len(qMgr.configDB) - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		obj, err := qMgr.GetQsfpState(int32(idx + 1))
		if err != nil {
			qMgr.logger.Err(fmt.Sprintln("Error getting the qsfp state for QsfpId:", idx+1))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func (qMgr *QsfpManager) GetQsfpConfig(id int32) (*objects.QsfpConfig, error) {
	var obj objects.QsfpConfig
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	numOfQsfp := int32(len(qMgr.configDB))
	if (id > numOfQsfp) && (id < int32(1)) {
		return nil, errors.New("Invalid QsfpId")
	}

	qsfpCfgEnt := qMgr.configDB[id-1]
	obj.QsfpId = id
	obj.AdminState = qsfpCfgEnt.AdminState
	obj.HigherAlarmTemperature = qsfpCfgEnt.HigherAlarmTemperature
	obj.HigherAlarmVoltage = qsfpCfgEnt.HigherAlarmVoltage
	obj.HigherAlarmRXPower = qsfpCfgEnt.HigherAlarmRXPower
	obj.HigherAlarmTXPower = qsfpCfgEnt.HigherAlarmTXPower
	obj.HigherAlarmTXBias = qsfpCfgEnt.HigherAlarmTXBias
	obj.HigherWarningTemperature = qsfpCfgEnt.HigherWarningTemperature
	obj.HigherWarningVoltage = qsfpCfgEnt.HigherWarningVoltage
	obj.HigherWarningRXPower = qsfpCfgEnt.HigherWarningRXPower
	obj.HigherWarningTXPower = qsfpCfgEnt.HigherWarningTXPower
	obj.HigherWarningTXBias = qsfpCfgEnt.HigherWarningTXBias
	obj.LowerAlarmTemperature = qsfpCfgEnt.LowerAlarmTemperature
	obj.LowerAlarmVoltage = qsfpCfgEnt.LowerAlarmVoltage
	obj.LowerAlarmRXPower = qsfpCfgEnt.LowerAlarmRXPower
	obj.LowerAlarmTXPower = qsfpCfgEnt.LowerAlarmTXPower
	obj.LowerAlarmTXBias = qsfpCfgEnt.LowerAlarmTXBias
	obj.LowerWarningTemperature = qsfpCfgEnt.LowerWarningTemperature
	obj.LowerWarningVoltage = qsfpCfgEnt.LowerWarningVoltage
	obj.LowerWarningRXPower = qsfpCfgEnt.LowerWarningRXPower
	obj.LowerWarningTXPower = qsfpCfgEnt.LowerWarningTXPower
	obj.LowerWarningTXBias = qsfpCfgEnt.LowerWarningTXBias

	return &obj, nil
}

func (qMgr *QsfpManager) GetBulkQsfpConfig(fromIdx int, cnt int) (*objects.QsfpConfigGetInfo, error) {
	var retObj objects.QsfpConfigGetInfo
	if qMgr.plugin == nil {
		return nil, errors.New("Invalid platform plugin")
	}
	numOfQsfp := len(qMgr.configDB)
	if fromIdx >= numOfQsfp {
		return nil, errors.New("Invalid range")
	}
	if fromIdx+cnt > numOfQsfp {
		retObj.EndIdx = numOfQsfp
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.EndIdx = fromIdx + cnt
		retObj.More = true
		retObj.Count = numOfQsfp - retObj.EndIdx + 1
	}
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		obj, err := qMgr.GetQsfpConfig(int32(idx + 1))
		if err != nil {
			qMgr.logger.Err(fmt.Sprintln("Error getting the Qsfp state for QsfpId:", idx+1))
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}

func genQsfpUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0
	for idx, val := range attrset {
		if true == val {
			switch idx {
			case 0:
				//QSFP Id
			case 1:
				mask |= objects.QSFP_UPDATE_ADMIN_STATE
			case 2:
				mask |= objects.QSFP_UPDATE_HIGHER_ALARM_TEMPERATURE
			case 3:
				mask |= objects.QSFP_UPDATE_HIGHER_ALARM_VOLTAGE
			case 4:
				mask |= objects.QSFP_UPDATE_HIGHER_ALARM_RX_POWER
			case 5:
				mask |= objects.QSFP_UPDATE_HIGHER_ALARM_TX_POWER
			case 6:
				mask |= objects.QSFP_UPDATE_HIGHER_ALARM_TX_BIAS
			case 7:
				mask |= objects.QSFP_UPDATE_HIGHER_WARN_TEMPERATURE
			case 8:
				mask |= objects.QSFP_UPDATE_HIGHER_WARN_VOLTAGE
			case 9:
				mask |= objects.QSFP_UPDATE_HIGHER_WARN_RX_POWER
			case 10:
				mask |= objects.QSFP_UPDATE_HIGHER_WARN_TX_POWER
			case 11:
				mask |= objects.QSFP_UPDATE_HIGHER_WARN_TX_BIAS
			case 12:
				mask |= objects.QSFP_UPDATE_LOWER_ALARM_TEMPERATURE
			case 13:
				mask |= objects.QSFP_UPDATE_LOWER_ALARM_VOLTAGE
			case 14:
				mask |= objects.QSFP_UPDATE_LOWER_ALARM_RX_POWER
			case 15:
				mask |= objects.QSFP_UPDATE_LOWER_ALARM_TX_POWER
			case 16:
				mask |= objects.QSFP_UPDATE_LOWER_ALARM_TX_BIAS
			case 17:
				mask |= objects.QSFP_UPDATE_LOWER_WARN_TEMPERATURE
			case 18:
				mask |= objects.QSFP_UPDATE_LOWER_WARN_VOLTAGE
			case 19:
				mask |= objects.QSFP_UPDATE_LOWER_WARN_RX_POWER
			case 20:
				mask |= objects.QSFP_UPDATE_LOWER_WARN_TX_POWER
			case 21:
				mask |= objects.QSFP_UPDATE_LOWER_WARN_TX_BIAS
			}
		}
	}
	return mask
}

func (qMgr *QsfpManager) UpdateQsfpConfig(oldCfg *objects.QsfpConfig, newCfg *objects.QsfpConfig, attrset []bool) (bool, error) {
	if qMgr.plugin == nil {
		return false, errors.New("Invalid platform plugin")
	}
	numOfQsfp := len(qMgr.configDB)
	if newCfg.QsfpId > int32(numOfQsfp) || newCfg.QsfpId < int32(1) {
		return false, errors.New("Invalid Qsfp Id")
	}
	qsfpCfgEnt := qMgr.configDB[newCfg.QsfpId-1]
	mask := genQsfpUpdateMask(attrset)
	if mask&objects.QSFP_UPDATE_ADMIN_STATE == objects.QSFP_UPDATE_ADMIN_STATE {
		qsfpCfgEnt.AdminState = newCfg.AdminState
	}
	if mask&objects.QSFP_UPDATE_HIGHER_ALARM_TEMPERATURE == objects.QSFP_UPDATE_HIGHER_ALARM_TEMPERATURE {
		qsfpCfgEnt.HigherAlarmTemperature = newCfg.HigherAlarmTemperature
	}
	if mask&objects.QSFP_UPDATE_HIGHER_ALARM_VOLTAGE == objects.QSFP_UPDATE_HIGHER_ALARM_VOLTAGE {
		qsfpCfgEnt.HigherAlarmVoltage = newCfg.HigherAlarmVoltage
	}
	if mask&objects.QSFP_UPDATE_HIGHER_ALARM_RX_POWER == objects.QSFP_UPDATE_HIGHER_ALARM_RX_POWER {
		qsfpCfgEnt.HigherAlarmRXPower = newCfg.HigherAlarmRXPower
	}
	if mask&objects.QSFP_UPDATE_HIGHER_ALARM_TX_POWER == objects.QSFP_UPDATE_HIGHER_ALARM_TX_POWER {
		qsfpCfgEnt.HigherAlarmTXPower = newCfg.HigherAlarmTXPower
	}
	if mask&objects.QSFP_UPDATE_HIGHER_ALARM_TX_BIAS == objects.QSFP_UPDATE_HIGHER_ALARM_TX_BIAS {
		qsfpCfgEnt.HigherAlarmTXBias = newCfg.HigherAlarmTXBias
	}
	if mask&objects.QSFP_UPDATE_HIGHER_WARN_TEMPERATURE == objects.QSFP_UPDATE_HIGHER_WARN_TEMPERATURE {
		qsfpCfgEnt.HigherWarningTemperature = newCfg.HigherWarningTemperature
	}
	if mask&objects.QSFP_UPDATE_HIGHER_WARN_VOLTAGE == objects.QSFP_UPDATE_HIGHER_WARN_VOLTAGE {
		qsfpCfgEnt.HigherWarningVoltage = newCfg.HigherWarningVoltage
	}
	if mask&objects.QSFP_UPDATE_HIGHER_WARN_RX_POWER == objects.QSFP_UPDATE_HIGHER_WARN_RX_POWER {
		qsfpCfgEnt.HigherWarningRXPower = newCfg.HigherWarningRXPower
	}
	if mask&objects.QSFP_UPDATE_HIGHER_WARN_TX_POWER == objects.QSFP_UPDATE_HIGHER_WARN_TX_POWER {
		qsfpCfgEnt.HigherWarningTXPower = newCfg.HigherWarningTXPower
	}
	if mask&objects.QSFP_UPDATE_HIGHER_WARN_TX_BIAS == objects.QSFP_UPDATE_HIGHER_WARN_TX_BIAS {
		qsfpCfgEnt.HigherWarningTXBias = newCfg.HigherWarningTXBias
	}
	if mask&objects.QSFP_UPDATE_LOWER_ALARM_TEMPERATURE == objects.QSFP_UPDATE_LOWER_ALARM_TEMPERATURE {
		qsfpCfgEnt.LowerAlarmTemperature = newCfg.LowerAlarmTemperature
	}
	if mask&objects.QSFP_UPDATE_LOWER_ALARM_VOLTAGE == objects.QSFP_UPDATE_LOWER_ALARM_VOLTAGE {
		qsfpCfgEnt.LowerAlarmVoltage = newCfg.LowerAlarmVoltage
	}
	if mask&objects.QSFP_UPDATE_LOWER_ALARM_RX_POWER == objects.QSFP_UPDATE_LOWER_ALARM_RX_POWER {
		qsfpCfgEnt.LowerAlarmRXPower = newCfg.LowerAlarmRXPower
	}
	if mask&objects.QSFP_UPDATE_LOWER_ALARM_TX_POWER == objects.QSFP_UPDATE_LOWER_ALARM_TX_POWER {
		qsfpCfgEnt.LowerAlarmTXPower = newCfg.LowerAlarmTXPower
	}
	if mask&objects.QSFP_UPDATE_LOWER_ALARM_TX_BIAS == objects.QSFP_UPDATE_LOWER_ALARM_TX_BIAS {
		qsfpCfgEnt.LowerAlarmTXBias = newCfg.LowerAlarmTXBias
	}
	if mask&objects.QSFP_UPDATE_LOWER_WARN_TEMPERATURE == objects.QSFP_UPDATE_LOWER_WARN_TEMPERATURE {
		qsfpCfgEnt.LowerWarningTemperature = newCfg.LowerWarningTemperature
	}
	if mask&objects.QSFP_UPDATE_LOWER_WARN_VOLTAGE == objects.QSFP_UPDATE_LOWER_WARN_VOLTAGE {
		qsfpCfgEnt.LowerWarningVoltage = newCfg.LowerWarningVoltage
	}
	if mask&objects.QSFP_UPDATE_LOWER_WARN_RX_POWER == objects.QSFP_UPDATE_LOWER_WARN_RX_POWER {
		qsfpCfgEnt.LowerWarningRXPower = newCfg.LowerWarningRXPower
	}
	if mask&objects.QSFP_UPDATE_LOWER_WARN_TX_POWER == objects.QSFP_UPDATE_LOWER_WARN_TX_POWER {
		qsfpCfgEnt.LowerWarningTXPower = newCfg.LowerWarningTXPower
	}
	if mask&objects.QSFP_UPDATE_LOWER_WARN_TX_BIAS == objects.QSFP_UPDATE_LOWER_WARN_TX_BIAS {
		qsfpCfgEnt.LowerWarningTXBias = newCfg.LowerWarningTXBias
	}

	return true, nil
}
