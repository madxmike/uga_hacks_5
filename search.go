package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

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

type SearchResult struct {
	Vendor    string     `json:"vendor"`
	Title     string     `json:"title"`
	Posted    *time.Time `json:"posted"`
	Price     string     `json:"price"`
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
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

	results := make([]SearchResult, 0)
	for _, harvester := range harvesters {
		harvest, err := harvester.Harvest()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, harvest...)
	}
	b, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
