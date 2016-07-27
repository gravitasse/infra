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
	"fmt"
	"infra/platformd/objects"
	"infra/platformd/pluginManager/onlp"
	"infra/platformd/pluginManager/openBMC"
	"infra/platformd/pluginManager/pluginCommon"
	"strings"
	"utils/logging"
)

type PluginIntf interface {
	Init() error
	DeInit() error
	GetFanState(fanId int32) (pluginCommon.FanState, error)
	GetFanConfig(fanId int32) (*objects.FanConfig, error)
	UpdateFanConfig(cfg *objects.FanConfig) (bool, error)
	GetMaxNumOfFans() int
	GetAllFanState(state []pluginCommon.FanState, count int) error
}

type ResourceManagers struct {
	*SysManager
	*FanManager
	*PsuManager
	*SfpManager
	*ThermalManager
	*LedManager
}

type PluginManager struct {
	*ResourceManagers
	logger logging.LoggerIntf
	plugin PluginIntf
}

func NewPluginMgr(pluginName string, initParams *pluginCommon.PluginInitParams) *PluginManager {
	var plugin PluginIntf
	pluginMgr := new(PluginManager)
	pluginMgr.ResourceManagers = new(ResourceManagers)
	pluginMgr.logger = initParams.Logger
	pluginName = strings.ToLower(pluginName)
	switch pluginName {
	case pluginCommon.ONLP_PLUGIN:
		fmt.Println("===== ONLP_PLUGIN =====")
		plugin = onlp.NewONLPPlugin(initParams)
		pluginMgr.plugin = plugin
	case pluginCommon.OpenBMC_PLUGIN:
		fmt.Println("===== OPENBMC_PLUGIN =====")
		plugin = openBMC.NewOpenBMCPlugin(initParams)
		pluginMgr.plugin = plugin
	default:
	}
	pluginMgr.SysManager = &SysMgr
	pluginMgr.FanManager = &FanMgr
	pluginMgr.PsuManager = &PsuMgr
	pluginMgr.SfpManager = &SfpMgr
	pluginMgr.ThermalManager = &ThermalMgr
	pluginMgr.LedManager = &LedMgr
	return pluginMgr
}

func (pMgr *PluginManager) Init() error {
	pMgr.plugin.Init()
	pMgr.SysManager.Init(pMgr.logger, pMgr.plugin)
	pMgr.FanManager.Init(pMgr.logger, pMgr.plugin)
	pMgr.PsuManager.Init(pMgr.logger, pMgr.plugin)
	pMgr.SfpManager.Init(pMgr.logger, pMgr.plugin)
	pMgr.ThermalManager.Init(pMgr.logger, pMgr.plugin)
	pMgr.LedManager.Init(pMgr.logger, pMgr.plugin)
	return nil
}

func (pMgr *PluginManager) Deinit() {
	pMgr.SysManager.Deinit()
	pMgr.FanManager.Deinit()
	pMgr.PsuManager.Deinit()
	pMgr.SfpManager.Deinit()
	pMgr.ThermalManager.Deinit()
	pMgr.LedManager.Deinit()
}
