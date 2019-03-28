package models

import "time"

//Post - struct of post
type Post struct {
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	ID       int       `json:"id"`
	IsEdited bool      `json:"isEdited"`
	Message  string    `json:"message"`
	Parent   int       `json:"parent"`
	Thread   int       `json:"thread"`
	Path     string    `json:"-"`
}

//PostFull - struct of postfull
type PostFull struct {
	Author *User   `json:"author"`
	Forum  *Forum  `json:"forum"`
	Post   *Post   `json:"post"`
	Thread *Thread `json:"thread"`
}

//PostUpdate - struct of postupdate
type PostUpdate struct {
	Message string `json:"message"`
}

//Posts - slice of post
type Posts []*Post
