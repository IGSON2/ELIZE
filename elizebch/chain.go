package elizebch

import (
	"elize/database"
	"elize/elizeutils"
	"fmt"
	"sync"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"current_difficulty"`
	m                 sync.Mutex
}

type storage interface {
	SaveBlock(hash string, data []byte)
	SaveBlockchain(data []byte)
	LastBlockPoint() []byte
	FindBlock(hash string) []byte
	DeleteBlocks()
}

var (
	elize     *blockchain
	syncOnce  sync.Once
	dbStorage storage = database.DB{}
)

func GetBlockchain() *blockchain {
	syncOnce.Do(func() {
		elize = &blockchain{}
		lastPoint := dbStorage.LastBlockPoint()
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
	newBlock := createBlock(b)
	b.NewestHash = newBlock.Hash
	b.CurrentDifficulty = newBlock.Difficulty
	b.Height = newBlock.Height

	dbStorage.SaveBlockchain(elizeutils.ToBytes(b))
	return newBlock
}

func (b *blockchain) Replace(newblocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()
	fmt.Println("Before", len(AllBlock()))
	b.CurrentDifficulty = newblocks[0].Difficulty
	b.Height = len(newblocks)
	b.NewestHash = newblocks[0].Hash
	dbStorage.SaveBlockchain(elizeutils.ToBytes(b))
	dbStorage.DeleteBlocks()
	for _, newblock := range newblocks {
		dbStorage.SaveBlock(newblock.Hash, elizeutils.ToBytes(newblock))
	}
	fmt.Println("After", len(AllBlock()))
}

func (b *blockchain) AddPeerBlock(newblock *Block) {
	b.m.Lock()
	m.m.Lock()
	defer b.m.Unlock()
	defer m.m.Unlock()

	b.Height += 1
	b.CurrentDifficulty = newblock.Difficulty
	b.NewestHash = newblock.Hash

	dbStorage.SaveBlockchain(elizeutils.ToBytes(b))
	dbStorage.SaveBlock(newblock.Hash, elizeutils.ToBytes(newblock))

	for _, tx := range newblock.Transactions {
		_, ok := m.Txs[tx.ID]
		if ok {
			delete(m.Txs, tx.ID)
		}
	}

}
