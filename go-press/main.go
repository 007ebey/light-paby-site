package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"word_press/auth"
	"word_press/database"
	"word_press/handlers"
	"word_press/services"
	"word_press/realtime"
)

func main() {
	database.InitDB()

	r := mux.NewRouter()

	// =========================
	// SERVICES + HANDLERS
	// =========================

	// Posts
	postRepo := services.NewPostRepository()
	postService := services.NewPostService(postRepo)

	// Auth
	userRepo := services.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService, postService)

	contactRepo := services.NewContactRepository()
	contactService := services.NewContactService(contactRepo)

	presenceManager := realtime.NewPresenceManager(6000)

	prayerService := services.NewPrayerService()
	presenceService := services.NewPresenceService(presenceManager)
	

	// Image
	imageService := services.NewImageService()

	// Admin
	adminHandler := handlers.NewAdminHandler(postService, imageService, contactService)

	// Comments
	commentService := services.NewCommentService()

	// Blog
	blogHandler := handlers.NewBlogHandler(postService, commentService)

	pageHandler := handlers.NewPageHandler(postService, contactService)

	prayerHandler := handlers.NewPrayerHandler(prayerService, presenceService)

	// =========================
	// PUBLIC ROUTES
	// =========================

	r.HandleFunc("/", pageHandler.Home).Methods("GET")
	r.HandleFunc("/contact", pageHandler.Contact).Methods("POST")

	// Auth
	r.HandleFunc("/login", authHandler.Login)
	r.HandleFunc("/logout", authHandler.Logout)

	// =========================
	// BLOG ROUTES (LOGGED-IN USERS)
	// =========================

	blog := r.PathPrefix("/blog").Subrouter()
	blog.Use(auth.RequireLogin)

	blog.HandleFunc("", pageHandler.Blog) // /blog
	blog.HandleFunc("/{slug}", blogHandler.BlogPost)
	blog.HandleFunc("/{slug}/comment", blogHandler.CreateComment).Methods("POST")
	blog.HandleFunc("/{slug}/comment/delete", blogHandler.DeleteComment).Methods("POST")

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

	r.HandleFunc("/vision-notes", prayerHandler.VisionNotes)
    r.HandleFunc("/vision/create", prayerHandler.CreateVisionNote)

	prayer := r.PathPrefix("/prayer").Subrouter()

	prayer.HandleFunc("", prayerHandler.Home)
    prayer.HandleFunc("/start", prayerHandler.Start).Methods("POST")
    prayer.HandleFunc("/end", prayerHandler.End).Methods("POST")
    prayer.HandleFunc("/heartbeat", prayerHandler.Heartbeat).Methods("POST")


    prayer.HandleFunc("/list", prayerHandler.PrayerList)
	prayer.HandleFunc("/list/create", prayerHandler.CreatePrayer).Methods("POST")
	prayer.HandleFunc("/list/delete", prayerHandler.DeletePrayer).Methods("POST")

    r.HandleFunc("/vision-notes", prayerHandler.VisionNotes)
    r.HandleFunc("/vision/create", prayerHandler.CreateVisionNote)
	
    prayer.HandleFunc("/fellow-prayer", prayerHandler.Fellow)

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