package p2p

import (
	"elizebch/elizeutils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	openPort := r.URL.Query().Get("openPort")

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	elizeutils.Errchk(err)
	result := strings.Split(r.RemoteAddr, ":")
	fmt.Println(r.RemoteAddr)
	initPeer(conn, result[0], openPort)
}

func Addpeer(ip, port, openPort string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", ip, port, openPort), nil)
	elizeutils.Errchk(err)
	initPeer(conn, ip, port)
}
