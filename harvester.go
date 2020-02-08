package main

type Harvester interface {
	Harvest() ([]SearchResult, error)
}
