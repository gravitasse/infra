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
	"models/events"
	"time"
	"utils/logging"
)

type FMGRServer struct {
	logger   *logging.Writer
	dbHdl    redis.Conn
	subHdl   redis.PubSubConn
	InitDone chan bool
}

func NewFMGRServer(logger *logging.Writer) *FMGRServer {
	fMgrServer := &FMGRServer{}
	fMgrServer.logger = logger
	fMgrServer.InitDone = make(chan bool)
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

func (server *FMGRServer) RedisSub() {
	for {
		switch n := server.subHdl.Receive().(type) {
		case redis.Message:
			var evt events.Event
			err := json.Unmarshal(n.Data, &evt)
			if err != nil {
				server.logger.Err(fmt.Sprintln("Unable to Unmarshal the byte stream", err))
				continue
			}
			server.logger.Info(fmt.Sprintln("OwnerId:", evt.OwnerId))
			server.logger.Info(fmt.Sprintln("OwnerName:", evt.OwnerName))
			server.logger.Info(fmt.Sprintln("EvtId:", evt.EvtId))
			server.logger.Info(fmt.Sprintln("EventName:", evt.EventName))
			server.logger.Info(fmt.Sprintln("Timestamp:", evt.TimeStamp))
			server.logger.Info(fmt.Sprintln("Description:", evt.Description))
			server.logger.Info(fmt.Sprintln("SrcObjName:", evt.SrcObjName))
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
			if err != nil {
				server.logger.Err(fmt.Sprintln("Unable to Unmarshal the byte stream", err))
				continue
			}
			server.logger.Info(fmt.Sprintln("Src Obj Key", obj))

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

func (server *FMGRServer) InitServer(paramDir string) {
	var err error
	server.dbHdl, err = server.dial()
	if err != nil {
		server.logger.Err(fmt.Sprintln(err))
		return
	}
	//defer server.dbHdl.Close()
	server.subHdl = redis.PubSubConn{Conn: server.dbHdl}

	server.subHdl.Subscribe("ASICD")
	server.subHdl.Subscribe("ARPD")
	go server.RedisSub()

}

func (server *FMGRServer) StartServer(paramDir string) {
	server.InitServer(paramDir)
	server.InitDone <- true
	for {
		time.Sleep(10 * time.Second)
	}
}
