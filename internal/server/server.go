package server

import (
	"fmt"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func StartServer() {
	http.HandleFunc("/", ViewHandler)
	fmt.Println("Starting server on http://localhost:3031")
	log.Fatal(http.ListenAndServe(":3031", nil))
}
