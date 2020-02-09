package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

type SearchHandler struct {
	craiglistCities []CraigslistCity
}

type searchOptions struct {
	UseCraigslist string `json:"use_craigslist"`
	Query         string `json:"query"`
	Bounds        string `json:"bounds"`
	MinPrice      string `json:"price_min"`
	MaxPrice      string `json:"price_max"`
}

type SearchResult struct {
	Vendor      string  `json:"vendor"`
	Title       string  `json:"title"`
	Posted      string  `json:"posted"`
	Price       string  `json:"price"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	URL         string  `json:"url"`
}

func (h *SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	opts, err := h.parseForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	harvesters := make([]Harvester, 0)
	if opts.UseCraigslist == "on" {
		bounds, err := h.parseBounds(opts.Bounds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
		harvesters = append(harvesters, &CraigslistHarvester{
			options: opts,
			cities:  FindAllCitiesWithin(h.craiglistCities, bounds),
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

func (h *SearchHandler) parseForm(r *http.Request) (searchOptions, error) {
	var searchOptions searchOptions
	err := json.NewDecoder(r.Body).Decode(&searchOptions)
	if err != nil {
		return searchOptions, errors.Wrap(err, "could not parse form")
	}
	return searchOptions, nil
}

func (h *SearchHandler) parseBounds(bounds string) ([]float64, error) {
	parsedBounds := make([]float64, 0, 4)
	split := strings.Split(bounds, ",")
	for _, bound := range split {
		parsed, err := strconv.ParseFloat(bound, 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse bounds")
		}
		parsedBounds = append(parsedBounds, parsed)
	}
	return parsedBounds, nil
}
