package pkg

import "time"

type LobbyID int

type Lobby struct {
	id       LobbyID
	host     PeerID
	members  map[PeerID]*Peer
	sealedAt time.Time
}

func NewLobby(host *Peer) *Lobby {
	return &Lobby{
		id:      LobbyID(NewID()),
		host:    host.id,
		members: make(map[PeerID]*Peer),
	}
}
