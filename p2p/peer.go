package p2p

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

var Peers peerMap = peerMap{v: make(map[string]*peer)}

type peerMap struct {
	v map[string]*peer
	m sync.Mutex
}

type peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte
}

func AllPeers(p *peerMap) []string {
	p.m.Lock()
	defer p.m.Unlock()
	var keys []string
	for key := range p.v {
		keys = append(keys, key)
	}
	return keys
}

func initPeer(initConn *websocket.Conn, ip, port string) *peer {
	key := fmt.Sprintf("%s:%s", ip, port)
	p := &peer{
		key:     key,
		address: ip,
		port:    port,
		conn:    initConn,
		inbox:   make(chan []byte),
	}
	go p.read()
	go p.write()
	Peers.v[key] = p
	return p
}

func (p *peer) close() {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	p.conn.Close()
	delete(Peers.v, p.key)
}

func (p *peer) write() {
	defer p.close()
	for {
		m, ok := <-p.inbox
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

func (p *peer) read() {
	defer p.close()
	for {
		m := Message{}
		err := p.conn.ReadJSON(&m)
		if err != nil {
			break
		}
		handleMessage(&m, p)
	}
}
