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

package faultMgr

import (
	"errors"
	"fmt"
	"strings"
)

func (fMgr *FaultManager) generateFaultObjKey(srcObjName string, srcObjKey interface{}) (FaultObjKey, error) {
	key := fmt.Sprintf("%v", srcObjKey)
	str := strings.Split(fmt.Sprintf("%v", key), "map[")
	key = strings.Split(str[1], "]")[0]
	srcObjUUID, err := fMgr.getUUID(srcObjName, key)
	if err != nil {
		fMgr.logger.Err("Unable to find the UUID of", srcObjName, srcObjKey, err)
		return "", errors.New(fmt.Sprintln("Unable to find the UUID of", srcObjName, srcObjKey, err))
	}
	return FaultObjKey(fmt.Sprintf("%s#%s#%s", srcObjName, key, srcObjUUID)), err
}

func getObjKey(srcObjName string, srcObjKey string) (str string) {
	str = srcObjName
	keyVal := strings.Split(srcObjKey, " ")
	for _, kv := range keyVal {
		val := strings.Split(strings.TrimSpace(kv), ":")[1]
		str = str + "#" + val
	}
	return str
}

func (fMgr *FaultManager) getUUID(srcObjName, srcObjKey string) (uuid string, err error) {
	objKey := getObjKey(srcObjName, srcObjKey)
	return fMgr.dbHdl.GetUUIDFromObjKey(objKey)
}

func getResolutionReason(reason Reason) string {
	switch reason {
	case AUTOCLEARED:
		return "Automatically Cleared"
	case FAULTDISABLED:
		return "Cleared because of FaultEnable(Enable=false) Action"
	case FAULTCLEARED:
		return "Cleared because of FaultClear Action"
	}
	return "Unknown"
}
