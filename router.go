package srvfrm

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var identityKey = "id"

func getRoot() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiHost := "https://example.com"
		c.JSON(200, gin.H{
			"api_version": "v1",
			"user_url":    apiHost + "/users",
		})
	}
}

func serverHeader(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Server", name)

		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
	}
}

func (srv *SrvFrm) loadRouter() (*gin.Engine, error) {
	if srv.cfg.Log.GinReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	if !srv.cfg.Log.Tty && srv.cfg.Log.GinLog != "" {
		f, err := os.OpenFile(srv.cfg.Log.GinLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			gin.DefaultWriter = io.MultiWriter(f)
		} else {
			log.Printf("Failed to log to file \"%s\", using default stderr", srv.cfg.Log.GinLog)
		}
	}

	router := gin.New()

	router.Use(gin.Recovery())

	//router.Use(gin.Logger())
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(srv.cfg.Log.TimeFormat),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.Use(serverHeader(srv.Name))

	if srv.RouterFunc != nil {
		srv.RouterFunc(router)
	}

	return router, nil
}
