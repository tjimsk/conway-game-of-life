package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var (
	_                 fmt.Stringer
	ErrCellOutOfBound = errors.New("Cell point is out of bounds")
)

type Grid struct {
	Generation int `json:"g"`
	Height     int `json:"h"`
	Width      int `json:"w"`

	cells      map[int]map[int]*Cell
	evoChan    chan Evolution
	updateChan chan []*Cell
}

func NewGrid(w int, h int) *Grid {
	g := &Grid{
		Width:  w,
		Height: h,

		cells:      map[int]map[int]*Cell{},
		evoChan:    make(chan Evolution),
		updateChan: make(chan []*Cell),
	}

	for i := 0; i < w; i++ {
		g.cells[i] = map[int]*Cell{}

		for j := 0; j < h; j++ {
			g.cells[i][j] = NewCell(Point{i, j})
		}
	}

	return g
}

func (g *Grid) CellAtPoint(p Point) (c *Cell, err error) {
	_, ok := g.cells[p.X]
	if !ok {
		return nil, ErrCellOutOfBound
	}

	c, ok = g.cells[p.X][p.Y]
	if !ok {
		return nil, ErrCellOutOfBound
	}

	return c, nil
}

func (g *Grid) activeCells() (cells []*Cell) {
	for _, cellRow := range g.cells {
		for _, c := range cellRow {
			if c.Active {
				cells = append(cells, c)
			}
		}
	}

	return cells
}

// An unstable cell is one located in a position adjacent to an active cell.
// When evolving an entire grid, this function is called first to determine
// the cells on which to apply conway's 4 rules of evolution.
func (g *Grid) unstableCells() (cells []*Cell) {
	aCells := g.activeCells()
	keys := map[string]bool{}

	for _, c := range aCells {
		p := Point{c.X, c.Y}
		pts := append(p.Adjacent(), p)

		for _, p := range pts {
			k := p.String()

			if _, exists := keys[k]; !exists {
				_c, err := g.CellAtPoint(p)
				if err == nil {
					cells = append(cells, _c)
				}

				keys[k] = true
			}
		}
	}

	return cells
}

func (g *Grid) evolveCell(c *Cell) {
	p := Point{c.X, c.Y}
	aCells := []*Cell{}

	for _, ap := range p.Adjacent() {
		ac, err := g.CellAtPoint(ap)
		if err != nil {
			continue
		}

		if ac.Active {
			aCells = append(aCells, ac)
		}
	}

	activeCount := len(aCells)

	switch {
	// condition #1 & #3: cell dies
	case c.Active && (activeCount < 2 || activeCount > 3):
		c.active = false
		c.color = Color{}
	//condition #2: cell stays alive
	case c.Active && (activeCount == 2 || activeCount == 3):
		c.active = true
		c.color = c.Color
	// condition #4: cell comes to life
	case !c.Active && activeCount == 3:
		c.active = true

		colors := []Color{}
		for _, _c := range aCells {
			colors = append(colors, _c.Color)
		}

		c.color = AverageColor(colors)
	}
}

func (g *Grid) StartEvolutions() {
	var cells = []*Cell{}

	for {
		t0 := time.Now()

		g.Generation++

		cells = g.unstableCells()

		for _, c := range cells {
			g.evolveCell(c)
		}

		for _, c := range cells {
			c.Flush()
		}

		t1 := time.Now()
		de := t1.Sub(t0)
		di := time.Duration(viper.GetInt("evoInterval")) * time.Millisecond

		g.evoChan <- NewEvolution(cells, g.Generation, de, di)

		time.Sleep(di)
	}
}
