package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

var (
	usersCount = 0
)

type User struct {
	Name  string `json:"n"`
	Color Color  `json:"c"`

	conn       *websocket.Conn
	evoChan    chan Evolution
	updateChan chan []*Cell
	closeChan  chan bool
}

func NewUserName() (name string) {
	name = fmt.Sprintf(`player%v`, usersCount)
	usersCount++

	return name
}
