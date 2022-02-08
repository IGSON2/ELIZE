package restapi

import (
	"elizebch/elizebch"
	"elizebch/elizeutils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var port string

type AddTxPayload struct {
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type BalanceResponse struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
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
	return []byte(fmt.Sprintf("http://localhost%s%s", port, *u)), nil
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         URLConverter("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         URLConverter("/blocks"),
			Method:      "GET",
			Description: "See All blocks",
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
			URL:         URLConverter("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
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
		json.NewEncoder(rw).Encode(errorResponse{elizebch.NotEnoughBalanceErr})
	} else {
		rw.WriteHeader(http.StatusCreated)
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(elizebch.ElizeMempool)
}

func Start(apiPort int) {
	port = fmt.Sprintf(":%d", apiPort)
	gorillaMux := mux.NewRouter()
	gorillaMux.Use(middleWare)
	gorillaMux.HandleFunc("/", documentation).Methods("GET")
	gorillaMux.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	gorillaMux.HandleFunc("/blocks/{hash:[a-f0-9]+}", oneblock).Methods("GET")
	gorillaMux.HandleFunc("/balance/{address}", balance).Methods("GET")
	gorillaMux.HandleFunc("/mempool", mempool).Methods("GET")
	gorillaMux.HandleFunc("/transaction", transaction).Methods("GET", "POST")
	fmt.Printf("Listening on http://localhost%s\n", port)
	elizeutils.Errchk(http.ListenAndServe(port, gorillaMux))
}
