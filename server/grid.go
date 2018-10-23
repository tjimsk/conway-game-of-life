package main

import (
	"fmt"
	"sync"
)

var (
	_ fmt.Stringer
)

type Grid struct {
	Width      int                   `json:"width"`
	Height     int                   `json:"height"`
	Cells      map[int]map[int]*Cell `json:"-"`
	Generation int                   `json:"generation"`
	mu         *sync.Mutex
}

func NewGrid(w, h int) *Grid {
	g := &Grid{
		Width:  w,
		Height: h,
		Cells:  map[int]map[int]*Cell{},
		mu:     &sync.Mutex{},
	}

	// instantiate every cell
	for i := 0; i < w; i++ {
		g.Cells[i] = map[int]*Cell{}

		for j := 0; j < h; j++ {
			g.Cells[i][j] = NewCell(i, j)
		}
	}

	return g
}

func (g *Grid) seed() {
	player0 := NewUser("player0")

	coordinates := []Coordinate{
		// blinker
		Coordinate{10, 10},
		Coordinate{11, 10},
		Coordinate{12, 10},
		// beacon
		Coordinate{10, 15},
		Coordinate{10, 16},
		Coordinate{11, 15},
		Coordinate{12, 18},
		Coordinate{13, 17},
		Coordinate{13, 18},
	}

	for _, coordinate := range coordinates {
		c := grid.cellAtCoordinate(coordinate)
		c.Agents[player0.Name] = player0
		c.Active = true
	}
}

func (g *Grid) cell(x, y int) *Cell {
	return g.Cells[x][y]
}

func (g *Grid) cellAtCoordinate(c Coordinate) *Cell {
	return g.Cells[c.X][c.Y]
}

func (g *Grid) cellsToEvolve() (cells []*Cell) {
	keys := map[string]bool{}
	activeCells := g.activeCells()

	for _, c := range activeCells {
		coordinates := append(c.adjacentCoordinates(), Coordinate{c.X, c.Y})

		for _, coordinate := range coordinates {
			k := coordinate.String()

			if _, exists := keys[k]; !exists {
				cells = append(cells, g.cellAtCoordinate(coordinate))
				keys[k] = true
			}
		}
	}

	return cells
}

func (g *Grid) activeCells() (activeCells []*Cell) {
	for _, row := range g.Cells {
		for _, c := range row {
			if c.Active {
				activeCells = append(activeCells, c)
			}
		}
	}

	return activeCells
}

func (g *Grid) activeAdjacentCells(c *Cell) (adjCells []*Cell) {
	for _, coordinate := range c.adjacentCoordinates() {
		adjCell := g.cellAtCoordinate(coordinate)

		if adjCell.Active {
			adjCells = append(adjCells, adjCell)
		}
	}

	return adjCells
}

func (g *Grid) nextGeneration() (gu GridUpdate) {
	updatedCells := []*Cell{}
	cells := g.cellsToEvolve()

	// evaluate every cell's next evolution
	for _, c := range cells {
		adjCells := g.activeAdjacentCells(c)
		activeCount := len(adjCells)

		// evaluate evolution state on cell copy
		var cellCopy Cell
		cellCopy = *c

		if c.Active {
			if activeCount < 2 { // condition #1: cell dies
				cellCopy.Active = false
				cellCopy.Agents = map[string]*User{}
			} else if activeCount == 2 || activeCount == 3 { //condition #2: cell lives
				cellCopy.Active = true
				cellCopy.appendAdjacentCellsAgents(adjCells)
			} else if activeCount > 3 { // condition #3: cell dies
				cellCopy.Active = false
				cellCopy.Agents = map[string]*User{}
			}
		} else {
			if activeCount == 3 { // condition #4: cell lives
				cellCopy.Active = true
				cellCopy.appendAdjacentCellsAgents(adjCells)
			}
		}

		updatedCells = append(updatedCells, &cellCopy)
	}

	// apply evolution updates
	for i, cellCopy := range updatedCells {
		c := g.cell(cellCopy.X, cellCopy.Y)
		c.Active = cellCopy.Active

		if c.Active {
			c.Color = cellCopy.AverageColor()
		}

		updatedCells[i] = c
	}

	g.Generation++

	return g.updateFromCells(updatedCells)
}

type CellUpdate struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Active bool   `json:"a"`
	Color  *Color `json:"c"`
}

type GridUpdate struct {
	Generation int          `json:"generation"`
	Cells      []CellUpdate `json:"cells"`
}

func (g *Grid) updateFromCells(updatedCells []*Cell) (gu GridUpdate) {
	gu.Generation = g.Generation

	for _, c := range updatedCells {
		gu.Cells = append(gu.Cells, CellUpdate{
			X:      c.X,
			Y:      c.Y,
			Active: c.Active,
			Color:  c.Color,
		})
	}

	return gu
}
