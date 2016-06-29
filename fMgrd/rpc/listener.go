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
	"errors"
	"fMgrd"
	"fmt"
	"infra/fMgrd/server"
	"utils/logging"
)

type FMGRHandler struct {
	server *server.FMGRServer
	logger *logging.Writer
}

func NewFMGRHandler(server *server.FMGRServer, logger *logging.Writer) *FMGRHandler {
	h := new(FMGRHandler)
	h.server = server
	h.logger = logger
	return h
}

func (h *FMGRHandler) CreateFMgrGlobal(conf *fMgrd.FMgrGlobal) (bool, error) {
	h.logger.Info(fmt.Sprintln("Received CreateFMgrGlobal call"))
	return true, nil
}

func (h *FMGRHandler) DeleteFMgrGlobal(conf *fMgrd.FMgrGlobal) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete FMgr config attrs:", conf))
	return true, nil
}

func (h *FMGRHandler) UpdateFMgrGlobal(origConf *fMgrd.FMgrGlobal, newConf *fMgrd.FMgrGlobal, attrset []bool, op []*fMgrd.PatchOpInfo) (bool, error) {
	return true, nil
}

func (h *FMGRHandler) convertFaultEntryToThrift(faultState server.FaultState) *fMgrd.FaultState {
	fault := fMgrd.NewFaultState()
	fault.OwnerId = int32(faultState.OwnerId)
	fault.EventId = int32(faultState.EventId)
	fault.OwnerName = faultState.OwnerName
	fault.EventName = faultState.EventName
	fault.SrcObjName = faultState.SrcObjName
	fault.Description = faultState.Description
	fault.OccuranceTime = faultState.OccuranceTime
	fault.SrcObjKey = faultState.SrcObjKey
	fault.ResolutionTime = faultState.ResolutionTime
	return fault
}

func (h *FMGRHandler) GetBulkFaultState(fromIndex fMgrd.Int, count fMgrd.Int) (*fMgrd.FaultStateGetInfo, error) {
	h.logger.Info(fmt.Sprintln("Get bulk call for Faults"))
	nextIdx, currCount, faultEntry := h.server.GetBulkFault(int(fromIndex), int(count))
	if faultEntry == nil {
		err := errors.New("Fault Manager server unable to fetch fault entry")
		return nil, err
	}
	faultEntryResponse := make([]*fMgrd.FaultState, len(faultEntry))
	for idx, item := range faultEntry {
		faultEntryResponse[idx] = h.convertFaultEntryToThrift(item)
	}

	faultEntryBulk := fMgrd.NewFaultStateGetInfo()
	faultEntryBulk.Count = fMgrd.Int(currCount)
	faultEntryBulk.StartIdx = fMgrd.Int(fromIndex)
	faultEntryBulk.EndIdx = fMgrd.Int(nextIdx)
	faultEntryBulk.FaultStateList = faultEntryResponse
	return faultEntryBulk, nil
}

func (h *FMGRHandler) GetFaultState(ownerId int32, eventId int32, ownerName string, eventName string, srcObjName string) (*fMgrd.FaultState, error) {
	return nil, nil
}
