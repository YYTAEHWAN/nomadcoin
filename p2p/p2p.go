package p2p

import (
	"fmt"
	"net/http"

	"learngo/github.com/nomadcoders/blockchain"
	"learngo/github.com/nomadcoders/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	ip := utils.Spliter(r.RemoteAddr, ":", 0)
	openPort := r.URL.Query().Get("openPort")
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return (ip != "" && openPort != "")
	}
	fmt.Printf("r.RemoteAddr : %s\n", r.RemoteAddr)
	fmt.Printf("Upgrade에서 받은 ip : %s\n", ip)
	fmt.Printf("%s click button to connect with you and wants an upgrade\n", openPort)
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	fmt.Printf("A : (upgrade 함수)\n")
	initPeer(conn, ip, openPort[1:]) // 3000이 요청한 4000의 주소를 얻어서 서버 연결하고 peer페이지에 등록
}

func AddPeers(address string, port, openPort string, broadcast bool) {
	// address 와 port는 POST로 받은 데이터 ,, openPort는 주소에 적혀있는 port 번호
	fmt.Printf("%s wants to connect to %s\n", openPort, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	peer := initPeer(conn, address, port) // 연결되면 ,, 이 peer는 요청받은 데이터 peer 주소야
	if !broadcast {                       // 처음이니? 그럼 친구들에게 알려주렴
		BroadcastNewPeer(peer) // 새로운 Peer 다른 Peer들에게 소개시켜주고
	}
	sendNewestBlock(peer) // 새로운 Peer에겐 최신 블록을 보내본다
}

func BroadcastNewBlock(b *blockchain.Block) {
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}

func BroadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.v {
		if key != newPeer.key { // 새로 연결된 peer가 아니면
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p) // 새로 연결된 peer 알리기
		}
	}
}
