package main

import (
	"flag"
	"github.com/conchrist/dogma/SocketServer"
	"runtime"
)

var config string

//START OMIT
func init() {
	//use maximum number of available processors.
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&config, "config", "config.gcfg", "Specify path to config file")
}

func main() {
	//start the whole thing.
	SocketServer.StartServer(config)
}

//END OMIT
