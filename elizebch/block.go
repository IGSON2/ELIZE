package elizebch

import (
	"elizebch/database"
	"elizebch/elizeutils"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	difficulty       = 2
	creatingInterval = 2
)

type Block struct {
	Height     int    `json:"height"`
	Data       string `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevhash,omitempty"`
	Nonce      int    `json:"nonce"`
	Difficulty int    `json:"difficulty"`
	TimeStamp  int64  `json:"timestamp"`
}

func (b *Block) createBlock(inputData string, chain *blockchain) {
	*b = Block{
		Height:     chain.Height + 1,
		Data:       inputData,
		PrevHash:   chain.NewestHash,
		Difficulty: chain.CurrentDifficulty,
	}
	b.mine()
	database.SaveBlock(b.Hash, elizeutils.ToBytes(b))
}

func (b *Block) mine() {
	target := strings.Repeat("0", difficulty)
	var hashedBlock string
	for !strings.HasPrefix(hashedBlock, target) {
		hashedBlock = elizeutils.Hash(b)
		b.Nonce++
	}
	b.TimeStamp = int64(time.Now().Unix())
	b.Hash = hashedBlock
	fmt.Println(b)
}

func FindBlock(hash string) (*Block, error) {
	var EmptyBlock = &Block{}
	data := database.OneBlock(hash)
	if data == nil {
		return nil, errors.New("this block doesn't exist")
	} else {
		elizeutils.FromBytes(EmptyBlock, []byte(data))
		return EmptyBlock, nil
	}
}

func AllBlock() []*Block {
	newestBlock, err := FindBlock(GetBlockchain().NewestHash)
	var allBlocks = []*Block{newestBlock}
	elizeutils.Errchk(err)
	for {
		newestBlock, err = FindBlock(newestBlock.PrevHash)
		if err != nil {
			break
		}
		allBlocks = append(allBlocks, newestBlock)
	}
	return allBlocks
}
