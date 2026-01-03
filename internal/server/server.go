package server

import (
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
	log.Println("Starting web UI on http://localhost:8080")

	fs := http.FileServer(http.Dir("./app/bloatfish-app/dist"))
	http.Handle("/", fs)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
