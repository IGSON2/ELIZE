package elizebch

import (
	"elizebch/database"
	"elizebch/elizeutils"
	"sync"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var (
	elize    *blockchain
	syncOnce sync.Once
)

func GetBlockchain() *blockchain {
	if elize == nil {
		syncOnce.Do(func() {
			elize = &blockchain{}
			lastPoint := database.LastBlockPoint()
			if lastPoint == nil {
				elize.AddBlock("GENESIS")
			} else {
				elize.restore(lastPoint)
			}
		})
	}
	return elize
}

func (b *blockchain) restore(lastPoint []byte) {
	elizeutils.FromBytes(b, lastPoint)
}

func (b *blockchain) AddBlock(inputData string) {
	block := createBlock(inputData, b.NewestHash, b.Height)
	b.NewestHash = block.Hash
	b.Height = block.Height
	database.SaveBlockchain(elizeutils.ToBytes(b))
}
