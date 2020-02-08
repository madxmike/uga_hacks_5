package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func TemplateLoader(path string, liveMode bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t, ok := r.Context().Value("template").(*template.Template)
			if !ok || liveMode {
				var err error
				t, err = template.New("template").Funcs(template.FuncMap{
					"derefString": func(s *string) string { return *s },
				}).ParseGlob(path)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			ctx := context.WithValue(r.Context(), "template", t)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ServeTemplate(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("serving " + name)
		t, ok := r.Context().Value("template").(*template.Template)
		if !ok {
			http.Error(w, fmt.Sprintf("could not load template (%s) from context", name), http.StatusInternalServerError)
			return
		}
		err := t.ExecuteTemplate(w, name, r.Context().Value("data"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ServeStatic(path string) http.HandlerFunc {
	fs := http.FileServer(http.Dir("." + path))
	fs = http.StripPrefix(path, fs)
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	}
}
