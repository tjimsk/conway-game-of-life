package main

import (
	"testing"
)

func TestEvolutionResults(t *testing.T) {
	g := NewGrid(100, 100)
	// u := g.NewUser(nil)

	type state struct {
		generation    int
		unstableCount int
		activePoints  []Point
	}

	type testData struct {
		initialPoints []Point
		states        []state
	}

	data := testData{
		initialPoints: []Point{
			Point{1, 4},
			Point{2, 4},
			Point{3, 4},
		},
		states: []state{
			state{
				generation:    1,
				unstableCount: 15,
				activePoints: []Point{
					Point{2, 3},
					Point{2, 4},
					Point{2, 5},
				},
			},
			state{
				generation:    2,
				unstableCount: 15,
				activePoints: []Point{
					Point{1, 4},
					Point{2, 4},
					Point{3, 4},
				},
			},
		},
	}

	for _, p := range data.initialPoints {
		c, err := g.CellAtPoint(p)
		if err != nil {
			t.Fail()
		}
		c.Active = true
	}

	for _, s := range data.states {
		evo := g.evolve()

		if s.generation != evo.Generation {
			t.Fail()
		}

		if len(evo.Cells) != s.unstableCount {
			t.Fail()
		}

		for _, p := range s.activePoints {
			c, err := g.CellAtPoint(p)
			if err != nil {
				t.Fail()
			}

			if !c.Active {
				t.Fail()
			}
		}
	}
}
