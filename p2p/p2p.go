package p2p

import (
	"elizebch/elizeutils"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	_, err := upgrader.Upgrade(rw, r, nil)
	elizeutils.Errchk(err)
}

func Addpeer() {

}
