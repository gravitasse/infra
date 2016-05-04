package server

import (
	"models"
	"sysd"
	"utils/dbutils"
)

func (server *SYSDServer) ReadIpAclConfigFromDB(dbHdl *dbutils.DBUtil) error {
	server.logger.Info("Reading Ip Acl Config From Db")
	if dbHdl != nil {
		var dbObj models.IpTableAcl
		objList, err := dbHdl.GetAllObjFromDb(dbObj)
		if err != nil {
			server.logger.Err("DB query failed for IpTableAcl config")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := sysd.NewIpTableAcl()
			dbObject := objList[idx].(models.IpTableAcl)
			models.ConvertsysdIpTableAclObjToThrift(&dbObject, obj)
			server.AddIpTableRule(obj, true /*restart*/)
		}
	}
	server.logger.Info("reading ip acl config done")
	return nil
}

func (server *SYSDServer) AddIpTableRule(ipaclConfig *sysd.IpTableAcl,
	restart bool) (bool, error) {
	return (server.sysdIpTableMgr.AddIpRule(ipaclConfig, restart))
}

func (server *SYSDServer) DelIpTableRule(ipaclConfig *sysd.IpTableAcl) (bool, error) {
	return (server.sysdIpTableMgr.DelIpRule(ipaclConfig))
}
