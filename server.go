package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
)

func StartServer(port string) {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.Logger)
	})

	server := http.Server{
		Handler: r,
	}

	log.Fatal(server.ListenAndServe())

}
