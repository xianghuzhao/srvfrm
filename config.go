package srvfrm

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var defaultConfig = `log:
  gin_release_mode: false
  tty: true
  time_format: "2006-01-02 15:04:05.000 -0700 MST"
  server_log: server.log
  gin_log: gin.log

database:
  host: '/var/run/postgresql'
  port: ''
  user: srvfrm
  password: ''
  dbname: srvfrm
  sslmode: disable

server:
  unix_socket: ''
  listen_addr: ''
  port: 8080
  api_host: https://example.com
`

// Config keeps the configuration
type Config struct {
	Log struct {
		GinReleaseMode bool   `yaml:"gin_release_mode"`
		Tty            bool   `yaml:"tty"`
		TimeFormat     string `yaml:"time_format"`
		ServerLog      string `yaml:"server_log"`
		GinLog         string `yaml:"gin_log"`
	} `yaml:"log"`

	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
		Sslmode  string `yaml:"sslmode"`
	} `yaml:"database"`

	Server struct {
		UnixSocket string `yaml:"unix_socket"`
		ListenAddr string `yaml:"listen_addr"`
		Port       int    `yaml:"port"`
		APIHost    string `yaml:"api_host"`
	} `yaml:"server"`

	App interface{} `yaml:"app"`
}

func (srv *SrvFrm) loadConfig(fn string) error {
	srv.cfg = &Config{}

	err := srv.loadDefaultConfig()
	if err != nil {
		return err
	}

	err = srv.loadFromFile(fn)
	if err != nil {
		return err
	}

	return nil
}

func (srv *SrvFrm) loadDefaultConfig() error {
	err := yaml.Unmarshal([]byte(srv.DefaultConfig), srv.cfg)
	if err != nil {
		log.Println("Could not unmarshal default config")
		return err
	}

	return nil
}

func (srv *SrvFrm) loadFromFile(fn string) error {
	buffer, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Printf("Config file \"%s\" read error: %s\n", fn, err)
		return nil
	}

	err = yaml.Unmarshal(buffer, srv.cfg)
	if err != nil {
		return err
	}
	return nil
}
