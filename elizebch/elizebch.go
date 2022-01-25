package elizebch

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"sync"
)

type Block struct {
	Height   int
	Data     string
	Hash     string
	PrevHash string
}

type blockchain struct {
	Blocks []*Block
}

var (
	elize    *blockchain
	syncOnce sync.Once
)

func GetBlockchain() *blockchain {
	if elize == nil {
		syncOnce.Do(func() {
			elize = &blockchain{}
			elize.AddBlock("GENESIS")
		})
	}
	return elize
}

func (b *blockchain) AddBlock(inputData string) {
	b.Blocks = append(b.Blocks, createBlock(inputData))
}

func createBlock(inputData string) *Block {
	newBlock := Block{len(elize.Blocks) + 1, inputData, "", ""}
	newBlock.hashingData()
	newBlock.getLastHash()
	return &newBlock
}

func (b *Block) hashingData() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash + strconv.Itoa(b.Height)))
	b.Hash = fmt.Sprintf("%x", hash)
}

func (b *Block) getLastHash() {
	if b.Height > 1 {
		b.PrevHash = elize.Blocks[b.Height-2].Hash
	} else {
		b.PrevHash = ""
	}
}
