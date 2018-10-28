package main

import (
	"github.com/gorilla/websocket"
)

const (
	MESSAGE_TYPE_USER_DETAILS      MessageType = iota // 0
	MESSAGE_TYPE_GRID_DETAILS                         // 1
	MESSAGE_TYPE_GRID_ACTIVE_CELLS                    // 2
	MESSAGE_TYPE_NEW_EVOLUTION                        // 3
	MESSAGE_TYPE_CELLS_UPDATE                         // 4
)

type MessageType int

type Message struct {
	Type    MessageType `json:"t"`
	Content interface{} `json:"c"`
}

type CellsUpdateMessage struct {
	Type  MessageType `json:"t"`
	Cells []*Cell     `json:"c"`
}

func NewUserDetailsMessage(u User) Message {
	return Message{
		Type:    MESSAGE_TYPE_USER_DETAILS,
		Content: u,
	}
}

func NewGridDetailsMessage(g *Grid) Message {
	return Message{
		Type:    MESSAGE_TYPE_GRID_DETAILS,
		Content: g,
	}
}

func NewGridActiveCellsMessage(cells []*Cell) Message {
	return Message{
		Type:    MESSAGE_TYPE_GRID_ACTIVE_CELLS,
		Content: cells,
	}
}

func NewEvolutionMessage(evo Evolution) Message {
	return Message{
		Type:    MESSAGE_TYPE_NEW_EVOLUTION,
		Content: evo,
	}
}

func NewCellsUpdateMessage(cells []*Cell) Message {
	return Message{
		Type:    MESSAGE_TYPE_CELLS_UPDATE,
		Content: cells,
	}
}

func WriteUserDetails(conn *websocket.Conn, u User) error {
	return conn.WriteJSON(NewUserDetailsMessage(u))
}

func WriteGridDetails(conn *websocket.Conn, g *Grid) error {
	return conn.WriteJSON(NewGridDetailsMessage(g))
}

func WriteGridActiveCells(conn *websocket.Conn, cells []*Cell) error {
	return conn.WriteJSON(NewGridActiveCellsMessage(cells))
}

func WriteEvolution(conn *websocket.Conn, evo Evolution) error {
	return conn.WriteJSON(NewEvolutionMessage(evo))
}

func WriteCellsUpdate(conn *websocket.Conn, cells []*Cell) error {
	return conn.WriteJSON(NewCellsUpdateMessage(cells))
}

func ReadCellsUpdate(conn *websocket.Conn) (msg CellsUpdateMessage, err error) {
	if err := conn.ReadJSON(&msg); err != nil {
		return msg, err
	}

	return msg, nil
}
