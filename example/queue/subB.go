package main

import (
	"flag"
	"fmt"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "websocket url")
var topic = flag.String("topic", "test", "subscribe topic")
var group = flag.String("group", "subB", "group name")

func main() {
	flag.Parse()

	url := fmt.Sprintf("ws://%s/ws/sub/%s/%s", *addr, *topic, *group)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		return
	}
	defer c.Close()

	fmt.Printf("wait message for group %s\n", *group)
	for {
		t, msg, e := c.ReadMessage()
		if e != nil {
			fmt.Printf("read message error: %v", e)
			return
		}
		fmt.Printf("recv: %s(type=%d)\n", msg, t)
	}
}
