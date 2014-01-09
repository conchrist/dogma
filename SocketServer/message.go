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
