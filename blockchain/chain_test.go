package blockchain

import (
	"reflect"
	"sync"
	"testing"

	"github.com/nomadcoders/utils"
)

type fakeDB struct {
	fakeGetCheckPointFromDb func() []byte
	fakeGetBlockHashFromDb  func() []byte
}

func (f fakeDB) GetBlockHashFromDb(hash string) []byte { // nico's FindBlock()
	return f.fakeGetBlockHashFromDb()
}
func (f fakeDB) GetCheckPointFromDb() []byte { // nico's LoadChain()
	return f.fakeGetCheckPointFromDb()
}
func (fakeDB) SaveBlock(hash string, data []byte) {} // nico's SaveBlock()
func (fakeDB) SaveCheckpoint(data []byte)         {} // nico's SaveChain()
func (fakeDB) DeleteAllBlocks()                   {} // nico's DeleteAllBlocks()

func TestBlockchain(t *testing.T) {
	t.Run("you should create blockchain", func(t *testing.T) {
		dbStorage = fakeDB{ // 여기는 가짜 DB역할
			fakeGetCheckPointFromDb: func() []byte {
				t.Log("return nil 실행완료")
				return nil
			},
		}
		bc := Blockchain()
		if bc.Height != 1 { // 테스트를 위해 Blockchain()함수가 최초 실행되면 Height = 0인 블록이 최초 생성되고
			// 블록체인 체크포인트를 가져오지 못했기 때문에 (최초생성하여 체크포인트가 없었기 때문에)
			// if nil 인 경우로 들어가 AddBlock()이 실행된다 == 여기서 Blockchain()을 호출한다
			t.Error("Blockchain() should create a blockchain 보다는 create new block 이 괜찮지 않나?")
		}
	})
	t.Run("you should restore blockchain", func(t *testing.T) {
		once = *new(sync.Once) // 이런 함수? 역할? 기능? 은 어떻게 알아내는 걸까 문서의 힘인가
		dbStorage = fakeDB{    // 여기는 가짜 DB역할
			fakeGetCheckPointFromDb: func() []byte {
				bc := &blockchain{NewstHash: "xx", Height: 2, CurrentDifficulty: 1}
				t.Log("restore fakeDB 실행됨")
				return utils.ToBytes(bc)
			},
		}
		bc := Blockchain()
		t.Log("bc := Blockchain() 실행됨")
		if bc.Height != 2 {
			t.Errorf("Blockchain Height should be %d but we got %d", 2, bc.Height)
		}
	})
}

/*
func TestBlocksSlice(t *testing.T) {
	fakeDB =
	dbStorage = fakeDB{
		fakeGetBlockHashFromDb: func() []byte {
			b := &Block{
				Height: 1,
			}
			t.Log("i have been called in BlocksSlice")
			return utils.ToBytes(b)
		},
	}

	bc := Blockchain()
	blocks := BlocksSlice(bc)
	if reflect.TypeOf(blocks) != reflect.TypeOf([]*Block{}) {
		t.Error("에러났습니다 BlockSLice")
	}
}
*/

// 그냥 따라 침
func TestBlocksSlice(t *testing.T) {
	fakeBlocks := 0
	dbStorage = fakeDB{
		fakeGetBlockHashFromDb: func() []byte {
			var b *Block
			if fakeBlocks == 0 {
				b = &Block{
					Height:   2,
					PrevHash: "xx",
				}
			}
			if fakeBlocks == 1 { // GenesisBlock에 도착했다는 의미
				b = &Block{
					Height: 1,
				}
			}
			fakeBlocks++ // 이게 반복실행 되는 것 같은데 0 -> 1 -> 2 조건 X 이도록 만든 코드인 듯
			t.Log("i have been called in BlocksSlice")
			return utils.ToBytes(b)
		},
	}
	bc := Blockchain()
	blocks := BlocksSlice(bc)
	if reflect.TypeOf(blocks) != reflect.TypeOf([]*Block{}) {
		t.Error("에러났습니다 BlockSLice")
	}
}
