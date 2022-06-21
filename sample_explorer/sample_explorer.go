package sample_explorer

import (
	"elize/elizebch"
	"elize/elizeutils"
	"elize/wallet"
	"fmt"
	"html/template"
	"net/http"
)

const (
	pagesLocation    string = "templates/pages/*.html"
	partialsLocation string = "templates/partials/*.html"
)

type homeData struct {
	PageTitle   string
	Blocks      []*elizebch.Block
	BlockHeight Info
}

type addData struct {
	PageTitle string
	Balance   Info
}

type memData struct {
	PageTitle string
	MemHeight Info
	Tx        []elizebch.Tx
}

type PostedData struct {
	User    string
	Balance float64
}

type Info struct {
	Name  string
	Value string
}

var (
	templates *template.Template
	port      string = ":5000"
)

func home(rw http.ResponseWriter, r *http.Request) {
	currentChain := elizebch.AllBlock()
	data := homeData{"HOME", currentChain, Info{"Block Height", fmt.Sprint(currentChain[0].Height)}}
	elizeutils.Errchk(templates.ExecuteTemplate(rw, "home", data))
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data := addData{"ADD", Info{"Your Balance", fmt.Sprint(elizebch.BalanceByAddress(wallet.Wallet().Address))}}
		elizeutils.Errchk(templates.ExecuteTemplate(rw, "add", data))
	case "POST":
		r.ParseForm()
		Pdata := PostedData{r.Form.Get("user"), elizeutils.ToInt(r.Form.Get("balance"))}
		_, err := elizebch.ElizeMempool().AddTxs(Pdata.User, Pdata.Balance)
		if err != nil {
			http.Redirect(rw, r, "/errors", http.StatusPermanentRedirect)
		} else {
			http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
		}
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data := memData{PageTitle: "MEMPOOL", MemHeight: Info{"Mempool Height", fmt.Sprint(len(elizebch.ElizeMempool().Txs))}}
		data.Tx = elizebch.ElizeMempool().AllMemTx()
		elizeutils.Errchk(templates.ExecuteTemplate(rw, "mempool", data))
	case "POST":
		elizebch.GetBlockchain().AddBlock()
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func errors(rw http.ResponseWriter, r *http.Request) {
	data := addData{"ERROR!", Info{"\U0001f622 Check your balance : ", fmt.Sprint(elizebch.BalanceByAddress(wallet.Wallet().Address))}}
	elizeutils.Errchk(templates.ExecuteTemplate(rw, "errors", data))
}

func Start(explorerPort int) {
	port = fmt.Sprintf(":%d", explorerPort)
	fmt.Printf("Listening on http://localhost%s\n", port)
	templates = template.Must(template.ParseGlob(pagesLocation))
	templates = template.Must(templates.ParseGlob(partialsLocation))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	http.HandleFunc("/mempool", mempool)
	http.HandleFunc("/errors", errors)
	http.ListenAndServe(port, nil)
}
