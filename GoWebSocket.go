package main

import (
	"github.com/christopherL91/GoWebSocket/SocketServer"
	"runtime"
)

func init() {
	//use maximum number of available processors.
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	//start the whole thing.
	SocketServer.StartServer()
}
