package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	chat "github.com/christopherL91/GoWebSocket2"
	"html/template"
	"net/http"
	"runtime"
)

type Server struct {
	clients      []*Client
	addClient    chan *Client
	removeClient chan *Client
	sendAll      chan *chat.MessageStruct
	messages     []*chat.MessageStruct
}

//hold all the data to manage to communicate with a client
func NewServer() *Server {
	g := Server{
		clients:      make([]*Client, 0),
		addClient:    make(chan *Client),
		removeClient: make(chan *Client),
		sendAll:      make(chan *chat.MessageStruct),
		messages:     make(chan *[]chat.MessageStruct, 0),
	}
	return &g
}

func serveMain(rw http.ResponseWriter, req *http.Request) {
	var template_file, _ = template.ParseFiles("Public/web.html")
	req.Header.Add("Content-Type", "application/javascript")
	template_file.Execute(rw, nil)
}

func websocketServer(ws *websocket.Conn) {
	ch := make(chan string)
	go read(ws)
	for {
		select {
		case c := <-ch:
			fmt.Println(c)
		default:
			//do nothing...
		}
	}
}

func init() {
	cores := flag.Int("cores", 1, "The number of cores used")
	flag.Parse()
	runtime.GOMAXPROCS(*cores)
}

func main() {
	http.Handle("/echo", websocket.Handler(websocketServer))
	http.HandleFunc("/", serveMain)
	http.HandleFunc("/Public/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
