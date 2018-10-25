package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Color struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

func NewColor(R, G, B int) *Color {
	return &Color{R, G, B}
}

func NewRandomColor() *Color {
	var (
		R = rand.Intn(255)
		G = rand.Intn(255)
		B = rand.Intn(255)
	)

	return NewColor(R, G, B)
}

func (c *Color) String() string {
	return fmt.Sprintf(`(%v, %v, %v)`, c.R, c.G, c.B)
}
