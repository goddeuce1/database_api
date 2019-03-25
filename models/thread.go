package models

import "time"

//Thread - struct of thread
type Thread struct {
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Forum   string    `json:"forum,omitempty"`
	ID      int32     `json:"id,omitempty"`
	Message string    `json:"message"`
	Slug    string    `json:"slug,omitempty"`
	Title   string    `json:"title"`
	Votes   int       `json:"votes,omitempty"`
}

//ThreadUpdate - struct of threadupdate
type ThreadUpdate struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

//Threads - slice of thread
type Threads []*Thread