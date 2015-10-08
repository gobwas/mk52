package main

import (
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
	"log"
	"network"
)

type Task struct {
	ws *websocket.Conn
	id string
}

func (self *Task) send(response []byte) (err error) {
	if _, err = self.ws.Write(response); err != nil {
		log.Printf("[%s] error sending response: %s", self.id, err)
	}

	return
}

func (self *Task) success(result float64) (err error) {
	response, err := proto.Marshal(&network.Response{self.id, "", result})
	if err != nil {
		log.Printf("[%s] error marshalling response: %s", self.id, err)
		return
	}

	return self.send(response)
}

func (self *Task) error(e error) (err error) {
	response, err := proto.Marshal(&network.Response{self.id, e.Error(), 0})
	if err != nil {
		log.Printf("[%s] error marshalling response: %s", self.id, err)
		return
	}

	return self.send(response)
}
