package restapi

import (
	"elizebch/elizebch"
	"elizebch/elizeutils"
	"elizebch/p2p"
	"elizebch/wallet"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var restPort string

type addPeerPayload struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type AddTxPayload struct {
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type BalanceResponse struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

type WalletAddress struct {
	Address string `json:"wallet_address"`
}

type URLConverter string

type URLDescription struct {
	URL         URLConverter `json:"url"`
	Method      string       `json:"method"`
	Description string       `json:"description"`
	Payload     string       `json:"payload,omitempty"`
}

type errorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func (u *URLConverter) MarshalText() (test []byte, err error) {
	return []byte(fmt.Sprintf("http://localhost%s%s", restPort, *u)), nil
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         URLConverter("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         URLConverter("/status"),
			Method:      "GET",
			Description: "See the Status of the Blockchain",
		},
		{
			URL:         URLConverter("/blocks"),
			Method:      "GET",
			Description: "See All Blocks",
		},
		{
			URL:         URLConverter("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         URLConverter("/blocks/{hash}"),
			Method:      "GET",
			Description: "See A Block",
		},
		{
			URL:         URLConverter("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
		{
			URL:         URLConverter("/ws"),
			Method:      "GET",
			Description: "Upgrade to WebSockets",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func jsonMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RequestURI)
		next.ServeHTTP(rw, r)
	})
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data := elizebch.AllBlock()
		json.NewEncoder(rw).Encode(data)
	case "POST":
		elizebch.GetBlockchain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	}
}

func oneblock(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashId := vars["hash"]
	resultBlock, err := elizebch.FindBlock(hashId)
	if err == nil {
		json.NewEncoder(rw).Encode(*resultBlock)
	} else {
		json.NewEncoder(rw).Encode(errorResponse{fmt.Sprintf("%s", err)})
	}
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	isTotalQuery := r.URL.Query().Get("total")
	switch isTotalQuery {
	case "true":
		balance := elizebch.BalanceByAddress(address)
		elizeutils.Errchk(json.NewEncoder(rw).Encode(BalanceResponse{address, balance}))
	default:
		elizeutils.Errchk(json.NewEncoder(rw).Encode(elizebch.UTxOutsByAddress(address)))
	}
}

func transaction(rw http.ResponseWriter, r *http.Request) {
	var txReqPayload AddTxPayload
	elizeutils.Errchk(json.NewDecoder(r.Body).Decode(&txReqPayload))
	err := elizebch.ElizeMempool.AddTxs(txReqPayload.To, txReqPayload.Amount)
	if err != nil {
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
	} else {
		rw.WriteHeader(http.StatusCreated)
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(elizebch.ElizeMempool)
}

func userWallet(rw http.ResponseWriter, r *http.Request) {
	address := WalletAddress{wallet.Wallet().Address}
	json.NewEncoder(rw).Encode(address)
}

func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		elizeutils.Errchk(json.NewEncoder(rw).Encode(p2p.Peers))
	case "POST":
		var peerPayload addPeerPayload
		json.NewDecoder(r.Body).Decode(&peerPayload)
		p2p.Addpeer(peerPayload.Ip, peerPayload.Port, restPort)
	}
}

func Start(apiPort int) {
	restPort = fmt.Sprintf(":%d", apiPort)
	gorillaMux := mux.NewRouter()
	gorillaMux.Use(jsonMiddleWare, loggerMiddleWare)
	gorillaMux.HandleFunc("/", documentation).Methods("GET")
	gorillaMux.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	gorillaMux.HandleFunc("/blocks/{hash:[a-f0-9]+}", oneblock).Methods("GET")
	gorillaMux.HandleFunc("/balance/{address}", balance).Methods("GET")
	gorillaMux.HandleFunc("/mempool", mempool).Methods("GET")
	gorillaMux.HandleFunc("/transaction", transaction).Methods("GET", "POST")
	gorillaMux.HandleFunc("/wallet", userWallet).Methods("GET", "POST")
	gorillaMux.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	gorillaMux.HandleFunc("/peers", peers).Methods("GET", "POST")
	fmt.Printf("Listening on http://localhost%s\n", restPort)
	elizeutils.Errchk(http.ListenAndServe(restPort, gorillaMux))
}
