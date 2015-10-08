package main

import (
	"calculator/lexer/base"
	"calculator/parser/rpn"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"network"
	"strings"
	"time"
)

var Host = flag.String("host", "localhost", "host to listen")
var Port = flag.Int("port", 5555, "port to listen")
var Route = flag.String("route", "mk52", "route to listen")
var Timeout = flag.Int("timeout", 10, "timeout in seconds")

func CalcServer(ws *websocket.Conn) {
	defer ws.Close()
	timeout := time.NewTimer(time.Second * time.Duration(*Timeout))

	for {
		select {
		case <-timeout.C:
			log.Printf("timeout exceeded for connection")
			return

		default:
			var req []byte
			err := websocket.Message.Receive(ws, &req)
			if err != nil {
				log.Printf("error receiving %s", err)
			}

			var request network.Request
			err = proto.Unmarshal(req, &request)
			if err != nil {
				log.Printf("error unmarshalling request: %s", err)
				return
			}

			log.Printf("recevied: %s", request)

			task := Task{ws, request.Id}

			calc := rpn.New(base.New(strings.NewReader(request.Expression)))
			expr, err := calc.Parse()
			if err != nil {
				log.Printf("could not parse expression: %s", err)
				task.error(err)
				return
			}

			result, err := expr.Evaluate()
			if err != nil {
				log.Printf("could not evaluate expression: %s", err)
				task.error(err)
				return
			}

			log.Printf("calculated result: %f", result)
			task.success(result)

			return
		}
	}
}

func main() {
	flag.Parse()
	http.Handle(fmt.Sprintf("/%s", *Route), websocket.Handler(CalcServer))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *Host, *Port), nil))
}
