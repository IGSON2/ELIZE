package p2p

import (
	"elizebch/elizebch"
	"elizebch/elizeutils"
	"encoding/json"
	"fmt"
)

type messageKind int

const (
	MessageNewestBlock messageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
)

type Message struct {
	Kind    messageKind
	Payload []byte
}

func makeMessage(kind messageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: elizeutils.ToJSON(payload),
	}
	return elizeutils.ToJSON(m)
}

func sendNewstBlock(p *peer) {
	b, err := elizebch.FindBlock(elizebch.GetBlockchain().NewestHash)
	elizeutils.Errchk(err)
	blockJson := makeMessage(MessageNewestBlock, b)
	p.inbox <- blockJson
}

func handleMessage(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		var transferBlock elizebch.Block
		elizeutils.Errchk(json.Unmarshal(m.Payload, &transferBlock))
		fmt.Println(transferBlock)
	}
}
