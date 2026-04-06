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
	r.HandleFunc("/", handlers.Home)          // only GET
    r.HandleFunc("/contact", handlers.Contact) // only POST
	r.HandleFunc("/index2", handlers.Home2)
	r.HandleFunc("/blog", auth.RequireLogin(handlers.Blog))
	r.HandleFunc("/blog/{slug}", auth.RequireLogin(handlers.BlogPost))
	r.HandleFunc("/blog/{slug}/comment", auth.RequireLogin(handlers.CreateCommentHandler))
	r.HandleFunc("/login", handlers.LoginHandler)
	r.HandleFunc("/logout", handlers.LogoutHandler)
	r.HandleFunc("/register", auth.RequireLogin(auth.RequireAdmin(handlers.RegisterHandler)))
	r.HandleFunc("/admin/sliders/create", auth.RequireLogin(handlers.SliderAdmin))
	r.HandleFunc("/admin/sliders", auth.RequireLogin(handlers.AdminSliders))
	r.HandleFunc("/admin/posts", auth.RequireLogin(auth.RequireAdmin(handlers.AdminPosts)))
	r.HandleFunc("/admin/contact/read", auth.RequireLogin(auth.RequireAdmin(handlers.MarkContactRead)))
	r.HandleFunc("/admin/contact/delete", auth.RequireLogin(auth.RequireAdmin(handlers.DeleteContact)))
	r.HandleFunc("/admin/contact", auth.RequireLogin(auth.RequireAdmin(handlers.ManageQueries)))
	r.HandleFunc("/admin/posts/create", auth.RequireLogin(auth.RequireAdmin(handlers.AdminCreatePost)))
	r.HandleFunc("/comment/delete", auth.RequireLogin(handlers.DeleteCommentHandler))

	r.HandleFunc("/admin/posts/edit/{id}", auth.RequireLogin(auth.RequireAdmin(handlers.AdminEditPost)))

	r.HandleFunc("/admin/posts/delete/{id}",
	 auth.RequireLogin(auth.RequireAdmin(handlers.AdminDeletePost)),
	).Methods("POST")

	r.PathPrefix("/style/").
    Handler(http.StripPrefix("/style/",
        http.FileServer(http.Dir("./static/slowave/style")),
    ))

	r.PathPrefix("/uploads/").
	Handler(http.StripPrefix("/uploads/",
		http.FileServer(http.Dir("./static/uploads")),
	))

	http.ListenAndServe(":8080", r)
}