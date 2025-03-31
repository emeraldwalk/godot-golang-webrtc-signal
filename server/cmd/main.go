package main

import (
	"flag"
	"log"
	"net/http"

	".com/emeraldwalk/learngo/pkg"
	".com/emeraldwalk/learngo/static"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":9000", "http service address")

// TODO: figure out largest expected message size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	flag.Parse()

	server := &pkg.Server{}
	server.Run()

	http.Handle("/", http.FileServer(http.FS(static.Assets)))
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test.html")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}

		server.InitPeer(ws)
	})

	log.Println("Starting server on http://localhost" + *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
