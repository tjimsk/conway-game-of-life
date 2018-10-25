package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	_ fmt.Stringer
)

type Grid struct {
	Width       int                   `json:"width"`
	Height      int                   `json:"height"`
	Cells       map[int]map[int]*Cell `json:"-"`
	AutoEvo     bool                  `json:"autoEvo"`
	EvoInterval time.Duration         `json:"evoInterval"`
	Generation  int                   `json:"generation"`
	mu          *sync.Mutex
}

func NewGrid(w int, h int, autoEvo bool, evoInterval time.Duration) *Grid {
	g := &Grid{
		Width:       w,
		Height:      h,
		Cells:       map[int]map[int]*Cell{},
		AutoEvo:     autoEvo,
		EvoInterval: evoInterval,
		mu:          &sync.Mutex{},
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
	player0 := NewUser(nil) // create a seed user since it needs a random color

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
	_, ok := g.Cells[x]
	if !ok {
		return nil
	}

	cell, ok := g.Cells[x][y]
	if !ok {
		return nil
	}

	return cell
}

func (g *Grid) cellAtCoordinate(c Coordinate) *Cell {
	return g.Cells[c.X][c.Y]
}

func (g *Grid) cellsToEvolve() (cells []*Cell) {
	keys := map[string]bool{}
	activeCells := g.activeCells()

	for _, c := range activeCells {
		coordinates := append(g.adjacentCoordinates(c), Coordinate{c.X, c.Y})

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

func (g *Grid) adjacentCoordinates(c *Cell) []Coordinate {
	coordinates := []Coordinate{}

	if c == nil {
		return coordinates
	}

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

func (g *Grid) activeAdjacentCells(c *Cell) (adjCells []*Cell) {
	for _, coordinate := range g.adjacentCoordinates(c) {
		adjCell := g.cellAtCoordinate(coordinate)

		if adjCell != nil && adjCell.Active {
			adjCells = append(adjCells, adjCell)
		}
	}

	return adjCells
}

func (g *Grid) bumpGeneration() {
	g.Generation++

	for _, row := range g.Cells {
		for _, c := range row {
			c.Generation++
		}
	}
}

func (g *Grid) ActivateCells(msgCells map[string]Cell, u *User) (cells []*Cell) {
	for _, _c := range msgCells {
		c := g.cell(_c.X, _c.Y)

		if c == nil {
			continue
		}

		c.Active = _c.Active
		if c.Active {
			// add user to cell's agents and set cell's color to user's
			c.Agents = map[string]*User{}
			c.Agents[u.Name] = u
			c.Color = u.Color
		}

		cells = append(cells, c)
	}

	return cells
}

func (g *Grid) Evolve() (gu GridUpdate) {
	updatedCells := []*Cell{}
	cells := g.cellsToEvolve()

	g.bumpGeneration()

	// evaluate every cell's next evolution
	for _, c := range cells {
		if c == nil {
			continue
		}

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

	return g.UpdateFromCells(updatedCells)
}

type GridUpdate struct {
	Generation int          `json:"generation"`
	Cells      []CellUpdate `json:"cells"`
}

func (g *Grid) UpdateFromCells(updatedCells []*Cell) (gu GridUpdate) {
	gu.Generation = g.Generation

	for _, c := range updatedCells {
		gu.Cells = append(gu.Cells, CellUpdate{
			X:          c.X,
			Y:          c.Y,
			Active:     c.Active,
			Color:      c.Color,
			Generation: c.Generation,
		})
	}

	return gu
}
