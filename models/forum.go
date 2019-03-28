package models

//Forum - struct of forum
type Forum struct {
	Posts   int    `json:"posts,omitempty"`
	Slug    string `json:"slug"`
	Threads int    `json:"threads,omitempty"`
	Title   string `json:"title"`
	User    string `json:"user"`
}
