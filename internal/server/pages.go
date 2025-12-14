package server

import (
	"net/http"
	"text/template"
)

const templateDir = "./internal/html/"

func ViewPage(w http.ResponseWriter, p *Page) {
	t, _ := template.ParseFiles(templateDir + "packages.html")
	t.Execute(w, p)
}
