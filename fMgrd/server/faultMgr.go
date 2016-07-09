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

package server

import (
	"errors"
	"fmt"
	"models/events"
	"time"
)

type FaultObjKey string

type FaultData struct {
	OwnerName      string
	EventName      string
	Desc           string
	TimeStamp      time.Time
	SrcObjName     string
	SrcObjKey      interface{}
	FaultListIdx   int
	FaultSeqNumber uint64
}

func generateFaultObjKey(ownerName string, srcObjName string, srcObjKey interface{}) FaultObjKey {
	keyMap, exist := events.EventKeyMap[ownerName]
	if !exist {
		return ""
	}
	obj, exist := keyMap[srcObjName]
	if !exist {
		return ""
	}
	obj = srcObjKey
	return FaultObjKey(fmt.Sprintf("%s %s", srcObjName, obj))
}

func (server *FMGRServer) processEvents(evt events.Event) error {
	fId := FaultId{
		DaemonId: int(evt.OwnerId),
		EventId:  int(evt.EvtId),
	}

	server.logger.Info("Process Events....")
	if _, exist := server.FaultEventMap[fId]; exist {
		err := server.processFaultEvents(evt)
		//server.logger.Info(fmt.Sprintln("====1 Fault Database ====", server.FaultDatabase))
		//server.logger.Info(fmt.Sprintln("====Ring Buffer Size===", server.FaultList.GetListOfEntriesFromRingBuffer()))
		return err
	}

	if ent, exist := server.NonFaultEventMap[fId]; exist {
		if ent.IsClearingEvent == true {
			err := server.processFaultClearingEvents(evt)
			//server.logger.Info(fmt.Sprintln("====2 Fault Database ====", server.FaultDatabase))
			//server.logger.Info(fmt.Sprintln("====Ring Buffer Size===", server.FaultList.GetListOfEntriesFromRingBuffer()))
			return err
		}
	}

	return nil
}

func (server *FMGRServer) processFaultEvents(evt events.Event) error {
	server.logger.Info(fmt.Sprintln("Processing Fault Events: ", evt))
	server.CreateEntryInFaultDatabase(evt)
	return nil
}

func (server *FMGRServer) processFaultClearingEvents(evt events.Event) error {
	server.logger.Info(fmt.Sprintln("Processing Non Fault Events: ", evt))
	server.DeleteEntryFromFaultDatabase(evt)
	return nil
}

func (server *FMGRServer) AddFaultEntryInList(evt events.Event) int {
	fId := FaultId{
		DaemonId: int(evt.OwnerId),
		EventId:  int(evt.EvtId),
	}

	fObjKey := generateFaultObjKey(evt.OwnerName, evt.SrcObjName, evt.SrcObjKey)

	fDBKey := FaultDatabaseKey{
		FaultId:        fId,
		FObjKey:        fObjKey,
		Resolved:       false,
		FaultSeqNumber: server.FaultSeqNumber,
		OwnerName:      evt.OwnerName,
		EventName:      evt.EventName,
		OccuranceTime:  evt.TimeStamp,
		Description:    evt.Description,
	}

	return server.FaultList.InsertIntoRingBuffer(fDBKey)
}

func (server *FMGRServer) CreateEntryInFaultDatabase(evt events.Event) error {
	fId := FaultId{
		DaemonId: int(evt.OwnerId),
		EventId:  int(evt.EvtId),
	}

	if server.FaultDatabase[fId] == nil {
		server.FaultDatabase[fId] = make(map[FaultObjKey]FaultData)
	}
	fDbEnt, _ := server.FaultDatabase[fId]
	fObjKey := generateFaultObjKey(evt.OwnerName, evt.SrcObjName, evt.SrcObjKey)
	if fObjKey == "" {
		err := errors.New("Error generating fault object key")
		return err
	}
	fDataEnt, exist := fDbEnt[fObjKey]
	if exist {
		server.logger.Info("Already have corresponding fault in fault Database")
		return nil
	}
	idx := server.AddFaultEntryInList(evt)
	if idx == -1 {
		server.logger.Err("Unable to add entry in fault database")
		err := errors.New("Unable to add entry in fault database")
		return err
	}
	fDataEnt.OwnerName = evt.OwnerName
	fDataEnt.EventName = evt.EventName
	fDataEnt.Desc = evt.Description
	fDataEnt.TimeStamp = evt.TimeStamp
	fDataEnt.SrcObjName = evt.SrcObjName
	fDataEnt.SrcObjKey = evt.SrcObjKey
	fDataEnt.FaultListIdx = idx
	fDataEnt.FaultSeqNumber = server.FaultSeqNumber
	server.FaultSeqNumber++
	fDbEnt[fObjKey] = fDataEnt
	server.FaultDatabase[fId] = fDbEnt
	return nil
}

func (server *FMGRServer) DeleteEntryFromFaultDatabase(evt events.Event) error {
	cFId := FaultId{
		DaemonId: int(evt.OwnerId),
		EventId:  int(evt.EvtId),
	}

	cFEnt, _ := server.NonFaultEventMap[cFId]
	fId := FaultId{
		DaemonId: cFEnt.FaultOwnerId,
		EventId:  cFEnt.FaultEventId,
	}

	fObjKey := generateFaultObjKey(evt.OwnerName, evt.SrcObjName, evt.SrcObjKey)
	if fObjKey == "" {
		err := errors.New("Error generating fault object key")
		return err
	}
	//server.logger.Info(fmt.Sprintln("fObjKey:", fObjKey, "fId:", fId, "FaultDatabase:", server.FaultDatabase))
	fDbEnt, exist := server.FaultDatabase[fId]
	if !exist {
		server.logger.Info(fmt.Sprintln("No such fault occured to be cleared, no entry found in fault database", evt))
		return nil
	} else {
		fDataEnt, exist := fDbEnt[fObjKey]
		if !exist {
			server.logger.Info(fmt.Sprintln("No such fault occured to be cleared, no entry found in fault database entry", evt))
			return nil
		}

		intf := server.FaultList.GetEntryFromRingBuffer(fDataEnt.FaultListIdx)
		fDBKey := intf.(FaultDatabaseKey)
		if fDataEnt.FaultSeqNumber == fDBKey.FaultSeqNumber {
			fDBKey.Resolved = true
			fDBKey.ResolutionTime = evt.TimeStamp
			server.FaultList.UpdateEntryInRingBuffer(fDBKey, fDataEnt.FaultListIdx)
		}
		delete(fDbEnt, fObjKey)
	}
	server.FaultDatabase[fId] = fDbEnt
	return nil
}
