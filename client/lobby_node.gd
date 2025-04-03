extends Node

const DEFAULT_SERVER_PORT = 9000

signal player_added(pid: int)
signal lobby_sealed(lobby_id: int)

var packed_scene = preload("lobby_node.tscn")

var peer: WebRTCMultiplayerPeer = WebRTCMultiplayerPeer.new()
var signal_ws_client: SignalWsClient = SignalWsClient.new()
var mesh_initialized := false

# UI
var Entrance: Container
var WaitingRoom: Container

var HostInput: LineEdit
var HostBtn: Button

var JoinInput: LineEdit
var JoinBtn: Button

var LobbyCode: Label
var StartGameBtn: Button

var is_host := false

func _ready():
	var scene = packed_scene.instantiate()
	add_child(scene)

	Entrance = scene.get_node("%Entrance")
	WaitingRoom = scene.get_node("%WaitingRoom")

	HostInput = scene.get_node("%HostInput")
	HostBtn = scene.get_node("%HostBtn")
	JoinInput = scene.get_node("%JoinInput")
	JoinBtn = scene.get_node("%JoinBtn")
	LobbyCode = scene.get_node("%LobbyCode")
	StartGameBtn = scene.get_node("%StartGameBtn")

	HostBtn.pressed.connect(_on_host_pressed)
	JoinBtn.pressed.connect(_on_join_pressed)
	StartGameBtn.pressed.connect(_on_start_game_pressed)

	signal_ws_client.connected.connect(_on_connected)
	signal_ws_client.lobby_hosted.connect(_on_lobby_hosted)
	signal_ws_client.lobby_joined.connect(_on_lobby_joined)
	signal_ws_client.lobby_sealed.connect(_on_lobby_sealed)
	signal_ws_client.peer_connected.connect(_on_peer_connected)

	signal_ws_client.answer_received.connect(_on_remote_description_received.bind("answer"))
	signal_ws_client.offer_received.connect(_on_remote_description_received.bind("offer"))
	signal_ws_client.candidate_received.connect(_on_candidate_received)

func _enter_waiting_room() -> void:
	StartGameBtn.set_visible(is_host)
	Entrance.hide()
	WaitingRoom.show()

func _process(_delta: float) -> void:
	signal_ws_client.poll()

func _get_server_url() -> String:
	var host = HostInput.text
	return "wss://" + host + ":" + str(DEFAULT_SERVER_PORT) + "/ws"

func _on_host_pressed():
	print("----")
	is_host = true
	signal_ws_client.connect_to_server(_get_server_url())

func _on_join_pressed():
	print("----")
	is_host = false
	var lobby_id = JoinInput.text.to_int()
	signal_ws_client.connect_to_server(_get_server_url(), lobby_id)

func _on_start_game_pressed():
	if not is_host:
		print("Only host can start the game")
		return

	signal_ws_client.seal_lobby()

func _on_connected(pid: int):
	print("[lobby] id received: ", pid, " creating mesh")
	peer.create_mesh(pid)
	multiplayer.set_multiplayer_peer(peer)
	mesh_initialized = true

	player_added.emit(pid)

func _on_lobby_hosted(pid: int, lobby_id: int):
	print("[lobby] ", peer.get_unique_id(), " lobby hosted: lobby:", lobby_id, ", peer:", pid)
	LobbyCode.text = str(lobby_id)
	_enter_waiting_room()

func _on_lobby_joined(pid: int, lobby_id: int):
	print("[lobby] ", peer.get_unique_id(), " lobby joined: lobby:", lobby_id, ", peer:", pid)
	LobbyCode.text = "Waiting..."
	_enter_waiting_room()

func _on_lobby_sealed(lobby_id: int):
	print("[lobby] ", peer.get_unique_id(), " lobby sealed:", lobby_id)
	lobby_sealed.emit(lobby_id)

func _on_remote_description_received(pid: int, sdp: String, type: String):
	if peer.has_peer(pid):
		print("[lobby] ", peer.get_unique_id(), " Setting remote description: ", type, ", ", pid)
		var peer_cn = peer.get_peer(pid).connection
		peer_cn.set_remote_description(type, sdp)

func _on_candidate_received(pid: int, mid: String, index: int, sdp: String):
	if peer.has_peer(pid):
		print("[lobby] ", peer.get_unique_id(), " Adding ice candidate: ", pid, ", ", mid, ", ", index, ", ", sdp)
		var peer_cn = peer.get_peer(pid).connection
		peer_cn.add_ice_candidate(mid, index, sdp)

func _on_peer_connected(pid: int):
	print("[lobby] ", peer.get_unique_id(), " peer connected: ", pid)

	var peer_cn: WebRTCPeerConnection = WebRTCPeerConnection.new()
	peer_cn.initialize()

	# Not needed for LAN only gaming. TBD whether I'll add WAN support.
	# peer_cn.initialize({
	# 	"iceServers": [ { "urls": ["stun:stun.l.google.com:19302"] } ]
	# })

	peer_cn.session_description_created.connect(
		func _on_session_description_created(type: String, sdp: String):
			if not peer.has_peer(pid):
				return

			peer_cn.set_local_description(type, sdp)

			match type:
				"offer":
					signal_ws_client.send_offer(pid, sdp)
				"answer":
					signal_ws_client.send_answer(pid, sdp)
				_:
					print("[lobby] ", peer.get_unique_id(), " Unknown session description type: ", type)
	)

	peer_cn.ice_candidate_created.connect(
		func _on_ice_candidate_created(mid: String, index: int, sdp: String):
			signal_ws_client.send_candidate(pid, mid, index, sdp)
	)

	print("[lobby] ", peer.get_unique_id(), " adding peer: ", pid)
	peer.add_peer(peer_cn, pid)

	player_added.emit(pid)

	# HostBtn shouldn't create offers
	if not is_host:
		print("[lobby] ", peer.get_unique_id(), " creating offer for: ", pid)
		peer_cn.create_offer()
