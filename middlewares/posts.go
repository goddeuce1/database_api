package middlewares

import (
	"strings"

	"../database"
	"../models"
	ops "../sqlops"
)

//PostIDDetailsPostMiddleware - updates post message by id
func PostIDDetailsPostMiddleware(message string, id string) (*models.Post, *models.Error) {
	row := database.App.DB.QueryRow(ops.PIDUGetPostByID, id)

	post := models.Post{}

	err := row.Scan(
		&post.Author,
		&post.Forum,
		&post.ID,
		&post.IsEdited,
		&post.Message,
		&post.Thread,
		&post.Created,
	)

	if err != nil {
		return nil, models.ErrPostNotFound
	}

	if message != "" && message != post.Message {
		_, err = database.App.DB.Exec(ops.PIDUUpdateMessage, message, true, id)

		if err != nil {
			return nil, models.ErrGlobal
		}

		post.IsEdited = true
		post.Message = message
	}

	return &post, nil
}

//PostIDDetailsGetMiddleware - returns post by id
func PostIDDetailsGetMiddleware(id, related string) (*map[string]interface{}, *models.Error) {
	row := database.App.DB.QueryRow(ops.PIDUGetPostByID, id)

	post := models.Post{}

	err := row.Scan(
		&post.Author,
		&post.Forum,
		&post.ID,
		&post.IsEdited,
		&post.Message,
		&post.Thread,
		&post.Created,
	)

	if err != nil {
		return nil, models.ErrPostNotFound
	}

	result := make(map[string]interface{})

	params := strings.Split(related, ",")

	for _, value := range params {

		if value == "forum" {
			row = database.App.DB.QueryRow(ops.PIDUGetForumByName, post.Forum)
			forum := models.Forum{}
			err = row.Scan(
				&forum.Posts,
				&forum.Slug,
				&forum.Threads,
				&forum.Title,
				&forum.User,
			)

			if err != nil {
				return nil, models.ErrPostNotFound
			}

			result["forum"] = forum
		}

		if value == "user" {
			row = database.App.DB.QueryRow(ops.PIDUGetUserByName, post.Author)
			user := models.User{}
			err = row.Scan(
				&user.About,
				&user.Email,
				&user.Fullname,
				&user.Nickname,
			)

			if err != nil {
				return nil, models.ErrPostNotFound
			}

			result["author"] = user
		}

		if value == "thread" {
			row = database.App.DB.QueryRow(ops.PIDUGetThreadByID, post.Thread)
			thread := models.Thread{}
			err = row.Scan(
				&thread.Author,
				&thread.Created,
				&thread.Forum,
				&thread.ID,
				&thread.Message,
				&thread.Slug,
				&thread.Title,
			)

			if err != nil {
				return nil, models.ErrPostNotFound
			}

			result["thread"] = thread
		}

	}

	result["post"] = post

	return &result, nil
}
