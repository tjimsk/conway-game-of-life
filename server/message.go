package main

const (
	MESSAGE_TYPE_USER_DETAILS      MessageType = iota // 0
	MESSAGE_TYPE_GRID_DETAILS                         // 1
	MESSAGE_TYPE_GRID_ACTIVE_CELLS                    // 2
	MESSAGE_TYPE_GRID_UPDATE                          // 3
	MESSAGE_TYPE_UPDATE_CELLS                         // 4
)

type MessageType int

type Message struct {
	Type    MessageType `json:"t"`
	Content interface{} `json:"c"`
}

type ActivateCellsMessage struct {
	Type  MessageType `json:"t"`
	Cells []Cell      `json:"c"`
	User  User        `json:"u"`
}
