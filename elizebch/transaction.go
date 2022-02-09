package elizebch

import (
	"elizebch/elizeutils"
	"errors"
	"time"
)

const (
	minerReward         int = 50
	NotEnoughBalanceErr     = "not enough balance"
)

type Tx struct {
	ID        string   `json:"id"`
	TimeStamp int64    `json:"timestamp"`
	TxIns     []*TxIn  `json:"txins"`
	TxOuts    []*TxOut `json:"txouts"`
}

type TxIn struct {
	ID    string `json:"id"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner   string  `json:"owner"`
	Balance float64 `json:"balance"`
}

type UTxOut struct {
	ID      string  `json:"id"`
	Index   int     `json:"index"`
	Balance float64 `json:"balance"`
}

type Mempool struct {
	Txs []*Tx `json:"tx"`
}

var ElizeMempool = Mempool{}

func (t *Tx) getId() {
	t.ID = elizeutils.Hash(t)
}

func CoinbaseTx(miner string) *Tx {
	var newTx *Tx = &Tx{
		TimeStamp: time.Now().Unix(),
		TxIns: []*TxIn{
			{"", -1, "COINBASE"},
		},
		TxOuts: []*TxOut{
			{Owner: miner, Balance: float64(minerReward)},
		},
	}
	newTx.getId()
	return newTx
}

func UTxOutsByAddress(address string) []*UTxOut {
	var uTxOuts []*UTxOut
	var spentTxOuts = make(map[string]bool)
	for _, block := range AllBlock() {
		for _, tx := range block.Transaction {
			for _, txIn := range tx.TxIns {
				if txIn.Owner == address {
					spentTxOuts[txIn.ID] = true
				}
			}
			for index, txOut := range tx.TxOuts {
				if txOut.Owner == address {
					if _, exist := spentTxOuts[tx.ID]; !exist {
						utxout := &UTxOut{tx.ID, index, txOut.Balance}
						if !isOnMempool(utxout) {
							uTxOuts = append(uTxOuts, utxout)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func BalanceByAddress(address string) float64 {
	var UnspendBalance float64
	uTxOuts := UTxOutsByAddress(address)
	for _, uTxOut := range uTxOuts {
		UnspendBalance += uTxOut.Balance
	}
	return UnspendBalance
}

func makeTxs(from, to string, amount float64) (*Tx, error) {
	if BalanceByAddress(from) < amount {
		return nil, errors.New(NotEnoughBalanceErr)
	}
	var (
		txOuts []*TxOut
		txIns  []*TxIn
		total  float64
	)

	uTxOuts := UTxOutsByAddress(from)

	for _, uTxOut := range uTxOuts {
		txIns = append(txIns, &TxIn{uTxOut.ID, uTxOut.Index, from})
		total += uTxOut.Balance
	}
	if change := (total - amount); change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOuts = append(txOuts, &TxOut{to, amount})

	tx := &Tx{
		TimeStamp: time.Now().Unix(),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *Mempool) AddTxs(to string, amount float64) error {
	memTX, err := makeTxs("igson", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, memTX)
	return nil
}

func isOnMempool(utxout *UTxOut) bool {
	exist := false
Outer:
	for _, tx := range ElizeMempool.Txs {
		for _, txin := range tx.TxIns {
			if txin.ID == utxout.ID && txin.Index == utxout.Index {
				exist = true
				break Outer
			}
		}
	}
	return exist
}
