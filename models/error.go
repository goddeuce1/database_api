package models

//Error - struct of error
type Error struct {
	Message string `json:"message,omitempty"`
}

//ErrGlobal -
var ErrGlobal = &Error{Message: "Something went wrong..."}

//ErrForumOwnerNotFound -
var ErrForumOwnerNotFound = &Error{Message: "Forum owner not found"}

//ErrForumAlreadyExists -
var ErrForumAlreadyExists = &Error{Message: "Forum already exists"}

//ErrThreadAlreadyExists -
var ErrThreadAlreadyExists = &Error{Message: "Thread already exists"}

//ErrForumOrAuthorNotFound -
var ErrForumOrAuthorNotFound = &Error{Message: "Forum or author not found"}

//ErrForumNotFound -
var ErrForumNotFound = &Error{Message: "Forum not found"}

//ErrThreadNotFound -
var ErrThreadNotFound = &Error{Message: "Thread not found"}

//ErrUserAlreadyExists -
var ErrUserAlreadyExists = &Error{Message: "User already exists"}

//ErrUserNotFound -
var ErrUserNotFound = &Error{Message: "User not found"}

//ErrSettingsConflict -
var ErrSettingsConflict = &Error{Message: "Can't change user settings (conflict)"}

//ErrParentNotFound -
var ErrParentNotFound = &Error{Message: "Parent not found"}

//ErrPostNotFound -
var ErrPostNotFound = &Error{Message: "Post not found"}
