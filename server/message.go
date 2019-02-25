package life

type Message interface{}

type PushStateMessage struct {
	ActiveCells []Cell `json:"activeCells"`
	Generation  int    `json:"generation"`
	Interval    int    `json:"interval"`
}

type PushIntervalMessage struct {
	Interval int    `json:"interval"`
	Player   Player `json:"player"`
}

type PushPlayerMessage struct {
	Player Player `json:"player"`
}

type RequestActivateMessage struct {
	Points []Point `json:"points"`
	Player Player  `json:"player"`
}

type RequestDeactivateMessage struct {
	Point  Point  `json:"point"`
	Player Player `json:"player"`
}

type RequestIntervalMessage struct {
	Interval int    `json:"interval"`
	Player   Player `json:"player"`
}

type RequestResetMessage struct {
	Player Player `json:"player"`
}

func NewPushStateMessage(g *Grid) PushStateMessage {
	cells := []Cell{}
	gridCells := g.cellsCopy()

	for p, _ := range g.liveCellsCopy() {
		cells = append(cells, gridCells[p])
	}
	return PushStateMessage{
		ActiveCells: cells,
		Generation:  g.Generation,
		Interval:    g.Interval,
	}
}

func NewPushIntervalMessage(g *Grid, player Player) PushIntervalMessage {
	return PushIntervalMessage{
		Interval: g.Interval,
		Player:   player,
	}
}

func NewPushPlayerMessage(player Player) PushPlayerMessage {
	return PushPlayerMessage{
		Player: player,
	}
}
