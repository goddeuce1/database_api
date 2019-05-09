package middlewares

import (
	"database/sql"
	"fmt"
	"park_base/park_db/database"
	"park_base/park_db/models"
	ops "park_base/park_db/sqlops"
	"strconv"
	"strings"

	"github.com/jackc/pgx"
)

//ThreadCreateMiddleware - create posts for thread
//PREPARED
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
			row := database.App.DB.QueryRow("TCMFindPostByParent", post.Parent, result.ID)
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

	queryUsers := strings.Builder{}
	queryUsers.WriteString("INSERT INTO fu_table(nickname, forum) VALUES")

	i := 0
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

		queryUsers.WriteString(fmt.Sprintf("('%s', '%s')", post.Author, post.Forum))

		if i < len(posts)-1 {
			queryUsers.WriteString(", ")
		}

		newposts = append(newposts, &post)
		i++

		if post.ID == 1500000 {
			flag = true
		}

	}

	queryUsers.WriteString("ON CONFLICT ON CONSTRAINT fu_table_constraint DO NOTHING")
	database.App.DB.Exec("TCMUpdateForumPostsCount", len(newposts), result.Forum)
	database.App.DB.Exec(queryUsers.String())

	if flag {
		database.App.DB.Exec("VACUUM ANALYZE")
	}

	return newposts, nil
}

//ThreadSlugVoteMiddleware - +-1 vote for thread
//PREPARED
func ThreadSlugVoteMiddleware(vote *models.Vote, slug string) (*models.Thread, *models.Error) {
	thread, err := ThreadDetailsGetMiddleware(slug)
	if err != nil {
		return nil, models.ErrThreadNotFound
	}

	_, err = UserProfileGetMiddleware(vote.Nickname)

	if err != nil {
		return nil, models.ErrUserNotFound
	}

	_ = database.App.DB.QueryRow("TSVoteByID", vote.Voice, vote.Nickname, thread.ID).Scan(&vote.Voice)

	id := strconv.FormatInt(thread.ID, 10)

	return ThreadDetailsGetMiddleware(id)
}

//ThreadDetailsGetMiddleware - get info about thread by slug/id
//PREPARED
func ThreadDetailsGetMiddleware(slug string) (*models.Thread, *models.Error) {
	thread := models.Thread{}
	var row *pgx.Row

	if id, error := strconv.Atoi(slug); error == nil {
		row = database.App.DB.QueryRow("TFByID", id)
	} else {
		row = database.App.DB.QueryRow("TFBySlug", slug)
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
//PREPARED
func ThreadDetailsPostMiddleware(threadUpdate *models.ThreadUpdate, slug string) (*models.Thread, *models.Error) {
	thread, err := ThreadDetailsGetMiddleware(slug)

	if err != nil {
		return nil, models.ErrThreadNotFound
	}

	row := database.App.DB.QueryRow(ops.TDPUpdateMessageID, threadUpdate.Message, threadUpdate.Title, thread.ID)

	_ = row.Scan(
		&thread.Message,
		&thread.Title,
	)

	return thread, nil
}

//ThreadPostsMiddleware - returns thread posts
//PREPARED
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
				rows, _ = database.App.DB.Query("TPSinceDescLimitTree", thread.ID, since, limit)
			case "parent_tree":
				rows, _ = database.App.DB.Query("TPSinceDescLimitParentTree", thread.ID, since, limit)
			default:
				rows, _ = database.App.DB.Query("TPSinceDescLimitFlat", thread.ID, since, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				rows, _ = database.App.DB.Query("TPSinceAscLimitTree", thread.ID, since, limit)
			case "parent_tree":
				rows, _ = database.App.DB.Query("TPSinceAscLimitParentTree", thread.ID, since, limit)
			default:
				rows, _ = database.App.DB.Query("TPSinceAscLimitFlat", thread.ID, since, limit)
			}
		}
	} else {
		if desc == "true" {
			switch string(sort) {
			case "tree":
				rows, _ = database.App.DB.Query("TPDescLimitTree", thread.ID, limit)
			case "parent_tree":
				rows, _ = database.App.DB.Query("TPDescLimitParentTree", thread.ID, limit)
			default:
				rows, _ = database.App.DB.Query("TPDescLimitFlat", thread.ID, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				rows, _ = database.App.DB.Query("TPAscLimitTree", thread.ID, limit)
			case "parent_tree":
				rows, _ = database.App.DB.Query("TPAscLimitParentTree", thread.ID, limit)
			default:
				rows, _ = database.App.DB.Query("TPAscLimitFlat", thread.ID, limit)
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
