package main

type User struct {
	Name  string `json:"name"`
	Color *Color `json:"color"`

	gridChan  chan GridUpdate
	eventChan chan bool
}

func NewUser(name string) *User {
	return &User{
		Name:      name,
		Color:     NewRandomColor(),
		gridChan:  make(chan GridUpdate),
		eventChan: make(chan bool),
	}
}
