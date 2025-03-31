package pkg

import (
	"fmt"
)

type Hub struct {
	peers      map[PeerID]*Peer
	lobbies    map[LobbyID]*Lobby
	peer_lobby map[PeerID]LobbyID

	// Channels for messages
	connect    chan *Peer
	disconnect chan *Peer
	peer_msg   chan *PeerMsg
}

func NewHub() *Hub {
	return &Hub{
		peers:      make(map[PeerID]*Peer),
		lobbies:    make(map[LobbyID]*Lobby),
		peer_lobby: make(map[PeerID]LobbyID),

		connect:    make(chan *Peer),
		disconnect: make(chan *Peer),
		peer_msg:   make(chan *PeerMsg),
	}
}

// Handle hub messages
func (hub *Hub) Run() {
	fmt.Println("[Hub] Run")

	defer func() {
		fmt.Println("[Hub] Run exiting")
	}()

	for {
		select {
		case peer := <-hub.connect:
			fmt.Println("[Hub] <- connect")
			hub.peers[peer.id] = peer

		case peer := <-hub.disconnect:
			fmt.Println("[Hub] <- disconnect")
			delete(hub.peers, peer.id)

		case peer_msg := <-hub.peer_msg:
			source_peer := hub.peers[peer_msg.sourceId]
			if source_peer == nil {
				fmt.Println("[Hub] Peer not found")
				continue
			}

			switch peer_msg.msg.msgType {
			case HOST:
				source_peer.send <- msg(int(source_peer.id), CONNECTED, nil)

				lobby := NewLobby(source_peer)
				hub.lobbies[lobby.id] = lobby

				source_peer.send <- msg(int(lobby.id), HOST, nil)

				lobby.members[source_peer.id] = source_peer
				hub.peer_lobby[source_peer.id] = lobby.id

			case JOIN:
				source_peer.send <- msg(int(source_peer.id), CONNECTED, nil)

				lobby := hub.lobbies[LobbyID(peer_msg.msg.id)]

				source_peer.send <- msg(int(lobby.id), JOIN, nil)

				for _, member := range lobby.members {
					source_peer.send <- msg(int(member.id), PEER_CONNECT, nil)
					member.send <- msg(int(source_peer.id), PEER_CONNECT, nil)
				}

				lobby.members[source_peer.id] = source_peer
				hub.peer_lobby[source_peer.id] = lobby.id

			case OFFER, ANSWER, CANDIDATE:
				target_id := peer_msg.msg.id

				peer := hub.peers[PeerID(target_id)]

				if peer == nil || hub.peer_lobby[peer.id] != hub.peer_lobby[source_peer.id] {
					fmt.Println("[Hub] Peer not found or not in same lobby")
					continue
				}

				peer.send <- msg(int(source_peer.id), peer_msg.msg.msgType, peer_msg.msg.data)
			}
		}
	}
}
