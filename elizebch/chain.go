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
	m                 sync.Mutex
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

func (b *blockchain) AddBlock() *Block {
	var newBlock Block
	newBlock.createBlock(b)
	b.NewestHash = newBlock.Hash
	b.CurrentDifficulty = newBlock.Difficulty
	b.Height = newBlock.Height

	database.SaveBlockchain(elizeutils.ToBytes(b))
	return &newBlock
}

func (b *blockchain) Replace(newblocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()
	fmt.Println("Before", len(AllBlock()))
	b.CurrentDifficulty = newblocks[0].Difficulty
	b.Height = len(newblocks)
	b.NewestHash = newblocks[0].Hash
	database.SaveBlockchain(elizeutils.ToBytes(b))
	database.EmptyBlockBucket()
	for _, newblock := range newblocks {
		database.SaveBlock(newblock.Hash, elizeutils.ToBytes(newblock))
	}
	fmt.Println("After", len(AllBlock()))
}

func (b *blockchain) AddPeerBlock(block *Block) {
	b.m.Lock()
	defer b.m.Unlock()

	b.Height += 1
	b.CurrentDifficulty = block.Difficulty
	b.NewestHash = block.Hash

	database.SaveBlockchain(elizeutils.ToBytes(b))
	database.SaveBlock(block.Hash, elizeutils.ToBytes(block))

	// mempool

}
