package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"word_press/auth"
	"word_press/models"
)

func AdminDashboard(w http.ResponseWriter, r *http.Request) {

	user, _ := auth.GetCurrentUser(r)

	data := map[string]interface{} {
		"User": user,
		"Now": time.Now()
	}

	tmpl := template.Must(template.ParseFiles)

}