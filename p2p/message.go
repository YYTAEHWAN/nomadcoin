package p2p

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nomadcoders/blockchain"
	"github.com/nomadcoders/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksResponse, blockchain.BlocksSlice(blockchain.Blockchain()))
	p.inbox <- m
}

func makeMessage(kind MessageKind, p interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(p),
	}
	return utils.ToJSON(m)
}

func sendNewestBlock(p *peer) {
	fmt.Printf("Sending NewestBlock to %s\n", p.key)
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewstHash)
	utils.HandleErr(err)
	//완성된 메세지를 채널을 통해 보내줄 것임
	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func handleMsg(m *Message, p *peer) { // 이 함수는 포트 3000에 의해서 실행되고 ,,,, 매개변수인 p *peer은 4000의 peer임
	switch m.Kind {
	// 이 switch문은 3000이 실행하는 것임
	case MessageNewestBlock:
		fmt.Printf("Received the newest block from %s\n", p.key)
		// 이 최신블록 메세지는 4000으로부터 온 것임 (4000의 최신블록이 담겨있음)
		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Println(payload)
		b, err := blockchain.FindBlock(blockchain.Blockchain().NewstHash)
		utils.HandleErr(err)
		if payload.Height >= b.Height {
			fmt.Printf("Requesting All Blocks from %s\n", p.key)
			requestAllBlocks(p)
		} else {
			fmt.Printf("Sending newest Block to  %s\n", p.key)
			sendNewestBlock(p)
		}
	case MessageAllBlocksRequest:
		fmt.Printf("%s wants All The Blocks\n", p.key)
		sendAllBlocks(p)
	case MessageAllBlocksResponse:
		fmt.Printf("Received All the Blcoks from %s\n", p.key)
		var payload []*blockchain.Block
		json.Unmarshal(m.Payload, &payload)
		blockchain.Blockchain().Replace(payload)
	case MessageNewBlockNotify:
		fmt.Printf("Received New Block from %s\n", p.key)
		var payload *blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().AddPeerBlock(payload)
	case MessageNewTxNotify:
		var payload *blockchain.Tx
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddPeerTxOnMem(payload)
	case MessageNewPeerNotify:
		var payload string
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		parts := strings.Split(payload, ":")
		AddPeers(parts[0], parts[1], parts[2], false)
	}
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageNewTxNotify, tx)
	p.inbox <- m
}

func notifyNewPeer(payload string, p *peer) {
	m := makeMessage(MessageNewPeerNotify, payload)
	p.inbox <- m
}
