package life

import (
	"github.com/spf13/viper"
)

type Message interface{}

type PushStateMessage struct {
	ActiveCells []Cell `json:"activeCells"`
	Height      int    `json:"height"`
	Width       int    `json:"width"`
	Generation  int    `json:"generation"`
	Paused      bool   `json:"paused"`
	Interval    int    `json:"interval"`
}

func NewPushStateMessage(g *Grid) PushStateMessage {
	cells := []Cell{}
	for _, row := range g.state {
		for _, c := range row {
			if c.Active {
				cells = append(cells, c)
			}
		}
	}
	return PushStateMessage{
		ActiveCells: cells,
		Height:      g.Height,
		Width:       g.Width,
		Generation:  g.Generation,
		Paused:      g.Paused,
		Interval:    viper.GetInt("interval"),
	}
}

type PushPauseMessage struct {
	Paused bool   `json:"paused"`
	Player Player `json:"player"`
}

func NewPushPauseMessage(g *Grid, player Player) PushPauseMessage {
	return PushPauseMessage{
		Paused: g.Paused,
		Player: player,
	}
}

type PushIntervalMessage struct {
	Interval int    `json:"interval"`
	Player   Player `json:"player"`
}

func NewPushIntervalMessage(g *Grid, player Player) PushIntervalMessage {
	return PushIntervalMessage{
		Interval: viper.GetInt("interval"),
		Player:   player,
	}
}

type PushPlayerMessage struct {
	Player Player `json:"player"`
}

func NewPushPlayerMessage(player Player) PushPlayerMessage {
	return PushPlayerMessage{
		Player: player,
	}
}

type RequestActivateMessage struct {
	Points []Point `json:"points"`
	Player Player  `json:"player"`
}

type RequestDeactivateMessage struct {
	Point  Point  `json:"point"`
	Player Player `json:"player"`
}

type RequestPauseMessage struct {
	Pause  bool   `json:"pause"`
	Player Player `json:"player"`
}

type RequestIntervalMessage struct {
	Interval int    `json:"interval"`
	Player   Player `json:"player"`
}

type RequestResetMessage struct {
	Player Player `json:"player"`
}
