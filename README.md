# SrvFrm

A very simple go HTTP API server framework.
Simplify the construction of a HTTP API server, mostly for my personal use.


## Example

```go
import (
	_ "github.com/lib/pq"
	"github.com/xianghuzhao/srvfrm"
)

func main() {
	srv := srvfrm.New("MyServer", "1.2.0")
	srv.Run()
}
```


## HTTP Server

```go
import "net/http"
```


## Gin

<https://github.com/gin-gonic/gin>


## Log

```go
import "log"
```


## Postgresql

```go
import "database/sql"
import _ "github.com/lib/pq"
```
