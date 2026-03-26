package handlers

import (
	"html/template"
	"net/http"
)

type RenderOptions struct {
	Page string
	Header string
	Data interface{}
}

func Render(w http.ResponseWriter, page string, data interface{}){

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	tmpl, err := template.New("layout.html").
		Funcs(funcMap).
		ParseFiles(
			"templates/slowave/layout.html",
			"templates/slowave/home-header.html",
			"templates/slowave/footer.html",
			"templates/slowave/"+page,
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

func RenderWithOpts(w http.ResponseWriter, opts RenderOptions){

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	tmpl, err := template.New("layout.html").
		Funcs(funcMap).
		ParseFiles(
			"templates/slowave/layout.html",
			"templates/slowave/" + opts.Header,
			"templates/slowave/footer.html",
			"templates/slowave/" + opts.Page,
		)

	if err != nil {
	  http.Error(w, err.Error(), 500)
	  return
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", opts.Data)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}