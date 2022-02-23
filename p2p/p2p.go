package p2p

import (
	"elizebch/elizebch"
	"elizebch/elizeutils"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	openPort := r.URL.Query().Get("openPort")
	ip := elizeutils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	fmt.Printf("%s wants an upgrade \n", r.Host)
	conn, err := upgrader.Upgrade(rw, r, nil)
	elizeutils.Errchk(err)
	initPeer(conn, ip, openPort)
}

func AddPeer(ip, port, openPort string, broadcast bool) {
	fmt.Printf("%s wants to connect to port %s\n", openPort, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", ip, port, openPort), nil)
	elizeutils.Errchk(err)
	p := initPeer(conn, ip, port)
	if broadcast {
		BroadcastNewPeer(p)
		return
	}
	sendNewstBlock(p)
}

func BroadcastNewBlock(b *elizebch.Block) {
	for _, p := range Peers.v {
		notifyNewblock(b, p)
	}
}

func BroadcastNewTx(t *elizebch.Tx) {
	for _, p := range Peers.v {
		notifyNewTx(t, p)
	}
}

func BroadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.v {
		fmt.Println("KEY : ", key)
		if len(Peers.v) == 1 || key != newPeer.key {
			address := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			fmt.Println("Notifyed Address : ", address)
			notifyNewPeer(address, p)
		}
	}
}
