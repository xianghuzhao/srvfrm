package srvfrm

import (
	"log"

	"database/sql"
)

func (srv *SrvFrm) loadDatabase(createTableFunc func(*sql.DB) error) error {
	cfgDb := srv.cfg.Database

	connStr := "host=" + cfgDb.Host
	if cfgDb.Port != "" {
		connStr += " port=" + cfgDb.Port
	}
	if cfgDb.User != "" {
		connStr += " user=" + cfgDb.User
	}
	if cfgDb.Password != "" {
		connStr += " password=" + cfgDb.Password
	}
	if cfgDb.Port != "" {
		connStr += " port=" + cfgDb.Port
	}
	if cfgDb.Dbname != "" {
		connStr += " dbname=" + cfgDb.Dbname
	}
	if cfgDb.Sslmode != "" {
		connStr += " sslmode=" + cfgDb.Sslmode
	}

	log.Printf("Connect to database: %s", connStr)

	var err error

	srv.db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect to database (%s): %s", connStr, err)
		return err
	}

	err = srv.db.Ping()
	if err != nil {
		return err
	}

	err = createTableFunc(srv.db)
	if err != nil {
		return err
	}

	return nil
}

func (srv *SrvFrm) destroyDatabase() error {
	log.Printf("Close database")
	return srv.db.Close()
}
