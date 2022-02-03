package elizebch

import (
	"elizebch/database"
	"elizebch/elizeutils"
	"errors"
	"fmt"
	"strings"
	"time"
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
		Difficulty: GetBlockchain().Difficulty(),
	}
	fmt.Println("Difficulty in createBlock() : ", b.Difficulty)
	b.mine()
	database.SaveBlock(b.Hash, elizeutils.ToBytes(b))
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	var hashedBlock string
	for !strings.HasPrefix(hashedBlock, target) {
		hashedBlock = elizeutils.Hash(b)
		b.Nonce++
		fmt.Println("Mining Hash : ", hashedBlock)
	}
	fmt.Println("Block in mine()", b)
	b.TimeStamp = int64(time.Now().Unix())
	b.Hash = hashedBlock
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
	fmt.Println("Block in AllBlock()", newestBlock)
	elizeutils.Errchk(err)
	var allBlocks = []*Block{newestBlock}
	for {
		newestBlock, err = FindBlock(newestBlock.PrevHash)
		if err != nil {
			break
		}
		allBlocks = append(allBlocks, newestBlock)
	}
	return allBlocks
}
