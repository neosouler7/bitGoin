package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/neosouler7/bitGoin/blockchain"
	"github.com/neosouler7/bitGoin/utils"
)

const (
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	utils.HandleErr(templates.ExecuteTemplate(rw, "home", data))
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.HandleErr(templates.ExecuteTemplate(rw, "add", nil))
	case "POST":
		utils.HandleErr(r.ParseForm())
		data := r.Form.Get("blockData")
		blockchain.GetBlockchain().AddBlock((data))
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(port int) {
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	fmt.Printf("listening to http://localhost:%d\n", port)
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}