package main

import (
	"encoding/json"
	"github.com/madxmike/go-craigslist"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
)

type CraigslistCity struct {
	Abbreviation     string  `json:"Abbreviation"`
	AreaID           int     `json:"AreaID"`
	Country          string  `json:"Country"`
	Description      string  `json:"Description"`
	Hostname         string  `json:"Hostname"`
	Latitude         float64 `json:"Latitude"`
	Longitude        float64 `json:"Longitude"`
	Region           string  `json:"Region"`
	ShortDescription string  `json:"ShortDescription"`
	SubAreas         []struct {
		Abbreviation     string `json:"Abbreviation"`
		Description      string `json:"Description"`
		ShortDescription string `json:"ShortDescription"`
		SubAreaID        int    `json:"SubAreaID"`
	} `json:"SubAreas,omitempty"`
	Timezone string `json:"Timezone"`
}

func LoadAllCities() ([]CraigslistCity, error) {
	cities := make([]CraigslistCity, 0)
	resp, err := http.Get("https://reference.craigslist.org/Areas")
	if err != nil {
		return cities, errors.Wrap(err, "could not load city data")

	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&cities)
	if err != nil {
		return cities, errors.Wrap(err, "could not parse city data")
	}
	return cities, nil
}

func FindAllCitiesWithin(cities []CraigslistCity, bounds []float64) []CraigslistCity {
	swLong := bounds[0]
	swLat := bounds[1]
	neLong := bounds[2]
	neLat := bounds[3]
	within := make([]CraigslistCity, 0, cap(cities))
	for _, city := range cities {
		eastBound := city.Longitude < neLong
		westBound := city.Longitude > swLong

		var inLong bool
		if neLong < swLong {
			inLong = eastBound || westBound
		} else {
			inLong = eastBound && westBound
		}

		inLat := city.Latitude > swLat && city.Latitude < neLat
		if inLat && inLong {
			within = append(within, city)
		}
	}
	return within
}

type CraigslistHarvester struct {
	options searchOptions
	cities  []CraigslistCity
}

func (h *CraigslistHarvester) Harvest() ([]SearchResult, error) {
	min, err := strconv.Atoi(h.options.MinPrice)
	if err != nil {
		return nil, err
	}
	max, err := strconv.Atoi(h.options.MaxPrice)
	if err != nil {
		return nil, err
	}
	results := make([]SearchResult, 0)
	opts := craigslist.SearchOptions{
		Category: "sss",
		Query:    h.options.Query,
		MinPrice: min,
		MaxPrice: max,
	}
	for _, city := range h.cities {
		result, err := craigslist.Search(city.Hostname, opts)
		if err != nil {
			return results, errors.Wrap(err, "could not load craigslist data")
		}
		for _, listing := range result.Listings {
			got, err := craigslist.GetListing(listing.URL)
			if err != nil {
				log.Println(err)
				continue
			}
			if got.Location == nil {
				continue
			}
			results = append(results, SearchResult{
				Vendor:    "Craigslist",
				Title:     got.Title,
				Posted:    got.PostedAt,
				Price:     strconv.Itoa(got.Price),
				Latitude:  got.Location.Lat,
				Longitude: got.Location.Lng,
			})
		}
	}
	log.Println(len(results))

	return results, nil
}
