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
func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
	// 블록체인 구조체를 바이트로 변환한 슬라이스 []byte를 인수로 넘겨준다
}
func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewstHash, b.Height+1)
	b.NewstHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockchain) BlocksSlice() []*Block {
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

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.BlocksSlice()
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

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		b.CurrentDifficulty = defaultDifficulty // 첫 난이도를 현재 난이도에 넣어준다
		return defaultDifficulty                // 첫 블록 생성이라면 난이도를 defaultDifficulty = 2로 설정
	} else if b.Height%difficultyInterval == 0 {
		return b.recalculateDifficulty()
	} else { // 블록 생성 간격이 5 이하라면 // 저장된 난이도를 변동없이 그대로 리턴한다
		return b.CurrentDifficulty
	}
}
func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			checkpoint := db.GetCheckPointFromDb()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				b.restore(checkpoint)
			}
		})
	}
	return b
}
