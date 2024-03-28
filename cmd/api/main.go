package main

import (
	"log"
	"marketplace/internal/config"
	"marketplace/internal/handler"
	"marketplace/internal/middleware"
	"marketplace/internal/store"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	db, err := store.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	h := handler.NewHandler(db)
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.Path("/").HandlerFunc(h.Home).Methods("GET")
	r.HandleFunc("/register", middleware.OptionalAuth(h.Register, db)).Methods("POST")
	r.HandleFunc("/login", middleware.OptionalAuth(h.Login, db)).Methods("POST")
	r.HandleFunc("/ads", middleware.Auth(h.PostAd, db)).Methods("POST")
	r.HandleFunc("/ads", middleware.OptionalAuth(h.GetAds, db)).Methods("GET")

	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
