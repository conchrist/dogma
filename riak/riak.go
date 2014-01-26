package main

import (
	riak "github.com/mrb/riakpbc"
	"log"
	"os/exec"
)

func generateUuid() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

type Data struct {
	Message string `riak:"index" json:"data"`
}

func main() {

	coder := riak.NewCoder("json", riak.JsonMarshaller, riak.JsonUnmarshaller)
	client := riak.NewClientWithCoder([]string{"127.0.0.1:8098"}, coder)
	defer client.Close()

	// Dial all the nodes.
	if err := client.Dial(); err != nil {
		log.Fatalf("Dialing failed: %v ", err)
	}

	log.Println("Connected to nodes")

	if _, err := client.SetClientId("GoWebSocket"); err != nil {
		log.Fatal("Setting ID failed. ", err)
	}

	id := generateUuid()
	testData := Data{
		Message: "hejsan",
	}

	if _, err := client.StoreStruct("posts", id, &testData); err != nil {
		//log.Println(rres)
		log.Fatal("Failed storing struct", err)
	}

	log.Println("i'm here")

	incomming := Data{}

	if _, err := client.FetchStruct("posts", id, &incomming); err != nil {
		log.Fatal("Failed to fetch struct", err)
	}
}
