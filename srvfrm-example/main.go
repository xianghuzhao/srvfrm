package main

import (
	_ "github.com/lib/pq"
	"github.com/xianghuzhao/srvfrm"
)

func main() {
	srv := srvfrm.New("MyServer", "1.2.0")
	srv.Run()
}
