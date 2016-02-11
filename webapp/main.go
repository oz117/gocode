package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("Views/view.html", "Views/edit.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type page struct {
	Title string
	Body  []byte
} // Content of a page

func (p *page) savepage() error {
	fileName := p.Title + ".txt"
	return ioutil.WriteFile(fileName, p.Body, 0600)
}

func loadpage(title string) (*page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	return &page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)

	if err != nil {
		return
	}
	p, err := loadpage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view.html", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)

	if err != nil {
		return
	}
	p, err := loadpage(title)
	if err != nil {
		p = &page{Title: title}
	}
	renderTemplate(w, "edit.html", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	body := r.FormValue("content")

	p := &page{Title: title, Body: []byte(body)}
	if err := p.savepage(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *page) {
	if err := templates.ExecuteTemplate(w, tmpl, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":4000", nil)
}
