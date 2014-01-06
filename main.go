package main

import (
	"flag"
	"github.com/christopherL91/GoWebSocket2/Chat"
	"html/template"
	"net/http"
	"runtime"
)

func serveMain(rw http.ResponseWriter, req *http.Request) {
	var template_file, _ = template.ParseFiles("Public/web.html")
	req.Header.Add("Content-Type", "application/javascript")
	template_file.Execute(rw, nil)
}

func init() {
	cores := flag.Int("cores", 1, "The number of cores used")
	flag.Parse()
	runtime.GOMAXPROCS(*cores)
}

func main() {
	server := Chat.NewServer("/echo")
	go server.Listen()
	http.HandleFunc("/", serveMain)
	http.HandleFunc("/Public/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
