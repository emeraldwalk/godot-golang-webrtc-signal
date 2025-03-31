package pkg

type LobbyID int

type Lobby struct {
	id      LobbyID
	host    PeerID
	members map[PeerID]*Peer
}

func NewLobby(host *Peer) *Lobby {
	return &Lobby{
		id:      LobbyID(NewID()),
		host:    host.id,
		members: make(map[PeerID]*Peer),
	}
}
