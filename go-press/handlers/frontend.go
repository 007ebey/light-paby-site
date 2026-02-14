package handlers

import (
	"net/http"
	"word_press/models"
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

func SliderAdmin(w http.ResponseWriter, r *http.Request) {
	Render(w, "admin/create-slider.html", map[string]interface{}{
		"SiteTitle": "Admin: Slide",
	})
}

func AdminSliders(w http.ResponseWriter, r *http.Request) {
	
}