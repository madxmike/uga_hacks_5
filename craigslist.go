package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/umahmood/haversine"
	"net/http"
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
	err = json.NewDecoder(resp.Body).Decode(cities)
	if err != nil {
		return cities, errors.Wrap(err, "could not parse city data")
	}
	return cities, nil
}

func FindAllCitiesWithinFrom(cities []CraigslistCity, miles float64, lat float64, long float64) []CraigslistCity {
	within := make([]CraigslistCity, 0, cap(cities))
	for _, city := range cities {
		cityPoint := haversine.Coord{
			Lat: city.Latitude,
			Lon: city.Longitude,
		}

		fromPoint := haversine.Coord{
			Lat: lat,
			Lon: long,
		}
		milesFrom, _ := haversine.Distance(cityPoint, fromPoint)
		if milesFrom <= miles {
			within = append(within, city)
		}
	}
	return within
}
