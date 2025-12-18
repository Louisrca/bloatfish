package server

import (
	"fmt"
	"log"
	"net/http"
)

type PackageInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	SizeStr string `json:"size_str"`
}
type DependencyReport struct {
	Declared     []PackageInfo `json:"declared"`
	Used         []PackageInfo `json:"used"`
	Unused       []PackageInfo `json:"unused"`
	Indirect     []PackageInfo `json:"indirect"`
	UnusedDev    []PackageInfo `json:"unused_dev"`
	Framework    string        `json:"framework"`
	TotalSize    int64         `json:"total_size"`
	TotalSizeStr string        `json:"total_size_str"`
	Errors       []string      `json:"errors"`
}

type Page struct {
	Title string
	Body  *DependencyReport
}

func StartServer() {
	http.HandleFunc("/", ViewHandler)
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./internal/app/static")),
		),
	)
	fmt.Println("Starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
