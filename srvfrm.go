package srvfrm

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// SrvFrm keeps the details of srvfrm
type SrvFrm struct {
	Name          string
	Version       string
	DefaultConfig string

	cfg *config
	db  *sql.DB

	preDBFunc  func(*sql.DB) error
	routerFunc func(*gin.Engine) error
}

// New creates a new SrvFrm instance
func New(name, version string) *SrvFrm {
	return &SrvFrm{
		Name:          name,
		Version:       version,
		DefaultConfig: defaultConfig,
	}
}

// SetDefaultConfig sets the default config
func (srv *SrvFrm) SetDefaultConfig(defaultConfig string) {
	srv.DefaultConfig = defaultConfig
}

// SetAppConfig sets the app config
func (srv *SrvFrm) SetAppConfig(appConfig interface{}) {
	srv.cfg.App = appConfig
}

// SetPreDBFunc sets the create table function
func (srv *SrvFrm) SetPreDBFunc(f func(*sql.DB) error) {
	srv.preDBFunc = f
}

// SetRouterFunc sets the custom router function
func (srv *SrvFrm) SetRouterFunc(f func(*gin.Engine) error) {
	srv.routerFunc = f
}

func (srv *SrvFrm) printVersion() {
	fmt.Printf("%s %s (built on SrvFrm %s)\n", srv.Name, srv.Version, Version)
}

// Run the service with srvfrm
func (srv *SrvFrm) Run() {
	flagVersion := flag.Bool("version", false, "Print server version")
	flagConfigFile := flag.String("config", "config.yml", "Configuration file path")
	flagDefaultConfig := flag.Bool("default", false, "Output default configuration")

	flag.Parse()

	if *flagVersion {
		srv.printVersion()
		return
	}

	if *flagDefaultConfig {
		fmt.Printf("%s", srv.DefaultConfig)
		return
	}

	err := srv.loadConfig(*flagConfigFile)
	if err != nil {
		log.Fatalln(err)
	}

	err = srv.loadLog()
	if err != nil {
		log.Fatalln(err)
	}
	defer destroyLog()

	err = srv.loadDatabase()
	if err != nil {
		log.Fatalln(err)
	}
	defer srv.destroyDatabase()

	srv.runServer()
}
