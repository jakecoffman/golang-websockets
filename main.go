package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"github.com/jakecoffman/golang-websockets/chat"
	"io/ioutil"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")
var indexFile = "index.html"
var h chat.Hub

func homeHandler(response []byte) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(response)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
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
	index, err := ioutil.ReadFile(indexFile)
	if err != nil {
		panic(err)
	}
	flag.Parse()
	h = chat.NewHub()

	go h.Run()

	http.HandleFunc("/", homeHandler(index))
	http.HandleFunc("/ws", wsHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
