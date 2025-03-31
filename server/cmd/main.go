package main

import (
	"flag"
	"log"
	"net/http"

	".com/emeraldwalk/learngo/pkg"
	".com/emeraldwalk/learngo/static"
	"github.com/gorilla/websocket"
)

var port = flag.String("port", ":9000", "http port")
var isHttps = flag.Bool("https", false, "use https")

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

	protocol := "http"
	if *isHttps {
		protocol = "https"
	}

	log.Printf("Starting server on %s://localhost%s", protocol, *port)
	var err error
	if *isHttps {
		err = http.ListenAndServeTLS(*port, "ssl.crt", "ssl.key", nil)
	} else {
		err = http.ListenAndServe(*port, nil)
	}

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
