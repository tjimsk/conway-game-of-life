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

func NewUser(conn *websocket.Conn) User {
	return User{
		Name:       NewUserName(),
		Color:      NewRandomColor(),
		conn:       conn,
		evoChan:    make(chan Evolution),
		updateChan: make(chan []*Cell),
		closeChan:  make(chan bool),
	}
}

func NewUserName() (name string) {
	name = fmt.Sprintf(`player%v`, usersCount)
	usersCount++

	return name
}

func GetUser(name string) (u User) {
	mu.Lock()
	u = users[name]
	mu.Unlock()

	return u
}

func RegisterUser(u User) {
	mu.Lock()
	users[u.Name] = u
	mu.Unlock()
}

func UnregisterUser(name string) {
	mu.Lock()
	delete(users, name)
	mu.Unlock()
}
