package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"word_press/auth"
	"word_press/database"
	"word_press/handlers"
	"word_press/services"
)

func main() {
	database.InitDB()

	r := mux.NewRouter()

	// =========================
	// SERVICES + HANDLERS
	// =========================

	// Auth
	userRepo := services.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Posts
	postRepo := services.NewPostRepository()
	postService := services.NewPostService(postRepo)

	contactRepo := services.NewContactRepository()
	contactService := services.NewContactService(contactRepo)

	// Image
	imageService := services.NewImageService()

	// Admin
	adminHandler := handlers.NewAdminHandler(postService, imageService, contactService)

	// Comments
	commentService := services.NewCommentService()

	// Blog
	blogHandler := handlers.NewBlogHandler(postService, commentService)

	// =========================
	// PUBLIC ROUTES
	// =========================

	r.HandleFunc("/", handlers.Home).Methods("GET")
	r.HandleFunc("/index2", handlers.Home2).Methods("GET")
	r.HandleFunc("/contact", handlers.Contact).Methods("POST")

	// Auth
	r.HandleFunc("/login", authHandler.Login)
	r.HandleFunc("/logout", authHandler.Logout)

	// =========================
	// BLOG ROUTES (LOGGED-IN USERS)
	// =========================

	blog := r.PathPrefix("/blog").Subrouter()
	blog.Use(auth.RequireLogin)

	blog.HandleFunc("", handlers.Blog) // /blog
	blog.HandleFunc("/{slug}", blogHandler.BlogPost)
	blog.HandleFunc("/{slug}/comment", blogHandler.CreateComment).Methods("POST")
	blog.HandleFunc("/comment/delete", blogHandler.DeleteComment).Methods("POST")

	// =========================
	// ADMIN ROUTES
	// =========================

	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(auth.RequireLogin)
	admin.Use(auth.RequireAdmin)

	// posts
	admin.HandleFunc("/posts", adminHandler.AdminPosts)
	admin.HandleFunc("/posts/create", adminHandler.AdminCreatePost)
	admin.HandleFunc("/posts/edit/{id}", adminHandler.AdminEditPost)
	admin.HandleFunc("/posts/delete/{id}", adminHandler.AdminDeletePost).Methods("POST")
	

	// contact
	admin.HandleFunc("/contact", adminHandler.ManageQueries)
	admin.HandleFunc("/contact/read", adminHandler.MarkContactRead).Methods("POST")
	admin.HandleFunc("/contact/delete", adminHandler.DeleteContact).Methods("POST")

	// register (admin only now)
	admin.HandleFunc("/register", authHandler.Register)

	// =========================
	// STATIC FILES
	// =========================

	r.PathPrefix("/style/").Handler(
		http.StripPrefix("/style/",
			http.FileServer(http.Dir("./static/slowave/style")),
		),
	)

	r.PathPrefix("/uploads/").Handler(
		http.StripPrefix("/uploads/",
			http.FileServer(http.Dir("./static/uploads")),
		),
	)

	// =========================
	// START SERVER
	// =========================

	http.ListenAndServe(":8080", r)
}