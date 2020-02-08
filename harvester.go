package main

import (
	"io"
)

type Harvester interface {
	Harvest(writer io.Writer) error
}
