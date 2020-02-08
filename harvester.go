package main

import (
	"github.com/pkg/errors"
	"github.com/sosedoff/go-craigslist"
	"io"
)

type Harvester interface {
	Harvest(writer io.Writer) error
}

type CraigslistHarvester struct {
	options searchOptions
}

func (h *CraigslistHarvester) Harvest(writer io.Writer) error {
	opts := craigslist.SearchOptions{
		Query:    h.options.Query,
		MinPrice: h.options.MinPrice,
		MaxPrice: h.options.MaxPrice,
	}
	results, err := craigslist.Search(h.options.Location, opts)
	if err != nil {
		return errors.Wrap(err, "could not load craigslist data")
	}
	var json string
	for _, listings := range results.Listings {
		json, err = listings.JSON()
		_, _ = writer.Write([]byte(json))
	}
	return nil
}
