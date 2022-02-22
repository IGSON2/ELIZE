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

func Addpeer(ip, port, openPort string) {
	fmt.Printf("%s wants to connect to port %s\n", openPort, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", ip, port, openPort[1:]), nil)
	elizeutils.Errchk(err)
	p := initPeer(conn, ip, port)
	sendNewstBlock(p) //:4000 -> :3000
}

func BrodcastNewblock(b *elizebch.Block) {
	for _, p := range Peers.v {
		notifyNewblock(b, p)
	}
}
