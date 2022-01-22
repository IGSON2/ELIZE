package blockchainn

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"sync"
)

type block struct {
	Height   int
	Data     string
	Hash     string
	PrevHash string
}

type blockchain struct {
	Blocks []*block
}

var (
	elize    *blockchain
	syncOnce sync.Once
)

func GetBlockchain() *blockchain {
	if elize == nil {
		syncOnce.Do(func() {
			elize.addBlock("GENESIS")
		})
	}
	elize.addBlock("SECOND")
	return elize
}

func (b *blockchain) addBlock(inputData string) {
	b.Blocks = append(b.Blocks, createBlock(inputData))
}

func createBlock(inputData string) *block {
	newBlock := block{len(elize.Blocks) + 1, inputData, "", ""}
	newBlock.hashingData()
	newBlock.getLastHash()
	return &newBlock
}

func (b *block) hashingData() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash + strconv.Itoa(b.Height)))
	b.Hash = fmt.Sprintf("%x", hash)
}

func (b *block) getLastHash() {
	if b.Height > 1 {
		b.PrevHash = elize.Blocks[b.Height-1].Hash
	} else {
		b.PrevHash = ""
	}
}
