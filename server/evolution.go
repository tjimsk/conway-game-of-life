package main

import (
	"time"
)

type Evolution struct {
	Cells      []*Cell       `json:"c"`
	Generation int           `json:"g"`
	Duration   time.Duration `json:"d"`
	Interval   time.Duration `json:"i"`
}

func NewEvolution(cells []*Cell, gen int, d time.Duration, i time.Duration) (e Evolution) {
	return Evolution{
		Cells:      cells,
		Generation: gen,
		Duration:   d,
		Interval:   i,
	}
}
