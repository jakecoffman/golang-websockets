package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"github.com/jakecoffman/golang-websockets/chat"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8000", "http service address")
var indexFile = "index.html"
var h chat.Hub

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a websocket. TODO: check origin
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}
	c := chat.NewConnection(ws)
	h.Register(c)
	defer func() { h.Unregister(c) }()
	go c.Writer()
	c.Reader(h)
}

func main() {
	flag.Parse()
	h = chat.NewHub()

	go h.Run()

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", wsHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
