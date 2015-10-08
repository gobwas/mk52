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

var Host = flag.String("host", "localhost", "host to listen")
var Port = flag.Int("port", 5555, "port to listen")
var Route = flag.String("route", "mk52", "route to listen")
var Timeout = flag.Int("timeout", 10, "timeout in seconds")

func main() {
	flag.Parse()

	origin := fmt.Sprintf("http://%s", *Host)
	url := fmt.Sprintf("ws://%s:%d/%s", *Host, *Port, *Route)

	expr := flag.Arg(1)
	if expr == "" {
		log.Fatal("expression is expected")
	}

	log.Printf("sending expression: %s", expr)

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

	timeout := time.NewTimer(time.Second * time.Duration(*Timeout))

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
