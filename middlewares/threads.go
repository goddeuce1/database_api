package middlewares

import (
	"strconv"
	"time"

	"../database"
	"../models"
	ops "../sqlops"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
)

//ThreadCreateMiddleware - create posts for thread
func ThreadCreateMiddleware(posts models.Posts, thread string) (models.Posts, *models.Error) {
	result, errorThread := ThreadFindBySlug(thread)

	if errorThread != nil {
		return nil, models.ErrThreadNotFound
	}

	created := time.Now().Format("2006-01-02 15:04:05")

	for _, index := range posts {
		_, errorUser := UserProfileGetMiddleware(index.Author)

		if errorUser != nil {
			return nil, models.ErrUserNotFound
		}

		if index.Parent != 0 {
			row := database.App.DB.QueryRow(ops.TCMFindPostByParent, index.Parent, result.ID)
			var parent int
			errorParent := row.Scan(&parent)

			if errorParent != nil {
				return nil, models.ErrParentNotFound
			}

		}

	}

	var row *pgx.Row

	if _, error := strconv.Atoi(thread); error == nil {
		row = database.App.DB.QueryRow(ops.TCMFindForumByThread, thread)
	} else {
		row = database.App.DB.QueryRow(ops.TCMFindForumBySlug, thread)
	}

	var forum string
	var id int
	err := row.Scan(&forum, &id)

	if err != nil {
		return posts, nil
	}

	for _, index := range posts {
		index.Forum = forum
		index.Thread = id

		var rows *pgx.Row

		rows = database.App.DB.QueryRow(
			ops.TCMInsertValues,
			index.Author,
			created,
			index.Forum,
			index.Message,
			index.Parent,
			index.Thread,
		)

		var id int
		var timeNow time.Time
		_ = rows.Scan(&id, &timeNow)
		index.ID = id
		index.Created = timeNow

		_, _ = database.App.DB.Exec(
			ops.TCMUpdateForumPostsCount,
			index.Forum,
		)

	}

	return posts, nil

}

//ThreadFindBySlug - find thread by slug (lol)
func ThreadFindBySlug(slug string) (*models.Thread, *models.Error) {
	var row *pgx.Row

	if id, error := strconv.Atoi(slug); error == nil {
		row = database.App.DB.QueryRow(ops.TFByID, id)
	} else {
		row = database.App.DB.QueryRow(ops.TFBySlug, slug)
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
		return nil, models.ErrThreadNotFound
	}

	return &thread, nil
}

