package main

import (
	"flag"
	"github.com/christopherL91/GoWebSocket/SocketServer"
	"html/template"
	"net/http"
	"runtime"
)

func serveMain(rw http.ResponseWriter, req *http.Request) {
	var template_file, _ = template.ParseFiles("client/src/views/index.html")
	req.Header.Add("Content-Type", "application/javascript")
	template_file.Execute(rw, nil)
}

func init() {
	cores := flag.Int("cores", 1, "The number of cores used")
	flag.Parse()
	runtime.GOMAXPROCS(*cores)
}

func main() {
	server := SocketServer.NewServer("/echo")
	go server.Listen()
	http.HandleFunc("/", serveMain)
	http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
		path := "client/" + r.URL.Path
		http.ServeFile(w, r, path)
	})
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
