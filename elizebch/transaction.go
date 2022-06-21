package elizebch

import (
	"elize/elizeutils"
	"elize/wallet"
	"errors"
	"sync"
	"time"
)

const (
	minerReward         int = 50
	NotEnoughBalanceErr     = "not enough balance"
	NotVerified             = "not verified"
)

type Tx struct {
	ID        string   `json:"id"`
	TimeStamp int64    `json:"timestamp"`
	TxIns     []*TxIn  `json:"txins"`
	TxOuts    []*TxOut `json:"txouts"`
}

type TxIn struct {
	TXID      string `json:"id"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

type UTxOut struct {
	TXID    string  `json:"id"`
	Index   int     `json:"index"`
	Balance float64 `json:"balance"`
}

type Mempool struct {
	Txs map[string]*Tx `json:"tx"`
	m   sync.Mutex
}

var m *Mempool
var memOnce sync.Once

func ElizeMempool() *Mempool {
	memOnce.Do(func() {
		m = &Mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

func (t *Tx) getId() {
	t.ID = elizeutils.Hash(t)
}

func CoinbaseTx(miner string) *Tx {
	newTx := Tx{
		ID:        "",
		TimeStamp: time.Now().Unix(),
		TxIns: []*TxIn{
			{"", -1, "COINBASE"},
		},
		TxOuts: []*TxOut{
			{Address: miner, Balance: float64(minerReward)},
		},
	}
	newTx.getId()
	return &newTx
}

func UTxOutsByAddress(address string) []*UTxOut {
	var uTxOuts []*UTxOut
	var spentTxOuts = make(map[string]bool)
	for _, tx := range AllTxs() {
		for _, txIn := range tx.TxIns {
			if txIn.Signature == "COINBASE" {
				break
			}
			if FindTxs(txIn.TXID).TxOuts[txIn.Index].Address == address {
				spentTxOuts[txIn.TXID] = true
			}
		}
		for index, txOut := range tx.TxOuts {
			if txOut.Address == address {
				if _, exist := spentTxOuts[tx.ID]; !exist {
					utxout := &UTxOut{tx.ID, index, txOut.Balance}
					if !isOnMempool(utxout) {
						uTxOuts = append(uTxOuts, utxout)
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
		if total >= amount {
			break
		}
		txIns = append(txIns, &TxIn{uTxOut.TXID, uTxOut.Index, ""})
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
	tx.sign()
	verify := txVerify(tx)
	if !verify {
		return nil, errors.New(NotVerified)
	}
	return tx, nil
}

func (m *Mempool) AddTxs(to string, amount float64) (*Tx, error) {
	memTX, err := makeTxs(wallet.Wallet().Address, to, amount)

	if err != nil {
		return nil, err
	}
	m.Txs[memTX.ID] = memTX
	return memTX, nil
}

func isOnMempool(utxout *UTxOut) bool {
	exist := false
Outer:
	for _, tx := range ElizeMempool().Txs {
		for _, txin := range tx.TxIns {
			if txin.TXID == utxout.TXID && txin.Index == utxout.Index {
				exist = true
				break Outer
			}
		}
	}
	return exist
}

func (m *Mempool) AllMemTx() []Tx {
	var txs []Tx
	for _, tx := range m.Txs {
		txs = append(txs, *tx)
	}

	return txs
}

func AllTxs() []*Tx {
	var Txs []*Tx
	for _, block := range AllBlock() {
		Txs = append(Txs, block.Transactions...)
	}
	return Txs
}

func FindTxs(txID string) *Tx {
	for _, tx := range AllTxs() {
		if txID == tx.ID {
			return tx
		}
	}
	return nil
}

func (t *Tx) sign() {
	for _, txin := range t.TxIns {
		txin.Signature = wallet.Sign(wallet.Wallet(), t.ID)
	}
}

func txVerify(t *Tx) bool {
	verify := false
	for _, txIn := range t.TxIns {
		prevTx := FindTxs(txIn.TXID)
		if prevTx == nil {
			verify = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		verify = wallet.Verify(txIn.Signature, t.ID, address)
		if verify {
			break
		}
	}
	return verify
}

func (m *Mempool) AddPeerTx(t *Tx) {
	m.m.Lock()
	defer m.m.Unlock()
	m.Txs[t.ID] = t
}
