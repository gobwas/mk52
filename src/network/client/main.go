package main

import (
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"network"
	"os"
	"time"
)

func main() {
	flag.Parse()

	origin := "http://localhost"
	url := "ws://localhost:5555/mk52"

	expr := flag.Arg(1)
	if expr == "" {
		log.Fatal("Pass expression")
	}

	log.Printf("Sending expression: %s", expr)

	ws, err := websocket.Dial(url, "", origin)
	defer ws.Close()
	if err != nil {
		log.Fatal(err)
	}

	uid := uuid.NewV4().String()

	msg, err := proto.Marshal(&network.Request{uid, expr})
	if err != nil {
		log.Fatalf("error marshalling: %s", err)
	}

	if _, err := ws.Write(msg); err != nil {
		log.Fatal(err)
	}

	timeout := time.NewTimer(time.Second * 5)

	for {
		select {

		// receive done request
		case <-timeout.C:
			log.Fatalf("timeout exceeded")
			return

		default:
			var msg []byte
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				switch err {
				case io.EOF:
					log.Fatalf("server closed connection: %s", err)
				default:
					log.Fatalf("error receiving: %s", err)
				}
			}

			var resp network.Response
			err = proto.Unmarshal(msg, &resp)
			if err != nil {
				log.Fatalf("error unmarshalling %s", err)
			}

			log.Printf("received response", resp)

			if resp.Error != "" {
				log.Fatalf("received error: %s", resp.Error)
			} else {
				fmt.Println(resp.Result)
				os.Exit(0)
			}
		}
	}
}
