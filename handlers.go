package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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

type SearchHandler struct {
	craiglistCities []CraigslistCity
}

type searchOptions struct {
	UseCraigslist bool
	Query         string
	Lat           float64
	Long          float64
	Miles         float64
	MinPrice      int
	MaxPrice      int
}

func (h *SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	opts, err := h.parseForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	harvesters := make([]Harvester, 0)

	if opts.UseCraigslist {
		harvesters = append(harvesters, &CraigslistHarvester{
			options: *opts,
			cities:  FindAllCitiesWithinFrom(h.craiglistCities, opts.Miles, opts.Lat, opts.Long),
		})
	}
}

func (h *SearchHandler) parseForm(r *http.Request) (*searchOptions, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, errors.Wrap(err, "could not parse search form")
	}

	minPrice, err := strconv.Atoi(r.FormValue("price_min"))
	if err != nil {
		return nil, errors.Wrap(err, "could not parse price_min")
	}
	maxPrice, err := strconv.Atoi(r.FormValue("price_max"))
	if err != nil {
		return nil, errors.Wrap(err, "could not parse price_max")
	}
	miles, err := strconv.Atoi(r.FormValue("miles"))
	if err != nil {
		return nil, errors.Wrap(err, "could not parse miles")
	}
	lat, err := strconv.Atoi(r.FormValue("lat"))
	if err != nil {
		return nil, errors.Wrap(err, "could not parse lat")
	}
	long, err := strconv.Atoi(r.FormValue("long"))
	if err != nil {
		return nil, errors.Wrap(err, "could not parse long")
	}
	useCraigslist, err := strconv.ParseBool(r.FormValue("use_craigslist"))
	if err != nil {
		return nil, errors.Wrap(err, "could not parse use_craigslist")
	}
	return &searchOptions{
		UseCraigslist: useCraigslist,
		Query:         r.FormValue("query"),
		Lat:           float64(lat),
		Long:          float64(long),
		Miles:         float64(miles),
		MinPrice:      minPrice,
		MaxPrice:      maxPrice,
	}, nil
}

type SearchResult struct {
	Vendor    string    `json:"vendor"`
	Title     string    `json:"title"`
	Posted    time.Time `json:"posted"`
	Price     string    `json:"price"`
	Latitude  string    `json:"latitude"`
	Longitude string    `json:"longitude"`
}
