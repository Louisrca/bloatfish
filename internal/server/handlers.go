package server

import (
	"encoding/json"
	"net/http"
	"os"
)

func loadPackageReport(reportName string) (*Page, error) {
	filename := "./" + reportName + "_report.json"

	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var report DependencyReport
	if err := json.Unmarshal(body, &report); err != nil {
		return nil, err
	}

	return &Page{
		Title: reportName,
		Body:  &report,
	}, nil
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadPackageReport("full_audit")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ViewPage(w, p)
}
