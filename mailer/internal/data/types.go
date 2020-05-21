package types

import "time"

type Subscriber struct {
	ID    int
	Name  string
	Email string
}

type Thread struct {
	ID         int
	Title      string
	Link       string
	Posts      int
	Votes      int
	Views      int
	DatePosted time.Time
	Seen       bool
}
