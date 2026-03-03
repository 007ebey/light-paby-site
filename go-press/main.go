package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"word_press/database"
	"word_press/handlers"
	"word_press/auth"
)

func main() {
	database.InitDB()

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Home)
	r.HandleFunc("/index2", handlers.Home2)
	r.HandleFunc("/blog", handlers.Blog)
	r.HandleFunc("/blog/{slug}", handlers.BlogPost)
	r.HandleFunc("/login", handlers.LoginHandler)
	r.HandleFunc("/register", handlers.RegisterHandler)
	r.HandleFunc("/admin/sliders/create", auth.RequireLogin(handlers.SliderAdmin))
	r.HandleFunc("/admin/sliders", auth.RequireLogin(handlers.AdminSliders))
	r.HandleFunc("/admin/posts", auth.RequireLogin(auth.RequireAdmin(handlers.AdminPosts)))
	r.HandleFunc("/admin/posts/create", auth.RequireLogin(auth.RequireAdmin(handlers.AdminCreatePost)))

	r.HandleFunc("/admin/posts/edit/{id}", auth.RequireLogin(auth.RequireAdmin(handlers.AdminEditPost)))

	r.PathPrefix("/style/").Handler(http.StripPrefix("/style/",
	 http.FileServer(http.Dir("./static/slowave/style")),
	))
	
	http.ListenAndServe(":8080", r)
}