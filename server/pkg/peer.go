package pkg

import (
	"github.com/gorilla/websocket"
)

type PeerID int

type Peer struct {
	id   PeerID
	ws   *websocket.Conn
	send chan []byte
}

func NewPeer(ws *websocket.Conn) *Peer {
	return &Peer{
		id: PeerID(NewID()),
		ws: ws,
		// Buffered channel
		send: make(chan []byte, 256),
	}
}
