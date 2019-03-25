package models

//Vote - struct of vote
type Vote struct {
	Nickname string `json:"nickname"`
	Slug     string `json:"slug,omitempty"`
	SlugID   int    `json:"slugid,omitempty"`
	Voice    int    `json:"voice"`
}
