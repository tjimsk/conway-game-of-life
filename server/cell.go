package main

import ()

type Cell struct {
	Active bool  `json:"a"`
	Color  Color `json:"c"`
	X      int   `json:"x"`
	Y      int   `json:"y"`

	active bool
	color  Color
	x      int
	y      int
}

func NewCell(p Point) *Cell {
	return &Cell{
		X: p.X,
		Y: p.Y,
	}
}

func (c *Cell) Flush() {
	c.Active = c.active
	c.Color = c.color
}
