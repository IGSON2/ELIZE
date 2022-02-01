package restapi

import (
	"elizebch/elizebch"
	"elizebch/elizeutils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	port string = ":4000"
)

type AddBlockBody struct {
	Message string
}

type URLConverter string

type URLDescription struct {
	URL         URLConverter `json:"url"`
	Method      string       `json:"method"`
	Description string       `json:"description"`
	Payload     string       `json:"payload,omitempty"`
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
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data := elizebch.Allblock()
		json.NewEncoder(rw).Encode(data)
	case "POST":
		var addData AddBlockBody
		elizeutils.Errchk(json.NewDecoder(r.Body).Decode(&addData))
		elizebch.GetBlockchain().AddBlock(addData.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func oneblock(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashId := vars["hash"]
	resultBlock, err := elizebch.FindOneblock(hashId)
	if err == nil {
		json.NewEncoder(rw).Encode(*resultBlock)
	} else {
		json.NewEncoder(rw).Encode(err)
	}
}

func Start() {
	gorillaMux := mux.NewRouter()
	gorillaMux.HandleFunc("/", documentation).Methods("GET")
	gorillaMux.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	gorillaMux.HandleFunc("/blocks/{hash:[a-f0-9]+}", oneblock).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)
	elizeutils.Errchk(http.ListenAndServe(port, gorillaMux))
}
