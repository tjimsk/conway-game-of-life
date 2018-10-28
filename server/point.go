package main

import (
	"fmt"
)

type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf(`%v,%v`, p.X, p.Y)
}

func (p Point) Adjacent() (pts []Point) {
	return append(pts,
		Point{p.X - 1, p.Y + 1},
		Point{p.X - 1, p.Y},
		Point{p.X - 1, p.Y - 1},
		Point{p.X, p.Y - 1},
		Point{p.X, p.Y + 1},
		Point{p.X + 1, p.Y - 1},
		Point{p.X + 1, p.Y},
		Point{p.X + 1, p.Y + 1},
	)
}
