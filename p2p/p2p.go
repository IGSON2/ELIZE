package p2p

import (
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
	conn, err := upgrader.Upgrade(rw, r, nil)
	elizeutils.Errchk(err)
	initPeer(conn, ip, openPort)
}

func Addpeer(ip, port, openPort string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", ip, port, openPort[1:]), nil)
	elizeutils.Errchk(err)
	p := initPeer(conn, ip, port)
	sendNewstBlock(p)
}
