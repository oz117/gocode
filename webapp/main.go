package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
} // Content of a page

func (p *Page) savePage() error {
	fileName := p.Title + ".txt"
	return ioutil.WriteFile(fileName, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]

	p, _ := loadPage(title)
	t, _ := template.ParseFiles("Views/view.html")
	t.Execute(w, p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, _ := template.ParseFiles("Views/edit.html")
	t.Execute(w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("content")

	p := &Page{Title: title, Body: []byte(body)}
	p.savePage()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":4000", nil)
}
