package elizebch

import (
	"reflect"
	"testing"
)

func TestCreateBlock(t *testing.T) {
	dbStorage = fakeDB{}
	ElizeMempool().Txs["Test"] = &Tx{}
	b := createBlock(&blockchain{NewestHash: "Hash", Height: 1, CurrentDifficulty: 1})
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("createBlock() should return an instance of a block")
	}
}
