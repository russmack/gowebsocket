package main

import (
	"fmt"
	"github.com/russmack/gowebsocket"
	"strconv"
	"time"
)

func main() {
	// Create a websocket server.
	ws := gowebsocket.NewWebsocketServer()
	// Add a route and corresponding handler.
	ws.Add("/testsock", testsockHandler)
	// Start the server.
	ws.Start()
}

// testsockHandler is the function that will receive client requests and send responses.
func testsockHandler(req []byte, respFn func(string)) {
	fmt.Println("Recieved request:", string(req))
	for i := 0; i < 10; i++ {
		respFn("Sending message #" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
	}
	time.Sleep(3 * time.Second)
}
