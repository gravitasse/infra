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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"models/events"
	"time"
	"utils/logging"
	"utils/ringBuffer"
)

type FaultDetail struct {
	RaiseFault       bool
	ClearingEventId  int
	ClearingDaemonId int
}

type EventStruct struct {
	EventId     int
	EventName   string
	Description string
	SrcObjName  string
	EventEnable bool
	IsFault     bool
	Fault       FaultDetail
}

type DaemonEvent struct {
	DaemonId          int
	DaemonName        string
	DaemonEventEnable bool
	EventList         []EventStruct
}

type EventJson struct {
	DaemonEvents []DaemonEvent
}

const (
	EventDir string = "/etc/flexswitch/"
)

type FaultId struct {
	DaemonId int
	EventId  int
}

type NonFaultData struct {
	IsClearingEvent bool
	FaultEventId    int
	FaultOwnerId    int
}

type EvtDetail struct {
	IsClearingEvent  bool
	IsFault          bool
	RaiseFault       bool
	ClearingEventId  int
	ClearingDaemonId int
}

type FaultDataMap map[FaultObjKey]FaultData

type FaultDatabaseKey struct {
	FaultId        FaultId
	FObjKey        FaultObjKey
	Resolved       bool
	FaultSeqNumber uint64
	OwnerName      string
	EventName      string
	OccuranceTime  time.Time
	ResolutionTime time.Time
	Description    string
}

type FMGRServer struct {
	logger           *logging.Writer
	dbHdl            redis.Conn
	subHdl           redis.PubSubConn
	FaultEventMap    map[FaultId]FaultDetail
	NonFaultEventMap map[FaultId]NonFaultData
	FaultDatabase    map[FaultId]FaultDataMap
	FaultList        *ringBuffer.RingBuffer
	FaultSeqNumber   uint64
	InitDone         chan bool
}

func NewFMGRServer(logger *logging.Writer) *FMGRServer {
	fMgrServer := &FMGRServer{}
	fMgrServer.logger = logger
	fMgrServer.InitDone = make(chan bool)
	fMgrServer.FaultEventMap = make(map[FaultId]FaultDetail)
	fMgrServer.NonFaultEventMap = make(map[FaultId]NonFaultData)
	fMgrServer.FaultDatabase = make(map[FaultId]FaultDataMap)
	fMgrServer.FaultList = new(ringBuffer.RingBuffer)
	fMgrServer.FaultList.SetRingBufferCapacity(100000) //Max 100000 entries in fault database
	fMgrServer.FaultSeqNumber = 0
	return fMgrServer
}

func (server *FMGRServer) dial() (redis.Conn, error) {
	retryCount := 0
	ticker := time.NewTicker(2 * time.Second)
	for _ = range ticker.C {
		retryCount += 1
		dbHdl, err := redis.Dial("tcp", ":6379")
		if err != nil {
			if retryCount%100 == 0 {
				server.logger.Err(fmt.Sprintln("Failed to dail out to Redis server. Retrying connection. Num of retries = ", retryCount))
			}
		} else {
			return dbHdl, nil
		}
	}
	err := errors.New("Error opening db handler")
	return nil, err
}

func (server *FMGRServer) Subscriber() {
	for {
		switch n := server.subHdl.Receive().(type) {
		case redis.Message:
			var evt events.Event
			err := json.Unmarshal(n.Data, &evt)
			if err != nil {
				server.logger.Err(fmt.Sprintln("Unable to Unmarshal the byte stream", err))
				continue
			}
			//server.logger.Info(fmt.Sprintln("OwnerId:", evt.OwnerId))
			//server.logger.Info(fmt.Sprintln("OwnerName:", evt.OwnerName))
			//server.logger.Info(fmt.Sprintln("EvtId:", evt.EvtId))
			//server.logger.Info(fmt.Sprintln("EventName:", evt.EventName))
			//server.logger.Info(fmt.Sprintln("Timestamp:", evt.TimeStamp))
			//server.logger.Info(fmt.Sprintln("Description:", evt.Description))
			//server.logger.Info(fmt.Sprintln("SrcObjName:", evt.SrcObjName))
			keyMap, exist := events.EventKeyMap[evt.OwnerName]
			if !exist {
				server.logger.Info("Key map not found given Event")
				continue
			}
			obj, exist := keyMap[evt.SrcObjName]
			if !exist {
				server.logger.Info("Src Object Name not found in Event Key Map")
				continue
			}
			obj = evt.SrcObjKey
			server.logger.Debug(fmt.Sprintln("Src Obj Key", obj))

			server.processEvents(evt)

		case redis.Subscription:
			if n.Count == 0 {
				return
			}
		case error:
			server.logger.Err(fmt.Sprintf("error: %v\n", n))
			return
		}
	}
	server.subHdl.Unsubscribe()
	server.subHdl.PUnsubscribe()
}

