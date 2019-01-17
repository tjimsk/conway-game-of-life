package life

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

type Grid struct {
	Generation int
	Height     int
	Width      int
	Paused     bool
	players    map[string]Player
	state      map[int]map[int]Cell
}

func NewGrid(width int, height int) *Grid {
	g := &Grid{
		Width:   width,
		Height:  height,
		players: map[string]Player{},
		state:   newState(width, height),
	}

	return g
}

func newState(width, height int) (state map[int]map[int]Cell) {
	state = map[int]map[int]Cell{}
	for x := 0; x < width; x++ {
		state[x] = map[int]Cell{}
		for y := 0; y < height; y++ {
			state[x][y] = Cell{
				X: x,
				Y: y,
			}
		}
	}
	return state
}

func newStateCopy(state map[int]map[int]Cell) (stateCopy map[int]map[int]Cell) {
	stateCopy = newState(len(state), len(state[0]))
	for x, row := range state {
		stateCopy[x] = map[int]Cell{}
		for y, _ := range row {
			stateCopy[x][y] = state[x][y]
		}
	}
	return stateCopy
}

func (g *Grid) Activate(points []Point, player Player) {
	state := newStateCopy(g.state)
	for _, p := range points {
		if p.X*p.Y < 0 || p.X >= g.Width || p.Y >= g.Height {
			continue
		}
		c := state[p.X][p.Y]
		c.Active = true
		c.Color = player.Color
		state[p.X][p.Y] = c
	}
	g.state = state
}

func (g *Grid) Evolve() {
	currentState := g.state
	nextState := newStateCopy(currentState)
	g.Generation++
	// get active points
	activePoints := []Point{}
	for x, col := range currentState {
		for y, c := range col {
			if c.Active {
				activePoints = append(activePoints, Point{x, y})
			}
		}
	}
	// get unstable points
	unstablePoints := []Point{}
	keys := map[string]bool{}
	for _, p := range activePoints {
		adjacentPoints := []Point{
			p,
			Point{p.X - 1, p.Y + 1},
			Point{p.X - 1, p.Y},
			Point{p.X - 1, p.Y - 1},
			Point{p.X, p.Y - 1},
			Point{p.X, p.Y + 1},
			Point{p.X + 1, p.Y - 1},
			Point{p.X + 1, p.Y},
			Point{p.X + 1, p.Y + 1},
		}
		for _, p := range adjacentPoints {
			if p.X*p.Y < 0 || p.X > g.Width || p.Y > g.Height {
				continue
			}
			k := fmt.Sprintf("%v,%v", p.X, p.Y)
			_, ok := keys[k]
			if ok {
				continue
			}
			keys[k] = true
			unstablePoints = append(unstablePoints, p)
		}
	}
	// check next evolution state of each unstable cell
	for _, p := range unstablePoints {
		if p.X*p.Y < 0 || p.X >= g.Width || p.Y >= g.Height {
			continue
		}
		adjacentPoints := []Point{
			Point{p.X - 1, p.Y + 1},
			Point{p.X - 1, p.Y},
			Point{p.X - 1, p.Y - 1},
			Point{p.X, p.Y - 1},
			Point{p.X, p.Y + 1},
			Point{p.X + 1, p.Y - 1},
			Point{p.X + 1, p.Y},
			Point{p.X + 1, p.Y + 1},
		}
		c := currentState[p.X][p.Y]
		active := c.Active
		count := 0
		activeCells := []Cell{}
		for _, ap := range adjacentPoints {
			if ap.X*ap.Y < 0 || ap.X >= g.Width || ap.Y >= g.Height {
				continue
			}
			if currentState[ap.X][ap.Y].Active {
				count++
				activeCells = append(activeCells, currentState[ap.X][ap.Y])
			}
		}
		switch {
		case active && (count < 2 || count > 3): // condition #1 & #3: cell dies
			c := nextState[p.X][p.Y]
			c.Active = false
			c.Color = Color{}
			nextState[p.X][p.Y] = c
		case active && (count == 2 || count == 3): //condition #2: cell stays alive
			c := nextState[p.X][p.Y]
			c.Active = true
			c.Color = c.Color
			nextState[p.X][p.Y] = c
			// log.Println("cell stays alive", p)
		case !active && count == 3: // condition #4: cell comes to life; average cell colors
			c := nextState[p.X][p.Y]
			c.Active = true
			colors := []Color{}
			for _, _c := range activeCells {
				colors = append(colors, _c.Color)
			}
			c.Color = averageColor(colors)
			nextState[p.X][p.Y] = c
			// log.Println("cell comes to life!", p)
		}
	}
	g.state = nextState
}

func (g *Grid) Deactivate(p Point, player Player) {
	state := newStateCopy(g.state)
	if p.X*p.Y < 0 || p.X >= g.Width || p.Y >= g.Height {
		return
	}
	c := state[p.X][p.Y]
	c.Active = false
	state[p.X][p.Y] = c

	g.state = state
}

func (g *Grid) AddPlayer(conn *websocket.Conn) (player Player) {
	player = NewPlayer(conn)
	g.players[player.Name] = player

	go player.ListenMessages()
	go player.ListenWebsocket()

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

func (g *Grid) SetPause(pause bool, player Player) {
	g.Paused = pause
}

func (g *Grid) SetInterval(interval int, player Player) {
	viper.Set("interval", interval)
}

func (g *Grid) Reset(player Player) {
	g.state = newState(g.Width, g.Height)
	g.Generation = 0
}

func (g *Grid) PushStateChange() {
	msg := NewPushStateMessage(g)
	for _, player := range g.players {
		player.messageChan <- msg
	}
}

func (g *Grid) PushPauseChange(changedBy Player) {
	msg := NewPushPauseMessage(g, changedBy)
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

type Cell struct {
	Active bool  `json:"-"`
	Color  Color `json:"c"`
	X      int   `json:"x"`
	Y      int   `json:"y"`
}

func newCell(p Point) Cell {
	return Cell{
		X: p.X,
		Y: p.Y,
	}
}

type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf(`(%v,%v)`, p.X, p.Y)
}

func (p Point) adjacentPoints() (ps []Point) {
	return append(ps,
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
