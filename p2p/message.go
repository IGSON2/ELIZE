package p2p

import (
	"elizebch/database"
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
	MessageNotifyNewBlock
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
	fmt.Printf("%s sending NewestBlock to %s\n", database.GetDBname(), p.key)
}

func requestAllblock(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func sendAllblocks(p *peer) {
	m := makeMessage(MessageAllBlocksResponse, elizebch.AllBlock())
	p.inbox <- m
}

func notifyNewblock(b *elizebch.Block, p *peer) {
	m := makeMessage(MessageNotifyNewBlock, b)
	p.inbox <- m
}

func handleMessage(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		fmt.Printf("Received the newest block from %s\n", p.key)
		var transferBlock elizebch.Block
		elizeutils.Errchk(json.Unmarshal(m.Payload, &transferBlock))
		lastBlock, err := elizebch.FindBlock(elizebch.GetBlockchain().NewestHash)
		elizeutils.Errchk(err)
		if transferBlock.Height >= lastBlock.Height {
			fmt.Printf("Requesting all blocks to %s\n", p.key)
			requestAllblock(p)
		} else {
			sendNewstBlock(p)
		}
	case MessageAllBlocksRequest:
		fmt.Printf("%s wants all the blocks.\n", p.key)
		sendAllblocks(p)

	case MessageAllBlocksResponse:
		fmt.Printf("Received all the blocks from %s\n", p.key)
		var payload []*elizebch.Block
		elizeutils.Errchk(json.Unmarshal(m.Payload, &payload))
		elizebch.GetBlockchain().Replace(payload)

	case MessageNotifyNewBlock:
		var payload *elizebch.Block
		elizeutils.Errchk(json.Unmarshal(m.Payload, &payload))
		elizebch.GetBlockchain().AddPeerBlock(payload)
	}
}
