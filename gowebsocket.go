// Package gowebsocket is a websocket library.
// Currently functional api - will improve and add OO alternative.
package gowebsocket

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

// WebsocketServer is the main struct representing the websocket server.
type WebsocketServer struct{}

// WebServer is a simple html server, useful for serving a websocket client.
type WebServer struct{}

// SocketHandler is an alias.
type SocketHandler func(*websocket.Conn)

// CustomHandler is an alias.
type CustomHandler func([]byte, func(string))

// NewWebsocketServer returns a new WebsocketServer.
func NewWebsocketServer() *WebsocketServer {
	return &WebsocketServer{}
}

// NewWebServer returns a new WebServer.
func NewWebServer() *WebServer {
	return &WebServer{}
}

// Add adds a route and an associated handler function.
func (s *WebsocketServer) Add(route string, handlerFn CustomHandler) {
	http.Handle(route, websocket.Handler(getSocketHandler(handlerFn)))
}

// Start starts the server listening for the specified routes.
func (s *WebsocketServer) Start() {
	fmt.Println("Starting websocket server...")
	// Serve an example client html page.
	w := NewWebServer()
	w.serveExampleClientPage()
	// Add builtin fixed routes - will probably remove.
	s.addBuiltinRoutes()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

// getSocketHandler returns a handler to be asscociated with a route.
// The handler wraps the client-provided handler, receives request data
// from the websocket, passes it on to the wrapped handler, along with a
// send function for the wrapped handler to use to send the response.
func getSocketHandler(myCustHandler CustomHandler) SocketHandler {
	return func(ws *websocket.Conn) {
		var in []byte
		if err := websocket.Message.Receive(ws, &in); err != nil {
			fmt.Println("Err,", err)
		}
		outFn := func(msg string) {
			websocket.Message.Send(ws, msg)
		}
		myCustHandler(in, outFn)
	}
}

/*
 * The below is non-essential and will be (re)moved at some point.
 */

// serveExampleClientPage serves a html page which will communicate with the websocket server.
func (s *WebServer) serveExampleClientPage() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "socketbasic.html")
	}
	http.HandleFunc("/", handler)
	// Putting this here for possible later use.
	//err := http.ListenAndServe(":8081", nil)
	//if err != nil {
	//	panic("ListenAndServe: " + err.Error())
	//}

}

// addBuiltinRoutes adds some routes that might be useful, eg for debugging.
func (s *WebsocketServer) addBuiltinRoutes() {
	// Http handlers - http://
	//http.HandleFunc("/", reqDump)
	//fs := http.FileServer(http.Dir("."))
	//http.Handle("/builtin/", http.StripPrefix("/builtin/", fs))

	//http.Handle("/", http.FileServer(http.Dir(".")))

	// Websocket handlers - ws://
	//http.Handle("/echo", websocket.Handler(webHandler))
	//http.Handle("/", websocket.Handler(webHandler))
}

// reqDump might be useful for debugging.
func reqDump(c http.ResponseWriter, req *http.Request) {
	fmt.Println("Received request for url:", req.URL)
	c.Write([]byte("Nice."))
}

// echoHandler might be useful for debugging.
func echoHandler(ws *websocket.Conn) {
	fmt.Println("Echoing.")
	io.Copy(ws, ws)
}

// webHandler is a sample fake data feed.
func webHandler(ws *websocket.Conn) {
	fmt.Println("rx")
	var in []byte
	if err := websocket.Message.Receive(ws, &in); err != nil {
		fmt.Println("Err,", err)
		return
	}

	for i := 0; i < 100; i++ {
		for j := 0; j < 10; j++ {
			rad := (math.Pi * 2) / 10 * float64(j)
			response := strconv.FormatFloat(rad, 'f', 6, 64)
			websocket.Message.Send(ws, response)
			time.Sleep(time.Millisecond * 100)
		}
	}
}
