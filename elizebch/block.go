package elizebch

import (
	"elize/elizeutils"
	"elize/wallet"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	Height       int    `json:"height"`
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevhash,omitempty"`
	Nonce        int    `json:"nonce"`
	Difficulty   int    `json:"difficulty"`
	TimeStamp    int64  `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

const (
	defaultDifficulty       = 2
	minuteInterval          = 2
	blockInterval           = 5
	allowedRange      int64 = 2
)

func createBlock(chain *blockchain) *Block {
	newBlock := &Block{
		Height:   chain.Height + 1,
		PrevHash: chain.NewestHash,
	}
	newBlock.setDifficulty()
	newBlock.mine()
	dbStorage.SaveBlock(newBlock.Hash, elizeutils.ToBytes(newBlock))
	return newBlock
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	var hashedBlock string
	for !strings.HasPrefix(hashedBlock, target) {
		hashedBlock = elizeutils.Hash(b)
		b.Nonce++
	}
	b.TimeStamp = int64(time.Now().Unix())
	b.ConfirmTx()
	b.Hash = hashedBlock
}

func (b *Block) setDifficulty() {
	if b.Height == 1 {
		b.Difficulty = defaultDifficulty
	} else if b.Height%blockInterval == 0 {
		b.recalculateDifficulty()
	} else {
		b.Difficulty = AllBlock()[0].Difficulty
	}
}

func (b *Block) recalculateDifficulty() {
	allblock := AllBlock()
	actualTime := (allblock[0].TimeStamp - allblock[blockInterval-2].TimeStamp) / 60
	expectedTime := int64(minuteInterval * blockInterval)
	if actualTime < expectedTime-allowedRange {
		b.Difficulty = elize.CurrentDifficulty + 1
		fmt.Println("BlockChain Difficulty has been increased.")
	} else if actualTime > expectedTime+allowedRange {
		b.Difficulty = elize.CurrentDifficulty - 1
		fmt.Println("BlockChain Difficulty has been decreased.")
	}
}

func FindBlock(hash string) (*Block, error) {
	var EmptyBlock = &Block{}
	data := dbStorage.FindBlock(hash)
	if data == nil {
		return nil, errors.New("this block doesn't exist")
	} else {
		elizeutils.FromBytes(EmptyBlock, []byte(data))
		return EmptyBlock, nil
	}
}

func AllBlock() []*Block {
	newestBlock, err := FindBlock(GetBlockchain().NewestHash)
	if err != nil {
		return []*Block{{Difficulty: defaultDifficulty}}
	}
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

func (b *Block) ConfirmTx() {
	var confirmedTx []*Tx
	for _, tx := range ElizeMempool().Txs {
		confirmedTx = append(confirmedTx, tx)
	}
	confirmedTx = append(confirmedTx, CoinbaseTx(wallet.Wallet().Address))
	ElizeMempool().Txs = make(map[string]*Tx)
	b.Transactions = confirmedTx
}
