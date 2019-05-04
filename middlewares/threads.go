package middlewares

import (
	"database/sql"
	"fmt"
	"log"
	"park_base/park_db/database"
	"park_base/park_db/models"
	ops "park_base/park_db/sqlops"
	"strconv"
	"strings"

	"github.com/jackc/pgx"
)

//ThreadCreateMiddleware - create posts for thread
func ThreadCreateMiddleware(posts models.Posts, thread string) (models.Posts, *models.Error) {
	result, errorThread := ThreadDetailsGetMiddleware(thread)

	if errorThread != nil {
		return nil, models.ErrThreadNotFound
	}

	if len(posts) == 0 {
		return models.Posts{}, nil
	}

	query := strings.Builder{}

	query.WriteString("INSERT INTO posts(author, message, parent, forum, thread) VALUES")

	for index, post := range posts {

		if post.Parent != 0 {
			row := database.App.DB.QueryRow(ops.TCMFindPostByParent, post.Parent, result.ID)
			parent := 0
			errorParent := row.Scan(&parent)

			if errorParent != nil {
				return nil, models.ErrParentNotFound
			}

		}

		_, getError := UserProfileGetMiddleware(post.Author)

		if getError != nil {
			return nil, models.ErrUserNotFound
		}

		if post.Parent != 0 {
			query.WriteString(fmt.Sprintf("('%s', '%s', '%d', '%s', '%d')", post.Author, post.Message, post.Parent, result.Forum, result.ID))
		} else {
			query.WriteString(fmt.Sprintf("('%s', '%s', NULL, '%s', '%d')", post.Author, post.Message, result.Forum, result.ID))
		}

		if index < len(posts)-1 {
			query.WriteString(", ")
		}

	}

	query.WriteString(" RETURNING id, thread, forum, created, isedited, author, message, parent")

	rows, _ := database.App.DB.Query(query.String())
	defer rows.Close()

	newposts := models.Posts{}
	flag := false

	for rows.Next() {
		post := models.Post{}
		var value sql.NullInt64
		_ = rows.Scan(
			&post.ID,
			&post.Thread,
			&post.Forum,
			&post.Created,
			&post.IsEdited,
			&post.Author,
			&post.Message,
			&value,
		)

		if value.Valid {
			post.Parent = value.Int64
		}

		newposts = append(newposts, &post)

		if post.ID == 1500000 {
			log.Println("vacuum will be")
			flag = true
		}

	}

	database.App.DB.Exec(ops.TCMUpdateForumPostsCount, len(newposts), result.Forum)

	if flag {
		log.Println("vacuum now!")
		database.App.DB.Exec("VACUUM ANALYZE")
	}

	return newposts, nil
}

//ThreadSlugVoteMiddleware - +-1 vote for thread
func ThreadSlugVoteMiddleware(vote *models.Vote, slug string) (*models.Thread, *models.Error) {
	var err error

	if id, error := strconv.Atoi(slug); error == nil {
		err = database.App.DB.QueryRow(ops.TSVoteByID, vote.Voice, vote.Nickname, id).Scan(&vote.Voice)
	} else {
		err = database.App.DB.QueryRow(ops.TSVoteBySlug, vote.Voice, vote.Nickname, slug).Scan(&vote.Voice)
	}

	if err != nil {
		switch err.(pgx.PgError).Code {
		case "23503":
			return nil, models.ErrThreadNotFound
		}
	}

	return ThreadDetailsGetMiddleware(slug)
}

//ThreadDetailsGetMiddleware - get info about thread by slug/id
func ThreadDetailsGetMiddleware(slug string) (*models.Thread, *models.Error) {
	thread := models.Thread{}
	var row *pgx.Row

	if id, error := strconv.Atoi(slug); error == nil {
		row = database.App.DB.QueryRow(ops.TFByID, id)
	} else {
		row = database.App.DB.QueryRow(ops.TFBySlug, slug)
	}

	err := row.Scan(
		&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.ID,
		&thread.Message,
		&thread.Slug,
		&thread.Title,
		&thread.Votes,
	)

	if err != nil {
		return nil, models.ErrThreadNotFound
	}

	return &thread, nil
}

//ThreadDetailsPostMiddleware - updates thread info
func ThreadDetailsPostMiddleware(threadUpdate *models.ThreadUpdate, slug string) (*models.Thread, *models.Error) {
	var row *pgx.Row

	if id, errorConvert := strconv.Atoi(slug); errorConvert == nil {
		row = database.App.DB.QueryRow(ops.TDPUpdateMessageID, threadUpdate.Message, threadUpdate.Title, id)
	} else {
		row = database.App.DB.QueryRow(ops.TDPUpdateMessageSlug, threadUpdate.Message, threadUpdate.Title, slug)
	}

	thread := models.Thread{}

	err := row.Scan(
		&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.ID,
		&thread.Message,
		&thread.Slug,
		&thread.Title,
		&thread.Votes,
	)

	if err != nil {

		if err.Error() == "no rows in result set" {
			return nil, models.ErrThreadNotFound
		}

	}

	return &thread, nil
}

//ThreadPostsMiddleware - returns thread posts
func ThreadPostsMiddleware(slug, limit, since, sort, desc string) (*models.Posts, *models.Error) {
	thread, err := ThreadDetailsGetMiddleware(slug)

	if err != nil {
		return nil, models.ErrThreadNotFound
	}

	var rows *pgx.Rows

	if since != "" {
		if desc == "true" {
			switch string(sort) {
			case "tree":
				rows, _ = database.App.DB.Query(ops.TPSinceDescLimitTree, thread.ID, since, limit)
			case "parent_tree":
				rows, _ = database.App.DB.Query(ops.TPSinceDescLimitParentTree, thread.ID, since, limit)
			default:
				rows, _ = database.App.DB.Query(ops.TPSinceDescLimitFlat, thread.ID, since, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				rows, _ = database.App.DB.Query(ops.TPSinceAscLimitTree, thread.ID, since, limit)
			case "parent_tree":
				rows, _ = database.App.DB.Query(ops.TPSinceAscLimitParentTree, thread.ID, since, limit)
			default:
				rows, _ = database.App.DB.Query(ops.TPSinceAscLimitFlat, thread.ID, since, limit)
			}
		}
	} else {
		if desc == "true" {
			switch string(sort) {
			case "tree":
				rows, _ = database.App.DB.Query(ops.TPDescLimitTree, thread.ID, limit)
			case "parent_tree":
				rows, _ = database.App.DB.Query(ops.TPDescLimitParentTree, thread.ID, limit)
			default:
				rows, _ = database.App.DB.Query(ops.TPDescLimitFlat, thread.ID, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				rows, _ = database.App.DB.Query(ops.TPAscLimitTree, thread.ID, limit)
			case "parent_tree":
				rows, _ = database.App.DB.Query(ops.TPAscLimitParentTree, thread.ID, limit)
			default:
				rows, _ = database.App.DB.Query(ops.TPAscLimitFlat, thread.ID, limit)
			}
		}
	}

	defer rows.Close()

	posts := models.Posts{}
	for rows.Next() {
		post := models.Post{}
		var value sql.NullInt64

		_ = rows.Scan(
			&post.ID,
			&post.Author,
			&value,
			&post.Message,
			&post.Forum,
			&post.Thread,
			&post.Created,
			&post.IsEdited,
		)

		if value.Valid {
			post.Parent = value.Int64
		}

		posts = append(posts, &post)
	}

	return &posts, nil
}
