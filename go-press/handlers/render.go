package handlers

import (
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, page string, data interface{}){

	tmpl, err := template.ParseFiles(
		"templates/slowave/layout.html",
		"templates/slowave/header.html",
		"templates/slowave/footer.html",
		"templates/slowave/" + page,
	)

	if err != nil {
	  http.Error(w, err.Error(), 500)
	  return
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", data)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
} 