func (server *FMGRServer) initFaultMgrDS() error {
	var evtJson EventJson
	evtMap := make(map[FaultId]EvtDetail)
	eventsFile := EventDir + "events.json"
	bytes, err := ioutil.ReadFile(eventsFile)
	if err != nil {
		server.logger.Err(fmt.Sprintln("Error in reading ", eventsFile, " file."))
		err := errors.New(fmt.Sprintln("Error in reading ", eventsFile, " file."))
		return err
	}

	err = json.Unmarshal(bytes, &evtJson)
	if err != nil {
		server.logger.Err(fmt.Sprintln("Errors in unmarshalling json file : ", eventsFile))
		err := errors.New(fmt.Sprintln("Errors in unmarshalling json file: ", eventsFile))
		return err
	}

	server.logger.Info(fmt.Sprintln("evtJson:", evtJson))
	for _, daemon := range evtJson.DaemonEvents {
		server.logger.Info(fmt.Sprintln("daemon.DaemonName:", daemon.DaemonName))
		for _, evt := range daemon.EventList {
			fId := FaultId{
				DaemonId: int(daemon.DaemonId),
				EventId:  int(evt.EventId),
			}
			evtEnt, exist := evtMap[fId]
			if exist {
				server.logger.Err(fmt.Sprintln("Duplicate entry found"))
				continue
			}
			if evt.IsFault == true {
				evtEnt.IsFault = true
				evtEnt.IsClearingEvent = false
				evtEnt.RaiseFault = evt.Fault.RaiseFault
				evtEnt.ClearingEventId = evt.Fault.ClearingEventId
				evtEnt.ClearingDaemonId = evt.Fault.ClearingDaemonId
			} else {
				evtEnt.IsFault = false
				evtEnt.IsClearingEvent = false
				evtEnt.RaiseFault = false
				evtEnt.ClearingEventId = -1
				evtEnt.ClearingDaemonId = -1
			}
			evtMap[fId] = evtEnt
		}
	}

	for fId, evt := range evtMap {
		if evt.IsFault == true {
			cFId := FaultId{
				DaemonId: evt.ClearingDaemonId,
				EventId:  evt.ClearingEventId,
			}
			cEvt, exist := evtMap[cFId]
			if !exist {
				server.logger.Err(fmt.Sprintln("No clearing event found for fault:", fId))
				continue
			}

			cEvt.IsClearingEvent = true
			evtMap[cFId] = cEvt
		}
	}

	for fId, evt := range evtMap {
		if evt.IsFault == true {
			evtEnt, _ := server.FaultEventMap[fId]
			evtEnt.RaiseFault = evt.RaiseFault
			evtEnt.ClearingEventId = evt.ClearingEventId
			evtEnt.ClearingDaemonId = evt.ClearingDaemonId
			server.FaultEventMap[fId] = evtEnt
			cFId := FaultId{
				DaemonId: evtEnt.ClearingDaemonId,
				EventId:  evtEnt.ClearingEventId,
			}
			cEvtEnt, _ := server.NonFaultEventMap[cFId]
			cEvtEnt.FaultOwnerId = fId.DaemonId
			cEvtEnt.FaultEventId = fId.EventId
			server.NonFaultEventMap[cFId] = cEvtEnt
		} else {
			evtEnt, _ := server.NonFaultEventMap[fId]
			evtEnt.IsClearingEvent = evt.IsClearingEvent
			server.NonFaultEventMap[fId] = evtEnt
		}
	}
	return nil
}

func (server *FMGRServer) InitServer(paramDir string) {
	var err error

	err = server.initFaultMgrDS()
	if err != nil {
		server.logger.Err(fmt.Sprintln(err))
		return
	}

	server.logger.Info(fmt.Sprintln("Non Fault Event:", server.NonFaultEventMap))
	server.logger.Info(fmt.Sprintln("Fault Event:", server.FaultEventMap))

	server.dbHdl, err = server.dial()
	if err != nil {
		server.logger.Err(fmt.Sprintln(err))
		return
	}
	//defer server.dbHdl.Close()
	server.subHdl = redis.PubSubConn{Conn: server.dbHdl}

	server.subHdl.Subscribe("ASICD")
	server.subHdl.Subscribe("ARPD")
	go server.Subscriber()

}

func (server *FMGRServer) StartServer(paramDir string) {
	server.InitServer(paramDir)
	server.InitDone <- true
	for {
		time.Sleep(10 * time.Second)
	}
}
