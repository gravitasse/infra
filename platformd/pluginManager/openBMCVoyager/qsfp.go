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

package openBMCVoyager

import (
	"errors"
	"fmt"
	"infra/platformd/pluginManager/pluginCommon"
)

/*
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include "openBMCVoyager.h"
#include "pluginCommon.h"
*/
import "C"

const (
	MAX_NUM_OF_QSFP int = 12
)

func (driver *openBMCVoyagerDriver) GetQsfpState(Id int32) (retObj pluginCommon.QsfpState, err error) {
	var qsfpInfo C.qsfp_info_t

	retval := int(C.GetQsfpState(&qsfpInfo, C.int(Id)))
	if retval < 0 {
		return retObj, errors.New(fmt.Sprintln("Unable to fetch qsft state of", Id))
	}
	retObj.VendorName = C.GoString(&qsfpInfo.VendorName[0])
	retObj.VendorOUI = C.GoString(&qsfpInfo.VendorOUI[0])
	retObj.VendorPartNumber = C.GoString(&qsfpInfo.VendorPN[0])
	retObj.VendorRevision = C.GoString(&qsfpInfo.VendorRev[0])
	retObj.VendorSerialNumber = C.GoString(&qsfpInfo.VendorSN[0])
	retObj.DataCode = C.GoString(&qsfpInfo.DataCode[0])
	retObj.Temperature = float64(qsfpInfo.Temperature)
	retObj.Voltage = float64(qsfpInfo.SupplyVoltage)
	retObj.RX1Power = float64(qsfpInfo.RX1Power)
	retObj.RX2Power = float64(qsfpInfo.RX2Power)
	retObj.RX3Power = float64(qsfpInfo.RX3Power)
	retObj.RX4Power = float64(qsfpInfo.RX4Power)
	retObj.TX1Power = float64(qsfpInfo.TX1Power)
	retObj.TX2Power = float64(qsfpInfo.TX2Power)
	retObj.TX3Power = float64(qsfpInfo.TX3Power)
	retObj.TX4Power = float64(qsfpInfo.TX4Power)
	retObj.TX1Bias = float64(qsfpInfo.TX1Bias)
	retObj.TX2Bias = float64(qsfpInfo.TX2Bias)
	retObj.TX3Bias = float64(qsfpInfo.TX3Bias)
	retObj.TX4Bias = float64(qsfpInfo.TX4Bias)
	return retObj, nil
}

func (driver *openBMCVoyagerDriver) GetMaxNumOfQsfp() int {
	driver.logger.Info("Inside OpenBMC Voyager: GetMaxNumOfQsfps()")
	return MAX_NUM_OF_QSFP
}
