package main

import (
	"fmt"
)

type Cell struct {
	X          int              `json:"x"`
	Y          int              `json:"y"`
	Generation int              `json:"-"`
	Agents     map[string]*User `json:"-"`
	Active     bool             `json:"active"`
	Color      *Color           `json:"color"`
}

func NewCell(x int, y int) *Cell {
	return &Cell{
		X:      x,
		Y:      y,
		Agents: map[string]*User{},
	}
}

func (c *Cell) AverageColor() *Color {
	if c.Color != nil {
		return c.Color
	}

	var (
		sR    int
		sG    int
		sB    int
		count = len(c.Agents)
	)

	for _, u := range c.Agents {
		sR = sR + u.Color.R
		sG = sG + u.Color.G
		sB = sB + u.Color.B
	}

	c.Color = NewColor(sR/count, sG/count, sB/count)

	return c.Color
}

func (c *Cell) adjacentCoordinates() []Coordinate {
	coordinates := []Coordinate{}

	coordinates = append(coordinates,
		Coordinate{c.X - 1, c.Y + 1},
		Coordinate{c.X - 1, c.Y},
		Coordinate{c.X - 1, c.Y - 1},
		Coordinate{c.X, c.Y - 1},
		Coordinate{c.X, c.Y + 1},
		Coordinate{c.X + 1, c.Y - 1},
		Coordinate{c.X + 1, c.Y},
		Coordinate{c.X + 1, c.Y + 1},
	)

	return coordinates
}

func (c *Cell) appendAdjacentCellsAgents(adjCells []*Cell) {
	for _, adjCell := range adjCells {
		for _, agent := range adjCell.Agents {
			c.Agents[agent.Name] = agent
		}
	}
}

type Coordinate struct {
	X int
	Y int
}

func (c Coordinate) String() string {
	return fmt.Sprintf(`%v,%v`, c.X, c.Y)
}
