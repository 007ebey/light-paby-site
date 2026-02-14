package handlers

import (
	"net/http"
	"time"
	"word_press/auth"
)

func AdminDashboard(w http.ResponseWriter, r *http.Request) {

	user, _ := auth.GetCurrentUser(r)

	data := map[string]interface{} {
		"SiteTitle": "Admin Dashboard",
		"User": user,
		"Now": time.Now(),
	}

	Render(w, "admin/dashboard.html", data)
}