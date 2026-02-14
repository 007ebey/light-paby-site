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
	r.HandleFunc("/index2", handlers.Home2)
	r.HandleFunc("/admin/sliders/create", handlers.SliderAdmin)

	r.PathPrefix("/style/").Handler(http.StripPrefix("/style/",
	 http.FileServer(http.Dir("./static/slowave/style")),
	))
	
	http.ListenAndServe(":8080", r)
}