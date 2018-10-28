package main

func (g *Grid) seed() {
	player0 := NewUser(nil) // create a seed user since it needs a random color

	points := []Point{
		// blinker
		Point{2, 3},
		Point{3, 3},
		Point{4, 3},
		// beacon
		Point{7, 2},
		Point{7, 3},
		Point{8, 2},
		Point{9, 5},
		Point{10, 4},
		Point{10, 5},
	}

	for _, point := range points {
		c, err := grid.CellAtPoint(point)
		if err != nil {
			continue
		}

		c.Active = true
		c.Color = player0.Color
	}
}
