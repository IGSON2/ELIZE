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
				fmt.Println("Init")
				elize.AddBlock("GENESIS")
			} else {
				fmt.Println("Restore")
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
	b = &blockchain{
		NewestHash:        newBlock.Hash,
		Height:            newBlock.Height,
		CurrentDifficulty: newBlock.Difficulty,
	}
	fmt.Println("BlockChain : ", b)
	database.SaveBlockchain(elizeutils.ToBytes(b))
}

func (b *blockchain) recalculateDifficulty() int {
	allblock := AllBlock()
	actualTime := (allblock[0].TimeStamp - allblock[minuteInterval-1].TimeStamp) / 60
	expectedTime := int64(minuteInterval * blockInterval)
	if actualTime < expectedTime-allowedRange {
		b.CurrentDifficulty++
		fmt.Println("BlockChain Difficulty has been increased.")
	} else if actualTime > expectedTime+allowedRange {
		b.CurrentDifficulty--
		fmt.Println("BlockChain Difficulty has been decreased.")
	}
	return b.CurrentDifficulty
}

func (b *blockchain) Difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%blockInterval == 0 {
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}
