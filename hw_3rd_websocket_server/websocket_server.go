package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var err error

func main() {
	http.HandleFunc("/", page_render)
	http.HandleFunc("/ws", ws_connect)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("ListenAndServ failed: ", err)
	}
}

func page_render(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./render.html")
	if err != nil {
		log.Fatalf("template.ParseFiles failed: ", err)
	}
	t.Execute(w, nil)

}

func ws_connect(w http.ResponseWriter, r *http.Request) {
	//判斷請求是否為websocket升級請求。
	if websocket.IsWebSocketUpgrade(r) {
		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			log.Fatalf("websocket upgrader.Upgrade failed: ", err)
		}
		conn.WriteMessage(websocket.TextMessage, []byte("The web socket is connected."))
		// 使用goroutine接收與回覆訊息
		go func() {
			for {
				t, c, err := conn.ReadMessage()
				if err != nil {
					// panic(err)
					// log.Fatalf("websocket conn.ReadMessage() failed: ", err)
					return
				}
				received := fmt.Sprintf("Message: \"%s\" received.", string(c))
				conn.WriteMessage(websocket.TextMessage, []byte(received))
				if t == -1 {
					return
				}
			}
		}()
	} else {
		fmt.Println("not connected")
	}
}
