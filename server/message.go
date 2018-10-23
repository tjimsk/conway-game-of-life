package main

const (
	OUTBOUND_MESSAGE_TYPE_USER              MessageType = iota // 0
	OUTBOUND_MESSAGE_TYPE_GRID                                 // 1
	OUTBOUND_MESSAGE_TYPE_GRID_ACTIVE_CELLS                    //2
	OUTBOUND_MESSAGE_TYPE_GRID_UPDATE                          //3
	OUTBOUND_MESSAGE_TYPE_EVENT_UPDATE                         //4
)

const (
	INBOUND_MESSAGE_TYPE_SET_NAME MessageType = iota
	INBOUND_MESSAGE_TYPE_SET_COLOR
	INBOUND_MESSAGE_TYPE_SET_AUTO_EVOLUTION
	INBOUND_MESSAGE_TYPE_SET_EVOLUTION_INTERVAL
	INBOUND_MESSAGE_TYPE_POST_COMMENT
)

type MessageType int

type Message struct {
	Type    int         `json:"type"`
	Content interface{} `json:"content"`
}
