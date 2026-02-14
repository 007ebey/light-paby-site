package handlers

import (
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	Render(w, "index.html", map[string]interface{}{
		"SiteTitle": "Cool Site",
	})
}

func Home2(w http.ResponseWriter, r *http.Request) {
	Render(w, "index2.html", map[string]interface{}{
		"SiteTitle": "Cool Site: Layout 2",
	})
}