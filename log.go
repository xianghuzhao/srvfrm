package srvfrm

import (
	"fmt"
	"log"
	"os"
)

var logFile *os.File

func (srv *SrvFrm) loadLog() error {
	log.SetPrefix(fmt.Sprintf("[%s] ", srv.Name))

	if srv.cfg.Log.Tty {
		return nil
	}

	logFile, err := os.OpenFile(srv.cfg.Log.ServerLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	log.SetOutput(logFile)

	return nil
}

func destroyLog() error {
	return logFile.Close()
}
