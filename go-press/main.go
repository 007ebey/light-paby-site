package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"word_press/database"
	"word_press/handlers"
)

func main() {
	database.InitDB()

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Home)
	r.HandleFunc("/post/{slug}", handlers.SinglePost)

	r.HandleFunc("/admin", handlers.AdminDashboard)
	r.HandleFunc("/admin/posts", handlers.AdminPosts)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	http.ListenAndServe(":8080", r)
}