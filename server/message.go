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

// const (
// 	MESSAGE_TYPE_USER_DETAILS      MessageType = iota // 0
// 	MESSAGE_TYPE_GRID_DETAILS                         // 1
// 	MESSAGE_TYPE_GRID_ACTIVE_CELLS                    // 2
// 	MESSAGE_TYPE_NEW_EVOLUTION                        // 3
// 	MESSAGE_TYPE_CELLS_UPDATE                         // 4
// )

// type MessageType int

// type Message struct {
// 	Type    MessageType `json:"t"`
// 	Content interface{} `json:"c"`
// }

// type CellsUpdateMessage struct {
// 	Type  MessageType `json:"t"`
// 	Cells []*Cell     `json:"c"`
// }

// func NewUserDetailsMessage(u User) Message {
// 	return Message{
// 		Type:    MESSAGE_TYPE_USER_DETAILS,
// 		Content: u,
// 	}
// }

// func NewGridDetailsMessage(g *Grid) Message {
// 	return Message{
// 		Type:    MESSAGE_TYPE_GRID_DETAILS,
// 		Content: g,
// 	}
// }

// func NewGridActiveCellsMessage(cells []Cell) Message {
// 	return Message{
// 		Type:    MESSAGE_TYPE_GRID_ACTIVE_CELLS,
// 		Content: cells,
// 	}
// }

// func NewEvolutionMessage(evo *Evolution) Message {
// 	return Message{
// 		Type:    MESSAGE_TYPE_NEW_EVOLUTION,
// 		Content: evo,
// 	}
// }

// func NewCellsUpdateMessage(cells []Cell) Message {
// 	return Message{
// 		Type:    MESSAGE_TYPE_CELLS_UPDATE,
// 		Content: cells,
// 	}
// }

// func WriteUserDetails(u User) error {
// 	return u.conn.WriteJSON(NewUserDetailsMessage(u))
// }

// func WriteGridDetails(u User, g *Grid) error {
// 	return u.conn.WriteJSON(NewGridDetailsMessage(g))
// }

// func WriteGridActiveCells(u User, cells []Cell) error {
// 	return u.conn.WriteJSON(NewGridActiveCellsMessage(cells))
// }

// func WriteEvolution(u User, evo *Evolution) error {
// 	return u.conn.WriteJSON(NewEvolutionMessage(evo))
// }

// func WriteCellsUpdate(u User, cells []Cell) error {
// 	return u.conn.WriteJSON(NewCellsUpdateMessage(cells))
// }

// type ActivateMessage struct {
// 	Points []Point
// 	Player User
// }

// type PauseMessage struct {
// 	Pause  bool
// 	Player User
// }
