package note

import "time"

type Note struct {
	Id      *int      `json:"id"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Tag     string    `json:"tag"`
}
