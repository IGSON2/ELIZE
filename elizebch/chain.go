package elizebch

import (
	"elizebch/database"
	"elizebch/elizeutils"
	"fmt"
	"sync"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"current_difficulty"`
}

var (
	elize    *blockchain
	syncOnce sync.Once
)

func GetBlockchain() *blockchain {
	syncOnce.Do(func() {
		elize = &blockchain{}
		lastPoint := database.LastBlockPoint()
		if lastPoint == nil {
			fmt.Println("Init")
			elize.AddBlock()
		} else {
			fmt.Println("Restore")
			elize.restore(lastPoint)
		}
	})
	return elize
}

func (b *blockchain) restore(lastPoint []byte) {
	elizeutils.FromBytes(b, lastPoint)
}

func (b *blockchain) AddBlock() {
	var newBlock Block
	newBlock.createBlock(b)
	b.NewestHash = newBlock.Hash
	b.CurrentDifficulty = newBlock.Difficulty
	b.Height = newBlock.Height

	database.SaveBlockchain(elizeutils.ToBytes(b))
}
