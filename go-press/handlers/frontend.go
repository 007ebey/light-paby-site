package handlers

import (
	"html/template"
	"net/http"
	"github.com/gorilla/mux"
	"word_press/models"
)

func Home(w http.ResponseWriter, r *http.Request) {
	Render(w, "index.html", map[string]interface{}{
		"SiteTitle": "Slowave CMS"
	})
}