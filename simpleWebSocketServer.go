package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// EventData as a sample data structure which is sent to client
type EventData struct {
	EventID     int
	Description string
	IsImportant bool
}

var addr = flag.String("addr", "localhost:8081", "http service address")

var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }} // use default options

func main() {
	flag.Parse()
	log.SetFlags(0)
	fmt.Println("Starting WebSocket server")
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}
		fmt.Printf("recv: %s ", message)
		response := EventData{12343, "Some Description", true}
		marshalledResponse, err := json.Marshal(response)
		duration := time.Duration(10) * time.Second // prepare pause for 10 seconds
		for i := 0; i < 10; i++ {
			time.Sleep(duration)
			err = c.WriteMessage(mt, marshalledResponse)
			if err != nil {
				fmt.Println("write:", err)
				break
			}
		}
	}
}
