package elizebch

import (
	"crypto/sha256"
	"elizebch/database"
	"elizebch/elizeutils"
	"errors"
	"fmt"
	"strconv"
)

type Block struct {
	Height   int    `json:"height"`
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevhash,omitempty"`
}

func createBlock(inputData, newestHash string, height int) *Block {
	newBlock := Block{elize.Height + 1, inputData, "", ""}
	newBlock.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(newBlock.Data+newBlock.PrevHash+strconv.Itoa(newBlock.Height))))
	newBlock.PrevHash = newestHash
	database.SaveBlock(newBlock.Hash, elizeutils.ToBytes(newBlock))
	return &newBlock
}

func FindBlock(hash string) (*Block, error) {
	var EmptyBlock = &Block{}
	data := database.OneBlock(hash)
	if data == nil {
		return nil, errors.New("this block doesn't exist")
	} else {
		elizeutils.FromBytes(EmptyBlock, []byte(hash))
		return EmptyBlock, nil
	}
}

func AllBlock() []*Block {
	newestBlock, err := FindBlock(elize.NewestHash)
	elizeutils.Errchk(err)
	var allBlocks []*Block
	for err != nil {
		newestBlock, err = FindBlock(newestBlock.PrevHash)
		allBlocks = append(allBlocks, newestBlock)
	}
	for _, tinyblock := range allBlocks {
		fmt.Println(tinyblock)
	}
	return allBlocks
}
