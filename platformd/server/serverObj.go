//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
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
	"infra/platformd/objects"
)

type ServerOpId int

const (
	GET_FAN_STATE ServerOpId = iota
	GET_BULK_FAN_STATE
	GET_FAN_CONFIG
	GET_BULK_FAN_CONFIG
	UPDATE_FAN_CONFIG
	GET_SFP_STATE
	GET_BULK_SFP_STATE
	GET_SFP_CONFIG
	GET_BULK_SFP_CONFIG
	UPDATE_SFP_CONFIG
	GET_PSU_STATE
	GET_BULK_PSU_STATE
	GET_PSU_CONFIG
	GET_BULK_PSU_CONFIG
	UPDATE_PSU_CONFIG
	GET_PLATFORM_STATE
	GET_BULK_PLATFORM_STATE
	GET_THERMAL_STATE
	GET_BULK_THERMAL_STATE
)

type ServerRequest struct {
	Op   ServerOpId
	Data interface{}
}

type GetBulkInArgs struct {
	FromIdx int
	Count   int
}

type GetFanStateInArgs struct {
	FanId int32
}

type GetFanStateOutArgs struct {
	Obj *objects.FanState
	Err error
}

type GetBulkFanStateOutArgs struct {
	BulkInfo *objects.FanStateGetInfo
	Err      error
}

type GetFanConfigInArgs struct {
	FanId int32
}

type GetFanConfigOutArgs struct {
	Obj *objects.FanConfig
	Err error
}

type GetBulkFanConfigOutArgs struct {
	BulkInfo *objects.FanConfigGetInfo
	Err      error
}

type UpdateFanConfigInArgs struct {
	FanOldCfg *objects.FanConfig
	FanNewCfg *objects.FanConfig
	AttrSet   []bool
}

type UpdateConfigOutArgs struct {
	RetVal bool
	Err    error
}

type GetSfpStateInArgs struct {
	SfpId int32
}

type GetSfpStateOutArgs struct {
	Obj *objects.SfpState
	Err error
}

type GetBulkSfpStateOutArgs struct {
	BulkInfo *objects.SfpStateGetInfo
	Err      error
}

type GetSfpConfigInArgs struct {
	SfpId int32
}

type GetSfpConfigOutArgs struct {
	Obj *objects.SfpConfig
	Err error
}

type GetBulkSfpConfigOutArgs struct {
	BulkInfo *objects.SfpConfigGetInfo
	Err      error
}

type UpdateSfpConfigInArgs struct {
	SfpOldCfg *objects.SfpConfig
	SfpNewCfg *objects.SfpConfig
	AttrSet   []bool
}

type GetPlatformStateInArgs struct {
	ObjName string
}

type GetPlatformStateOutArgs struct {
	Obj *objects.PlatformState
	Err error
}

type GetBulkPlatformStateOutArgs struct {
	BulkInfo *objects.PlatformStateGetInfo
	Err      error
}

type GetThermalStateInArgs struct {
	ThermalId int32
}

type GetThermalStateOutArgs struct {
	Obj *objects.ThermalState
	Err error
}

type GetBulkThermalStateOutArgs struct {
	BulkInfo *objects.ThermalStateGetInfo
	Err      error
}