//ThreadSlugVoteMiddleware - +-1 vote for thread
func ThreadSlugVoteMiddleware(vote *models.Vote, slug string) (*models.Thread, *models.Error) {
	var err error
	var getPrevVoice int32
	nickname := &pgtype.Varchar{}
	prev := &pgtype.Int4{}
	threadID := &pgtype.Int4{}
	votes := &pgtype.Int4{}

	if id, error := strconv.Atoi(slug); error == nil {
		err = database.App.DB.QueryRow(ops.TSVSelectVoteByID, id, vote.Nickname).Scan(threadID, votes, prev, nickname)
	} else {
		err = database.App.DB.QueryRow(ops.TSVSelectVoteBySlug, slug, vote.Nickname).Scan(threadID, votes, prev, nickname)
	}

	if err != nil || nickname.Status != pgtype.Present || threadID.Status != pgtype.Present {
		return nil, models.ErrThreadNotFound
	}

	if prev.Status == pgtype.Present {
		getPrevVoice = int32(prev.Int)
		_, err = database.App.DB.Exec(ops.TSVUpdateVote, vote.Voice, threadID.Int, nickname.String)
	} else {
		_, err = database.App.DB.Exec(ops.TSVInsertVote, vote.Voice, threadID.Int, nickname.String)
	}

	getNewVotes := votes.Int + (int32(vote.Voice) - getPrevVoice)

	if err != nil {
		return nil, models.ErrThreadNotFound
	}

	getThreadSlug := &pgtype.Varchar{}
	thread := &models.Thread{}

	err = database.App.DB.QueryRow(ops.TSVUpdateVotes, getNewVotes, threadID.Int).Scan(getThreadSlug, &thread.Title, &thread.ID, &thread.Votes, &thread.Author, &thread.Created, &thread.Forum, &thread.Message)

	thread.Slug = getThreadSlug.String

	if err != nil {
		return nil, models.ErrThreadNotFound
	}

	return thread, nil
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
	thread, err := ThreadFindBySlug(slug)

	if err != nil {
		return nil, models.ErrThreadNotFound
	}

	if threadUpdate.Message != "" && threadUpdate.Message != thread.Message {
		var errorExec error

		if id, errorConvert := strconv.Atoi(slug); errorConvert == nil {
			_, errorExec = database.App.DB.Exec(ops.TDPUpdateMessageID, threadUpdate.Message, id)
		} else {
			_, errorExec = database.App.DB.Exec(ops.TDPUpdateMessageSlug, threadUpdate.Message, slug)
		}

		if errorExec != nil {
			return nil, models.ErrGlobal
		}

		thread.Message = threadUpdate.Message
	}

	if threadUpdate.Title != "" && threadUpdate.Title != thread.Title {
		var errorTitle error

		if id, errorConvert := strconv.Atoi(slug); errorConvert == nil {
			_, errorTitle = database.App.DB.Exec(ops.TDPUpdateTitleID, threadUpdate.Title, id)
		} else {
			_, errorTitle = database.App.DB.Exec(ops.TDPUpdateTitleSlug, threadUpdate.Title, slug)
		}

		if errorTitle != nil {
			return nil, models.ErrGlobal
		}

		thread.Title = threadUpdate.Title
	}

	return thread, nil
}

//ThreadPostsMiddleware - returns thread posts
func ThreadPostsMiddleware(slug, limit, since, sort, desc string) (*models.Posts, *models.Error) {
	thread, err := ThreadFindBySlug(slug)

	if err != nil {
		return nil, models.ErrThreadNotFound
	}

	var rows *pgx.Rows
	var error error

	if since != "" {
		if desc == "true" {
			switch string(sort) {
			case "tree":
				rows, error = database.App.DB.Query(ops.TPSinceDescLimitTree, thread.ID, since, limit)
			case "parent_tree":
				rows, error = database.App.DB.Query(ops.TPSinceDescLimitParentTree, thread.ID, since, limit)
			default:
				rows, error = database.App.DB.Query(ops.TPSinceDescLimitFlat, thread.ID, since, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				rows, error = database.App.DB.Query(ops.TPSinceAscLimitTree, thread.ID, since, limit)
			case "parent_tree":
				rows, error = database.App.DB.Query(ops.TPSinceAscLimitParentTree, thread.ID, since, limit)
			default:
				rows, error = database.App.DB.Query(ops.TPSinceAscLimitFlat, thread.ID, since, limit)
			}
		}
	} else {
		if desc == "true" {
			switch string(sort) {
			case "tree":
				rows, error = database.App.DB.Query(ops.TPDescLimitTree, thread.ID, limit)
			case "parent_tree":
				rows, error = database.App.DB.Query(ops.TPDescLimitParentTree, thread.ID, limit)
			default:
				rows, error = database.App.DB.Query(ops.TPDescLimitFlat, thread.ID, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				rows, error = database.App.DB.Query(ops.TPAscLimitTree, thread.ID, limit)
			case "parent_tree":
				rows, error = database.App.DB.Query(ops.TPAscLimitParentTree, thread.ID, limit)
			default:
				rows, error = database.App.DB.Query(ops.TPAscLimitFlat, thread.ID, limit)
			}
		}
	}

	if error != nil {
		return nil, models.ErrGlobal
	}

	defer rows.Close()

	posts := models.Posts{}
	for rows.Next() {
		post := models.Post{}

		if error = rows.Scan(
			&post.ID,
			&post.Author,
			&post.Parent,
			&post.Message,
			&post.Forum,
			&post.Thread,
			&post.Created,
		); err != nil {
			return nil, models.ErrGlobal
		}
		posts = append(posts, &post)
	}

	return &posts, nil
}
