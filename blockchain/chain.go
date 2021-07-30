package blockchain

import (
	"learngo/github.com/nomadcoders/db"
	"learngo/github.com/nomadcoders/utils"
	"sync"
)

const (
	defaultDifficulty  int = 2 // 첫 블록 생성의 난이도는 2로 초기화
	difficultyInterval int = 5 // 블록 간격 5개에 한번씩 난이도 조절
	blockInterval      int = 2 // 블록 생성 간격 2분
	allowedRange       int = 2 // 블록 생성 간격의 여유 범위 +- 2분
)

type blockchain struct {
	NewstHash         string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

var b *blockchain // 이건 blockchain package 안에서만 사용할 수 있음!
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewstHash, b.Height+1, getdifficulty(b))
	b.NewstHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
}

func persistBlockchain(b *blockchain) {
	db.SaveCheckpoint(utils.ToBytes(b))
	// 블록체인 구조체를 바이트로 변환한 슬라이스 []byte를 인수로 넘겨준다
}

func BlocksSlice(b *blockchain) []*Block {
	var Blocks []*Block
	hashCursor := b.NewstHash
	for {
		block, _ := FindBlock(hashCursor) // FindBlock이 *Block, err를 리턴하기 때문에
		Blocks = append(Blocks, block)    // 가장 최근의 블록부터 넣기 때문에 Blocks의 첫 블록은 가장 최근의 블록임
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return Blocks
}

func recalculateDifficulty(b *blockchain) int {
	allBlocks := BlocksSlice(b)
	newestBlock := allBlocks[0]
	lastRecalculateBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp - lastRecalculateBlock.Timestamp) / 60
	// 유닉스 시간에 60을 나눠줬으니까 분 단위 시간으로 바뀜
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	} else { // 우리가 예측해논 5*2분 과 actualTime이 같으면
		return b.CurrentDifficulty
	}
}

func getdifficulty(b *blockchain) int {
	if b.Height == 0 {
		// b.CurrentDifficulty = defaultDifficulty // 첫 난이도를 현재 난이도에 넣어준다 // 이건 그냥 내가 임의로 넣어준 코드
		//근데 생각하기론 필요하다고 생각이 들긴 하는데...
		// 나중에 다시 보고 고쳐보는 걸로
		return defaultDifficulty // 첫 블록 생성이라면 난이도를 defaultDifficulty = 2로 설정
	} else if b.Height%difficultyInterval == 0 {
		return recalculateDifficulty(b)
	} else { // 블록 생성 간격이 5 이하라면 // 저장된 난이도를 변동없이 그대로 리턴한다
		return b.CurrentDifficulty
	}
}

func UTxOutsByAddress(address string, b *blockchain) []*UTxOut { // 함수를 public으로 export 되게 해놓은 이유는 이 함수를 API에서 불러올 것이기 때문이다
	var uTxOuts []*UTxOut
	sTxOuts := make(map[string]bool)
	for _, block := range BlocksSlice(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					sTxOuts[input.TxID] = true
				}
			}

			for index, output := range tx.TxOuts {
				// 이걸 빠뜨림 소유자랑 넣은 주소랑 같으면 // 그래서 jiwon에도 똑같이 output이 생성되어 있었나보다
				if output.Owner == address {
					if _, ok := sTxOuts[tx.Id]; !ok {
						uTxOut := &UTxOut{tx.Id, index, output.Amount}
						if !IsOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func TotalBalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		checkpoint := db.GetCheckPointFromDb()
		if checkpoint == nil {
			b.AddBlock()
		} else {
			b.restore(checkpoint)
		}
	})
	return b
}
