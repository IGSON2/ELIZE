package explorer

import (
	"elizebch/elizebch"
	"elizebch/elizeutils"
	"fmt"
	"html/template"
	"net/http"
)

const (
	port             string = ":5000"
	pagesLocation    string = "templates/pages/*.gohtml"
	partialsLocation string = "templates/partials/*.gohtml"
)

type excuteData struct {
	PageTitle   string
	Blocks      []*elizebch.Block
	BlockHeight int
}

var templates *template.Template

func home(rw http.ResponseWriter, r *http.Request) {
	currentChain := elizebch.GetBlockchain().Blocks
	data := excuteData{"HOME", currentChain, currentChain[0].Height}
	elizeutils.Errchk(templates.ExecuteTemplate(rw, "home", data))
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data := excuteData{"ADD", nil, 0}
		elizeutils.Errchk(templates.ExecuteTemplate(rw, "add", data))
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		elizebch.GetBlockchain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start() {
	fmt.Printf("Listening on http://localhost%s\n", port)
	templates = template.Must(template.ParseGlob(pagesLocation))
	templates = template.Must(templates.ParseGlob(partialsLocation))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	http.ListenAndServe(port, nil)
}
