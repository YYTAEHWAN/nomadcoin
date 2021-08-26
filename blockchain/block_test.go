package blockchain

import (
	"reflect"
	"testing"

	"github.com/nomadcoders/utils"
)

func TestCreateBlock(t *testing.T) {
	dbStorage = fakeDB{}
	Mempool().Txs["test"] = &Tx{}
	b := createBlock("x", 1, 1)
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("creatBlock() should return an instance of Block")
	}
}

func TestFindBlock(t *testing.T) {
	t.Run("Block not found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeGetBlockHashFromDb: func() []byte {
				t.Log("i have been called11")
				return nil
			},
		}
		t.Log("실행됨?11")
		_, err := FindBlock("xx")
		if err == nil {
			t.Error("The block should not be found")
		}
	})
	t.Run("Block is found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeGetBlockHashFromDb: func() []byte {
				b := &Block{}
				t.Log("i have been called22")
				return utils.ToBytes(b)
			},
		}
		t.Log("실행됨?22")
		block, _ := FindBlock("xx")
		if reflect.TypeOf(block) != reflect.TypeOf(&Block{}) {
			t.Error("block should be found")
		}
	})
}
