package srvfrm

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type serverOnUnixSocket struct {
	http.Server
	UnixSocket string
}

func (srv *serverOnUnixSocket) Start() {
	os.Remove(srv.UnixSocket)

	unixListener, err := net.Listen("unix", srv.UnixSocket)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.Chmod(srv.UnixSocket, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Starting server on unix socket: %s\n", srv.UnixSocket)

	go func() {
		if err := srv.Serve(unixListener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen on unix socket error: %s\n", err)
		}
	}()
}

func (srv *serverOnUnixSocket) Stop(ctx context.Context) {
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server unix socket shutdown error: %s\n", err)
	}
}

type serverOnTCPPort struct {
	http.Server
	ListenAddr string
	Port       int
}

func (srv *serverOnTCPPort) Start() {
	srv.Addr = fmt.Sprintf("%s:%d", srv.ListenAddr, srv.Port)

	log.Printf("Starting server on TCP port: %s\n", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen on TCP port error: %s\n", err)
		}
	}()
}

func (srv *serverOnTCPPort) Stop(ctx context.Context) {
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server tcp port shutdown error: %s\n", err)
	}
}

func (srv *SrvFrm) runServer() {
	r, err := srv.loadRouter()
	if err != nil {
		log.Fatalf("Could not create new router: %s\n", err)
	}

	ctx := context.Background()

	var srvUnixSocket *serverOnUnixSocket
	if srv.cfg.Server.UnixSocket != "" {
		srvUnixSocket = &serverOnUnixSocket{
			UnixSocket: srv.cfg.Server.UnixSocket,
		}
		srvUnixSocket.Handler = r
		srvUnixSocket.Start()
	}

	var srvTCPPort *serverOnTCPPort
	if srv.cfg.Server.Port != 0 {
		srvTCPPort = &serverOnTCPPort{
			ListenAddr: srv.cfg.Server.ListenAddr,
			Port:       srv.cfg.Server.Port,
		}
		srvTCPPort.Handler = r
		srvTCPPort.Start()
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown server...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if srvUnixSocket != nil {
		srvUnixSocket.Stop(ctx)
	}
	if srvTCPPort != nil {
		srvTCPPort.Stop(ctx)
	}

	log.Println("Server exiting")
}
