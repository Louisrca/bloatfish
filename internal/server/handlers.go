package server

import (
	"net/http"
	"os"
)

func loadPackageReport(title string) *Page {
	filename := "./" + title + "_report.json"
	body, _ := os.ReadFile(filename)
	return &Page{Title: title, Body: body}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPackageReport("unused_packages")
	ViewPage(w, p)
}
