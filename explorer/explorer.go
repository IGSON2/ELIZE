package explorer

import (
	"elizebch/elizebch"
	"elizebch/elizeutils"
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
	BlockHeight int
}

type memData struct {
	PageTitle string
	TxIns     []elizebch.TxIn
	TxOuts    []elizebch.TxOut
}

type PostedData struct {
	User    string
	Balance float64
}

var (
	templates *template.Template
	port      string = ":5000"
)

func home(rw http.ResponseWriter, r *http.Request) {
	currentChain := elizebch.AllBlock()
	data := homeData{"HOME", currentChain, currentChain[0].Height}
	elizeutils.Errchk(templates.ExecuteTemplate(rw, "home", data))
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data := homeData{"ADD", nil, elizebch.GetBlockchain().Height}
		elizeutils.Errchk(templates.ExecuteTemplate(rw, "add", data))
	case "POST":
		r.ParseForm()
		Pdata := PostedData{r.Form.Get("user"), elizeutils.ToInt(r.Form.Get("balance"))}
		elizebch.ElizeMempool().AddTxs(Pdata.User, Pdata.Balance)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data := memData{PageTitle: "MEMPOOL"}
		data.TxIns, data.TxOuts = elizebch.ElizeMempool().AllMemTx()
		elizeutils.Errchk(templates.ExecuteTemplate(rw, "mempool", data))
	case "POST":
		elizebch.GetBlockchain().AddBlock()
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(explorerPort int) {
	port = fmt.Sprintf(":%d", explorerPort)
	fmt.Printf("Listening on http://localhost%s\n", port)
	templates = template.Must(template.ParseGlob(pagesLocation))
	templates = template.Must(templates.ParseGlob(partialsLocation))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	http.HandleFunc("/mempool", mempool)
	http.ListenAndServe(port, nil)
}
