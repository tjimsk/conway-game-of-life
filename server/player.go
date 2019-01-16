package life

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gorilla/websocket"
)

var playerCount = 0

type Player struct {
	Name           string    `json:"name"`
	Color          Color     `json:"color"`
	DisconnectChan chan bool `json:"-"`
	conn           *websocket.Conn
	messageChan    chan Message
}

func NewPlayer(conn *websocket.Conn) Player {
	playerCount++

	return Player{
		Name:           fmt.Sprintf(`player%v`, playerCount),
		Color:          newRandomColor(),
		conn:           conn,
		messageChan:    make(chan Message),
		DisconnectChan: make(chan bool),
	}
}

func (player Player) ListenMessages() {
	for {
		msg := <-player.messageChan
		if err := player.conn.WriteJSON(msg); err != nil {
			log.Println("WriteMessage:", err)
			player.DisconnectChan <- true
			break
		}
	}
}

// keep listening on messages sent from the client websocket connection
// deprecated: client side messages are sent as http requests
func (player Player) ListenWebsocket() {
	for {
		msg := struct{}{}
		if err := player.conn.ReadJSON(&msg); err != nil {
			log.Println("ReadJSON:", err)
			player.DisconnectChan <- true
			break
		}
	}
}

type Color struct {
	R int `json:"R"`
	G int `json:"G"`
	B int `json:"B"`
}

func newColor(R, G, B int) Color {
	return Color{R, G, B}
}

func newRandomColor() Color {
	var (
		R = rand.Intn(255)
		G = rand.Intn(255)
		B = rand.Intn(255)
	)

	return newColor(R, G, B)
}

func averageColor(colors []Color) Color {
	var (
		r     int
		g     int
		b     int
		count = len(colors)
	)

	for _, c := range colors {
		r = r + c.R
		g = g + c.G
		b = b + c.B
	}

	return Color{
		R: r / count,
		G: g / count,
		B: b / count,
	}
}

func (c Color) String() string {
	return fmt.Sprintf(`(%v, %v, %v)`, c.R, c.G, c.B)
}
