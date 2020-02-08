package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func StartServer(port string) {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(TemplateLoader("./static/pages/*", true))
		r.Use(middleware.DefaultCompress)
		r.Get("/", ServeTemplate("index"))
		r.Route("/view", func(r chi.Router) {
			r.Get("/", ServeTemplate("view"))
		})
	})
	r.Mount("/static/", ServeStatic("/static/"))

	server := http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	panic(server.ListenAndServe())
}
