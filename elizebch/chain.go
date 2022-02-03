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

const (
	defaultDifficulty       = 2
	minuteInterval          = 2
	blockInterval           = 5
	allowedRange      int64 = 2
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
	var newBlock Block
	newBlock.createBlock(inputData, b)
	b = &blockchain{NewestHash: newBlock.Hash,
		Height: newBlock.Height,
	}
	fmt.Println("Before : ", newBlock.Difficulty)
	b.recalculateDifficulty()
	fmt.Println("After : ", newBlock.Difficulty)
	database.SaveBlockchain(elizeutils.ToBytes(b))
}

func (b *blockchain) recalculateDifficulty() {
	if b.Height == 1 {
		b.CurrentDifficulty = defaultDifficulty
	} else if b.Height%blockInterval == 0 {
		allblock := AllBlock()
		actualTime := allblock[0].TimeStamp - allblock[minuteInterval-1].TimeStamp
		expectedTime := int64(minuteInterval * blockInterval)
		if actualTime < expectedTime-allowedRange {
			b.CurrentDifficulty++
			fmt.Println("BlockChain Difficulty has been increased.")
		} else if actualTime > expectedTime+allowedRange {
			b.CurrentDifficulty--
			fmt.Println("BlockChain Difficulty has been decreased.")
		}
	}
}
