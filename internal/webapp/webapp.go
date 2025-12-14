package webapp

import "os"

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) Save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func LoadPage(title string) *Page {
	filename := title + ".txt"
	body, _ := os.ReadFile(filename)
	return &Page{Title: title, Body: body}
}

func LoadPackageReport(title string) *Page {
	filename := "./" + title + "_report.json"
	body, _ := os.ReadFile(filename)
	return &Page{Title: title, Body: body}
}
