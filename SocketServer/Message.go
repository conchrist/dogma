package SocketServer

import (
	"encoding/json"
	"errors"
	"fmt"
)

type MessageStruct struct {
	From    string `json:"from"`
	Message string `json:"message"`
	Time    int    `json:"time"`
}

//returns message struct, decoded json.
func ParseMessage(m []byte) (*MessageStruct, error) {
	tmp := new(MessageStruct)
	err := json.Unmarshal(m, tmp)

	if err != nil {
		//return empty message
		return new(MessageStruct), errors.New("Could not parse json")
	}
	return tmp, nil
}

//returns encoded message.
func NewMessage(from string, message string, timestamp int) ([]byte, error) {

	g := MessageStruct{
		From:    from,
		Message: message,
		Time:    timestamp,
	}

	out, err := json.Marshal(g)

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("Could not create json")
	}
	return out, nil
}
