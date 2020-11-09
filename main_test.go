package main

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func Test_ServerStart(t *testing.T) {
	u := url.URL{Scheme: "ws", Host: "localhost:8088", Path: ""}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05"))); err != nil {
			log.Println(err)
		}
		n, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		fmt.Printf("Received: %d.\n", n)
		fmt.Printf("Received: %s.\n", message)
		time.Sleep(5 * time.Second)
	}
}
