package main

import (
	"github.com/christopherL91/GoWebSocket/SocketServer"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	SocketServer.StartServer()
}
