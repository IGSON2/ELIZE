package elizebch

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type Block struct {
	Height   int    `json:"height"`
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevhash,omitempty"`
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
	b.Blocks = append([]*Block{createBlock(inputData)}, b.Blocks...)
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
		b.PrevHash = elize.Blocks[0].Hash
	} else {
		b.PrevHash = ""
	}
}

func Allblock() []*Block {
	return GetBlockchain().Blocks
}

func FindOneblock(hash string) (*Block, error) {
	for _, block := range GetBlockchain().Blocks {
		if hash == block.Hash {
			return block, nil
		}
	}
	return nil, errors.New("this block doesn't exist")
}
