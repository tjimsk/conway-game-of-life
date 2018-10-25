package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

var (
	usersCount = 0
)

type User struct {
	Name  string `json:"name"`
	Color *Color `json:"color"`

	conn           *websocket.Conn
	gridUpdateChan chan GridUpdate
	endChan        chan bool
}

func NewUser(conn *websocket.Conn) *User {
	return &User{
		Name:           NewUserName(),
		Color:          NewRandomColor(),
		conn:           conn,
		gridUpdateChan: make(chan GridUpdate),
		endChan:        make(chan bool),
	}
}

func NewUserName() (name string) {
	name = fmt.Sprintf(`player%v`, usersCount)
	usersCount++

	return name
}

func UnregisterUser(name string) {
	mu.Lock()
	delete(users, name)
	mu.Unlock()
}

func RegisterUser(u *User) {
	mu.Lock()
	users[u.Name] = u
	mu.Unlock()
}

func GetUser(name string) (u *User) {
	mu.Lock()
	u = users[name]
	mu.Unlock()

	return u
}

func (u *User) SendUserDetails() error {
	return u.conn.WriteJSON(Message{
		Type:    MESSAGE_TYPE_USER_DETAILS,
		Content: u,
	})
}

func (u *User) SendGridDetails(g *Grid) error {
	return u.conn.WriteJSON(Message{
		Type:    MESSAGE_TYPE_GRID_DETAILS,
		Content: g,
	})
}

func (u *User) SendGridActiveCells(g *Grid) error {
	return u.conn.WriteJSON(Message{
		Type:    MESSAGE_TYPE_GRID_ACTIVE_CELLS,
		Content: g.activeCells(),
	})
}

func (u *User) SendGridUpdate(gu GridUpdate) error {
	return u.conn.WriteJSON(Message{
		Type:    MESSAGE_TYPE_GRID_UPDATE,
		Content: gu,
	})
}
