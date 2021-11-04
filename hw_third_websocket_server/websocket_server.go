package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var err error

func main() {
	http.HandleFunc("/ws", wsHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServ: ", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	time.Sleep(time.Second)
	if err = c.WriteJSON("hello world"); err != nil {
		log.Println(err)
	}
	defer c.Close()
}
