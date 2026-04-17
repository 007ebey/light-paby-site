package handlers

import (
	"html/template"
	"net/http"
)

type RenderOptions struct {
	Page string
	Header string
	Footer string
	Data interface{}
}

func RenderWithOpts(w http.ResponseWriter, opts RenderOptions){

	if (opts.Header == "default") {
		opts.Header = "header.html"
	}

	if (opts.Footer == "default") {
		opts.Footer = "footer.html"
	}

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"seq": func(start, end int) []int {
			var s []int
			for i := start; i <= end; i++ {
				s = append(s, i)
			}
			return s
		},
	}

	tmpl, err := template.New("layout.html").
		Funcs(funcMap).
		ParseFiles(
			"templates/slowave/layout.html",
			"templates/slowave/" + opts.Header,
			"templates/slowave/" + opts.Footer,
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