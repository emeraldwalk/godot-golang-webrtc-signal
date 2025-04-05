package pkg

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type PeerID int

type Peer struct {
	id       PeerID
	ws       *websocket.Conn
	send     chan []byte
	isClosed bool
}

func (peer *Peer) close() {
	if peer.isClosed {
		return
	}
	fmt.Println("[Peer] Closing peer", int(peer.id))
	peer.isClosed = true
	peer.ws.Close()
}

func NewPeer(ws *websocket.Conn) *Peer {
	return &Peer{
		id: PeerID(NewID()),
		ws: ws,
		// Buffered channel
		send: make(chan []byte, 256),
	}
}
