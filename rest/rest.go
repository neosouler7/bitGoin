package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neosouler7/bitGoin/blockchain"
	"github.com/neosouler7/bitGoin/utils"
	"github.com/neosouler7/bitGoin/wallet"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) { // use Go's auto implemented interface
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"` // field struct tag
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type addTxPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
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
			Description: "See the status of the Blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET, POST",
			Description: "[GET] Get all blocks, [POST] Add a block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See a block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for and address",
		},
		{
			URL:         url("/mempool"),
			Method:      "GET",
			Description: "Check current memory pool",
		},
		{
			URL:         url("/wallet"),
			Method:      "GET",
			Description: "Check my wallet address",
		},
	}
	err := json.NewEncoder(rw).Encode(data) // replace Marshalling & return to writer
	utils.HandleErr(err)
}

func status(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blockchain()))
	}
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain())))
	case "POST":
		blockchain.Blockchain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		hash := vars["hash"]
		block, err := blockchain.FindBlock(hash)
		encoder := json.NewEncoder(rw)
		if err == blockchain.ErrNotFound {
			utils.HandleErr(encoder.Encode(errorResponse{fmt.Sprintf("%s", err)}))
		} else {
			utils.HandleErr(encoder.Encode(block))
		}
	}
}

func balance(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		address := vars["address"]
		total := r.URL.Query().Get("total")
		switch total {
		case "true":
			amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
			utils.HandleErr(json.NewEncoder(rw).Encode(balanceResponse{address, amount}))
		default:
			utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain())))
		}
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
	}
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		address := wallet.Wallet().Address
		utils.HandleErr(json.NewEncoder(rw).Encode(struct {
			Address string `json:"address"`
		}{Address: address}))
	}
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addTxPayload
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
		err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			utils.HandleErr(json.NewEncoder(rw).Encode(errorResponse{err.Error()}))
			return
		}
		rw.WriteHeader(http.StatusCreated)
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application.json")
		next.ServeHTTP(rw, r)
	})
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleware) // use middleware
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
