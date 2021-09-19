package main

import "time"

type message struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	WrittenAt   string    `json:"written_at"`
}
