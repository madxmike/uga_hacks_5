package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
		r.Route("/about", func(r chi.Router) {
			r.Get("/", ServeTemplate("about"))
		})

		r.Route("/api", func(r chi.Router) {
			cities, err := LoadAllCities()
			if err != nil {
				panic(err)
			}
			r.Method("POST", "/search", &SearchHandler{craiglistCities: cities})
		})

	})
	r.Mount("/static/", ServeStatic("/static/"))

	server := http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	panic(server.ListenAndServe())
}
