package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nomadcoders/blockchain"
	"github.com/nomadcoders/p2p"
	"github.com/nomadcoders/utils"
	"github.com/nomadcoders/wallet"

	"github.com/gorilla/mux"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type AddBlockBody struct {
	Message string
}

type totalBalanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type addTxpayload struct {
	To     string
	Amount int
}

type MyWalletResponse struct {
	Address string `json:"address"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type AddpeerPayload struct {
	Address string
	Port    string
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.BlocksSlice(blockchain.Blockchain()))
	case "POST":
		newBlock := blockchain.Blockchain().AddBlock()
		p2p.BroadcastNewBlock(newBlock)
		rw.WriteHeader(http.StatusCreated)
	}
}
func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the Blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{address}"),
			Method:      "GET",
			Description: "Amount balance for a address",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to WebSockets",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}

}

func status(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.Blockchain(), rw)
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.TotalBalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(rw).Encode(totalBalanceResponse{address, amount})
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain())))
	}

}

func jsonContentTypeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	//blockchain.Txstatus(blockchain.Mempool(), rw)
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxpayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	// 만들어진 Tx를 다른 노드들에게도 전달
	// 이 코드는 REST API?로 만든 Tx를 다른 노드들에게 보내는 Tx만든 노드 시점이야
	p2p.BroadcastNewTx(tx)
	// NewEncoder이 실행되고 난 뒤에 WriteHeader이 또다시 실행되면 ㅠㅠㅠ 오류가 뜹니다
	rw.WriteHeader(http.StatusCreated)
}

func mywallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(MyWalletResponse{Address: address})
}

func getPeerServer(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var APayload *AddpeerPayload
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&APayload))
		fmt.Printf("REST API로 받은 데이터는 %s , %s\n", APayload.Address, APayload.Port)
		p2p.AddPeers(APayload.Address, APayload.Port, port[1:], true) // 여기서 port는 지금 열려있는 Openport 번호, Apayload는 보내온 Address and Port
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	}
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleWare, loggerMiddleWare)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", mywallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", getPeerServer).Methods("GET", "POST")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
