//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
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
	"fMgrd"
	"fmt"
	//"infra/fMgrd/server"
)

func (h *FMGRHandler) CreateFMgrGlobal(conf *fMgrd.FMgrGlobal) (bool, error) {
	h.logger.Info(fmt.Sprintln("Received CreateFMgrGlobal call"))
	/*
		err := h.SendSetArpGlobalConfig(int(conf.Timeout))
		if err != nil {
			return false, err
		}
	*/
	return true, nil
}

func (h *FMGRHandler) DeleteFMgrGlobal(conf *fMgrd.FMgrGlobal) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete FMgr config attrs:", conf))
	return true, nil
}

func (h *FMGRHandler) UpdateFMgrGlobal(origConf *fMgrd.FMgrGlobal, newConf *fMgrd.FMgrGlobal, attrset []bool, op []*fMgrd.PatchOpInfo) (bool, error) {
	/*
	   h.logger.Info(fmt.Sprintln("Original Arp config attrs:", origConf))
	   h.logger.Info(fmt.Sprintln("New Arp config attrs:", newConf))
	   err := h.SendUpdateArpGlobalConfig(int(newConf.Timeout))
	   if err != nil {
	           return false, err
	   }
	*/
	return true, nil
}
