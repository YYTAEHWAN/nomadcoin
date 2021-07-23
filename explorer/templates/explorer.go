package explorer

import (
	"fmt"
	"learngo/github.com/nomadcoders/blockchain"
	"log"
	"net/http"
	"text/template"
)

var templates *template.Template

const (
	templateDir string = "explorer/templates/"
)

// template 는 이미 import 되고 있기 때문이다
type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", nil}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		blockchain.Blockchain().AddBlock()
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}
func Start(port int) {
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	// 위엣 줄에 의해서 templates는 Object가 되었고
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	// 위엣줄에 의해서 template Object가 partial까지 load하게 되었다.
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
