package model

import "time"

type Event struct {
	ID          int
	Title       string
	EventTime   time.Time
	Duration    time.Duration
	Description string
	UserID      int
	// TODO
}
