package server

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"sysd"
)

func (server *SYSDServer) ReadIpAclConfigFromDB(dbHdl *sql.DB) {
	server.logger.Info("Reading Ip Acl Config From Db")
	dbCmd := "select * from IpTableAcl"
	rows, err := dbHdl.Query(dbCmd)
	if err != nil {
		server.logger.Err(fmt.Sprintln("query to db failed", err))
		server.dbUserCh <- 1
		return
	}
	for rows.Next() {
		var config sysd.IpTableAcl
		err = rows.Scan(&config.Name, &config.PhysicalPort,
			&config.Action, &config.IpAddr, &config.Protocol,
			&config.Port)
		if err != nil {
			server.logger.Err(fmt.Sprintln("scanning rows failed", err))
		} else {
			server.AddIpTableRule(&config, true /*restart*/)
		}
	}
	server.logger.Info("reading ip acl config done")
	server.dbUserCh <- 1
}

func (server *SYSDServer) AddIpTableRule(ipaclConfig *sysd.IpTableAcl,
	restart bool) (bool, error) {
	return (server.sysdIpTableMgr.AddIpRule(ipaclConfig, restart))
}

func (server *SYSDServer) DelIpTableRule(ipaclConfig *sysd.IpTableAcl) (bool, error) {
	return (server.sysdIpTableMgr.DelIpRule(ipaclConfig))
}
