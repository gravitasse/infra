package rpc

import (
	"errors"
	"fmt"
	"infra/sysd/server"
	"sysd"
	"utils/logging"
)

type SYSDHandler struct {
	server *server.SYSDServer
	logger *logging.Writer
}

func NewSYSDHandler(logger *logging.Writer, server *server.SYSDServer) *SYSDHandler {
	h := new(SYSDHandler)
	h.server = server
	h.logger = logger
	return h
}

func (h *SYSDHandler) SendGlobalLoggingConfig(gLoggingConfig *sysd.SystemLogging) bool {
	var Logging bool
	if gLoggingConfig.SystemLogging == "on" {
		Logging = true
	} else {
		Logging = false
	}
	gConf := server.GlobalLoggingConfig{
		Enable: Logging,
	}
	h.server.GlobalLoggingConfigCh <- gConf
	return true
}

func (h *SYSDHandler) SendComponentLoggingConfig(cLoggingConfig *sysd.ComponentLogging) bool {
	cConf := server.ComponentLoggingConfig{
		Component: cLoggingConfig.Module,
		Level:     logging.ConvertLevelStrToVal(cLoggingConfig.Level),
	}
	h.server.ComponentLoggingConfigCh <- cConf
	return true
}

func (h *SYSDHandler) UpdateSystemLogging(origConf *sysd.SystemLogging, newConf *sysd.SystemLogging, attrset []bool) (bool, error) {
	h.logger.Info(fmt.Sprintln("Original global config attrs:", origConf))
	if newConf == nil {
		err := errors.New("Invalid global Configuration")
		return false, err
	}
	h.logger.Info(fmt.Sprintln("Update global config attrs:", newConf))
	return h.SendGlobalLoggingConfig(newConf), nil
}

func (h *SYSDHandler) UpdateComponentLogging(origConf *sysd.ComponentLogging, newConf *sysd.ComponentLogging, attrset []bool) (bool, error) {
	h.logger.Info(fmt.Sprintln("Original component config attrs:", origConf))
	if newConf == nil {
		err := errors.New("Invalid component Configuration")
		return false, err
	}
	h.logger.Info(fmt.Sprintln("Update component config attrs:", newConf))
	return h.SendComponentLoggingConfig(newConf), nil
}

func (h *SYSDHandler) CreateSystemLogging(gLoggingConf *sysd.SystemLogging) (bool, error) {
	h.logger.Info(fmt.Sprintln("Create global config attrs:", gLoggingConf))
	return h.SendGlobalLoggingConfig(gLoggingConf), nil
}

func (h *SYSDHandler) CreateComponentLogging(cLoggingConf *sysd.ComponentLogging) (bool, error) {
	h.logger.Info(fmt.Sprintln("Create component config attrs:", cLoggingConf))
	return h.SendComponentLoggingConfig(cLoggingConf), nil
}

func (h *SYSDHandler) DeleteSystemLogging(gLoggingConf *sysd.SystemLogging) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete global config attrs:", gLoggingConf))
	return true, nil
}

func (h *SYSDHandler) DeleteComponentLogging(cLoggingConf *sysd.ComponentLogging) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete component config attrs:", cLoggingConf))
	return true, nil
}

func (h *SYSDHandler) CreateIpTableAcl(ipaclConfig *sysd.IpTableAcl) (bool, error) {
	h.logger.Info("Create Ip Table rule " + ipaclConfig.Name)
	h.server.IptableAddCh <- ipaclConfig
	return true, nil
	//return (h.server.AddIpTableRule(ipaclConfig, false /* non - restart*/))
}

func (h *SYSDHandler) UpdateIpTableAcl(origConf *sysd.IpTableAcl,
	newConf *sysd.IpTableAcl, attrset []bool) (bool, error) {
	return true, nil
}

func (h *SYSDHandler) DeleteIpTableAcl(ipaclConfig *sysd.IpTableAcl) (bool, error) {
	h.logger.Info("Delete Ip Table rule " + ipaclConfig.Name)
	//return (h.server.DelIpTableRule(ipaclConfig))
	h.server.IptableDelCh <- ipaclConfig
	return true, nil
}
