package main

import (
	"fmt"
	"time"
)

const (
	EVENT_USER_JOINS EventType = iota
	EVENT_USER_QUITS
	EVENT_AUTO_EVO_ENABLED
	EVENT_AUTO_EVO_DISABLED
)

type EventType int

type Event struct {
	Type  EventType
	Time  time.Time
	Agent *User
}

func (e *Event) String() string {
	switch e.Type {
	case EVENT_USER_JOINS:
		return fmt.Sprintf("%v joined.", e.Agent)
	case EVENT_USER_QUITS:
		return fmt.Sprintf("%v left.", e.Agent)
	case EVENT_AUTO_EVO_ENABLED:
		return "Automatic evolution enabled."
	case EVENT_AUTO_EVO_DISABLED:
		return "Automatic evolution disabled."
	default:
		return "Unrecognized event type"
	}
}
