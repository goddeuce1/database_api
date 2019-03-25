package models

//Forum - struct of forum
type Forum struct {
	Posts   int64  `json:"posts,omitempty"`
	Slug    string `json:"slug"`
	Threads int64  `json:"threads,omitempty"`
	Title   string `json:"title"`
	User    string `json:"user"`
}
