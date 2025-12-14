package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Louisrca/bloatfish/internal/webapp"
)

func ViewHandler(w http.ResponseWriter, r *http.Request) {

	p := webapp.LoadPackageReport("unused_packages")
	fmt.Printf("hello", p.Body)
	fmt.Fprintf(w, "<h1>Editing %s</h1>"+
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>",
		p.Title, p.Title, string(p.Body))
}

func StartServer() {
	http.HandleFunc("/view/", ViewHandler)
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
