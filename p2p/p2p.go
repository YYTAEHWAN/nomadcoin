package p2p

import (
	"fmt"
	"net/http"
	"time"

	"learngo/github.com/nomadcoders/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)

	for {
		_, p, err := conn.ReadMessage() // readMessage 하는건 block이다? 무슨말이지
		// fmt.Println("message arrived!") // blocking 하는지 확인하기 위해
		if err != nil {
			break
		}
		fmt.Printf("Just got:%s\n\n", p)
		time.Sleep(5 * time.Second)
		message := fmt.Sprintf("We also think that: %s", p)
		utils.HandleErr(conn.WriteMessage(websocket.TextMessage, []byte(message)))

	}

}
