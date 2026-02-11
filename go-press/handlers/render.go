package handlers

import (
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, page string, data interface{}){

	tmpl := template.Must(template.ParseFiles(
		"templates/slowave/layout.html"
		"templates/slowave/header.html"
		"template/slowave/" + page
	))

	tmpl.ExecuteTemplate(w, "layout.html", data)
} 