package blockchain

import (
	"errors"
	"learngo/github.com/nomadcoders/utils"
	"learngo/github.com/nomadcoders/wallet"
	"time"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct { // Tx = transaction
	ID        string   `json:"id"` // tx의 해시값?
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct { // TxIn = transaction Input
	TxID      string `json:"txID"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct { // TxOut = transaction Output2
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type UTxOut struct { //TxID의 TxOuts 중 Index번째 TxOut을 뜻함, 게다가 사용되지 않은 TxOut임
	TxID   string
	Index  int
	Amount int
}

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.ID, *wallet.Wallet())
		// 만드려는 블록에 담긴 모든 txIn에 싸인을 한다
	}
}

func vaildate(tx *Tx) bool {
	vaild := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			vaild = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		vaild = wallet.Verify(txIn.Signature, tx.ID, address)
		// 여기서 Verify()에 넣는 ID는 txIn.TxID == prevTx.ID  이 2가지는 모두 같은 ID임 ////// tx.ID 이거는 다를껄?
		if !vaild { //vaild == false 보다 더 나은 방식
			break
		}
	}
	return vaild
}

func makeCoinbaseTx(address string) *Tx {
	// 채굴자를 주소로 삼는 코인베이스 거래내역을 생성해서 Tx 포인터를 리턴할 것이다
	txIns := []*TxIn{
		{"", -1, "COINBASE"}, // coinbase transaction이기 때문에
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

var ErrorNoMoney = errors.New("not enought 돈")
var ErrorNotVaild = errors.New("transaction not vaild")

func makeTx(address string, to string, amount int) (*Tx, error) {
	if TotalBalanceByAddress(address, Blockchain()) < amount {
		return nil, ErrorNoMoney
	}

	var TxIns []*TxIn
	var TxOuts []*TxOut
	uTxOuts := UTxOutsByAddress(address, Blockchain())
	total := 0

	for _, uTxOut := range uTxOuts { // 사용될 TxOut들을 TxIn으로 옮기는 과정  //// 근데 이거  복사야 아니면 참조에 의한 변경이야?
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, address} // 이 address은 하드코딩이므로 나중에 바꿔버릴꺼야~
		TxIns = append(TxIns, txIn)                       // 의문 첫번째 - 위엣줄 코딩을 하게 되면 저 기록들이 옮겨져? 아니면 복사만 되어서 들어가는거야?
		total += uTxOut.Amount
	}
	if change := total - amount; change > 0 { // 거스름돈 받는 코드
		changeTxOut := &TxOut{address, change}
		TxOuts = append(TxOuts, changeTxOut)
	}
	TotxOut := &TxOut{to, amount}
	TxOuts = append(TxOuts, TotxOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     TxIns,
		TxOuts:    TxOuts,
	}
	tx.getId()
	tx.sign()
	vaild := vaildate(tx)
	if !vaild {
		return nil, ErrorNotVaild
	}
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	//transaction을 mempool에 추가해줄 뿐, 거래를 만들진 않는다
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToComfirm() []*Tx {
	coinBase := makeCoinbaseTx(wallet.Wallet().Address)
	txs := m.Txs
	txs = append(txs, coinBase)
	m.Txs = nil // 밈풀 초기화
	return txs
}

func IsOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			// mempool에 들어있는 inputs들은 모두 사용될 예정인 TxOut들이지 그러니 input과 같은 ID, Index가 있다면 삐삑 넌 걸렸어 UTxOuts로 다시 못들어가
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
}
