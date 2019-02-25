package life

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Grid struct {
	Generation int
	Interval   int
	players    map[string]Player
	liveCells  map[Point]struct{}
	cells      map[Point]Cell
	mu         sync.Mutex
}

func NewGrid(interval int) *Grid {
	return &Grid{
		Interval:  interval,
		players:   map[string]Player{},
		liveCells: map[Point]struct{}{},
		cells:     map[Point]Cell{},
	}
}

func (g *Grid) Evolve() {
	var (
		liveCells      = g.liveCellsCopy()
		cells          = g.cellsCopy()
		unstablePoints = map[Point]bool{}
	)

	// add all unstable points
	for p, _ := range liveCells {
		unstablePoints[p] = false
		for _, ap := range adjacentPoints(p) {
			unstablePoints[ap] = false
		}
	}

	// compute next generation cell state
	for p, _ := range unstablePoints {
		count, colors := 0, []Color{}

		for _, ap := range adjacentPoints(p) {
			if _, ok := liveCells[ap]; ok {
				count++
				colors = append(colors, cells[ap].Color)
			}
		}
		_, live := liveCells[p]

		switch {
		case live && (count < 2 || count > 3): // condition #1 & #3: cell dies
			unstablePoints[p] = false
		case live && (count == 2 || count == 3): //condition #2: cell stays alive
			unstablePoints[p] = true
		case !live && count == 3: // condition #4: cell comes to life; average cell colors
			unstablePoints[p] = true
			cells[p] = Cell{averageColor(colors), p}
		}
	}

	for p, live := range unstablePoints {
		if live && p.X < 100 && p.Y < 100 && p.X > -2 && p.Y > -2 {
			liveCells[p] = struct{}{}
		} else {
			delete(liveCells, p)
		}
	}

	g.setLiveCells(liveCells)
	g.setCells(cells)

	g.Generation++
}

func (g *Grid) Activate(points []Point, player Player) {
	liveCells := g.liveCellsCopy()
	cells := g.cellsCopy()

	for _, p := range points {
		liveCells[p] = struct{}{}
		cells[p] = Cell{player.Color, p}
	}

	g.setLiveCells(liveCells)
	g.setCells(cells)
}

func (g *Grid) Deactivate(p Point, player Player) {
	g.mu.Lock()
	delete(g.liveCells, p)
	g.mu.Unlock()
}

func (g *Grid) AddPlayer(conn *websocket.Conn) (player Player) {
	player = NewPlayer(conn)
	g.players[player.Name] = player

	return player
}

func (g *Grid) RemovePlayer(player Player) {
	delete(g.players, player.Name)
	log.Printf("%v removed", player.Name)
}

func (g *Grid) PlayerConnected(player Player) bool {
	_, ok := g.players[player.Name]
	return ok
}

func (g *Grid) NoConnectedUser() bool {
	return len(g.players) == 0
}

func (g *Grid) SetInterval(interval int, player Player) {
	g.mu.Lock()
	g.Interval = interval
	g.mu.Unlock()
}

func (g *Grid) Reset(player Player) {
	g.mu.Lock()
	g.liveCells = map[Point]struct{}{}
	g.Generation = 0
	g.mu.Unlock()
}

// accessors
func (g *Grid) liveCellsCopy() map[Point]struct{} {
	liveCells := map[Point]struct{}{}
	g.mu.Lock()
	for p, _ := range g.liveCells {
		liveCells[p] = struct{}{}
	}
	g.mu.Unlock()
	return liveCells
}

func (g *Grid) cellsCopy() map[Point]Cell {
	cells := map[Point]Cell{}
	g.mu.Lock()
	for p, c := range g.cells {
		cells[p] = c
	}
	g.mu.Unlock()
	return cells
}

func (g *Grid) setLiveCells(liveCells map[Point]struct{}) {
	g.mu.Lock()
	g.liveCells = liveCells
	g.mu.Unlock()
}

func (g *Grid) setCells(cells map[Point]Cell) {
	g.mu.Lock()
	g.cells = cells
	g.mu.Unlock()
}

// message push to websocket
func (g *Grid) PushStateChange() {
	msg := NewPushStateMessage(g)
	for _, player := range g.players {
		player.messageChan <- msg
	}
}

func (g *Grid) PushIntervalChange(changedBy Player) {
	msg := NewPushIntervalMessage(g, changedBy)
	for _, player := range g.players {
		player.messageChan <- msg
	}
}

func (g *Grid) PushPlayer(player Player) {
	msg := NewPushPlayerMessage(player)
	player.messageChan <- msg
}

func (g *Grid) PushState(player Player) {
	msg := NewPushStateMessage(g)
	player.messageChan <- msg
}

// cell
type Cell struct {
	Color Color `json:"c"`
	Point Point `json:"p"`
}

type Point struct {
	X int
	Y int
}

func adjacentPoints(p Point) (pts []Point) {
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
