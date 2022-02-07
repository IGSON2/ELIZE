package elizebch

import (
	"elizebch/elizeutils"
	"time"
)

const (
	minerReward int = 50
)

type Tx struct {
	Id        string   `json:"id"`
	TimeStamp int64    `json:"timestamp"`
	TxIns     []*TxIn  `json:"txins"`
	TxOuts    []*TxOut `json:"txouts"`
}

type TxIn struct {
	Owner   string  `json:"owner"`
	Balance float64 `json:"balance"`
}

type TxOut struct {
	Owner   string  `json:"owner"`
	Balance float64 `json:"balance"`
}

func (t *Tx) getId() {
	t.Id = elizeutils.Hash(t)
}

func CoinbaseTx(miner string) *Tx {
	var newTx *Tx = &Tx{
		TimeStamp: time.Now().Unix(),
		TxIns: []*TxIn{
			{Owner: "COINBASE", Balance: float64(minerReward)},
		},
		TxOuts: []*TxOut{
			{Owner: miner, Balance: float64(minerReward)},
		},
	}
	newTx.getId()
	return newTx
}

func TxOuts() []*TxOut {
	var tempTxOuts []*TxOut
	allblocks := AllBlock()
	for _, block := range allblocks {
		for _, txs := range block.Transaction {
			tempTxOuts = append(tempTxOuts, txs.TxOuts...)
		}
	}
	return tempTxOuts
}

func TxOutsByAddress(address string) []*TxOut {
	var tempTxOuts []*TxOut
	for _, txout := range TxOuts() {
		if address == txout.Owner {
			tempTxOuts = append(tempTxOuts, txout)
		}
	}
	return tempTxOuts
}

func BalanceByAddress(address string) float64 {
	var tempBalance float64
	txOutsByAddress := TxOutsByAddress(address)
	for _, txOut := range txOutsByAddress {
		if txOut.Owner == address {
			tempBalance += txOut.Balance
		}
	}
	return tempBalance
}
