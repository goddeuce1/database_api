package middlewares

import (
	"bytes"
	"fmt"
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
		fmt.Println(err)
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
		index.ID = int64(id)
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

	prevVoice := &pgtype.Int4{}
	threadID := &pgtype.Int4{}
	threadVotes := &pgtype.Int4{}
	userNickname := &pgtype.Varchar{}

	if id, error := strconv.Atoi(slug); error == nil {
		err = database.App.DB.QueryRow(ops.SQLSelectThreadAndVoteByID, id, vote.Nickname).Scan(prevVoice, threadID, threadVotes, userNickname)
	} else {
		err = database.App.DB.QueryRow(ops.SQLSelectThreadAndVoteBySlug, slug, vote.Nickname).Scan(prevVoice, threadID, threadVotes, userNickname)
	}

	if err != nil {
		return nil, models.ErrThreadNotFound
	}

	if threadID.Status != pgtype.Present || userNickname.Status != pgtype.Present {
		return nil, models.ErrThreadNotFound
	}

	var prevVoiceInt int32

	if prevVoice.Status == pgtype.Present {
		prevVoiceInt = int32(prevVoice.Int)
		_, err = database.App.DB.Exec(ops.SQLUpdateVote, threadID.Int, userNickname.String, vote.Voice)
	} else {
		_, err = database.App.DB.Exec(ops.SQLInsertVote, threadID.Int, userNickname.String, vote.Voice)
	}

	newVotes := threadVotes.Int + (int32(vote.Voice) - prevVoiceInt)

	if err != nil {
		return nil, models.ErrThreadNotFound
	}

	thread := &models.Thread{}
	slugNullable := &pgtype.Varchar{}
	err = database.App.DB.QueryRow(ops.SQLUpdateThreadVotes, newVotes, threadID.Int).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Message, slugNullable, &thread.Title, &thread.ID, &thread.Votes)
	thread.Slug = slugNullable.String

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
func ThreadPostsMiddleware(slugOrID string, limit, since, sort, desc []byte) (*models.Posts, *models.Error) {
	thread, err := ThreadFindBySlug(slugOrID)

	if err != nil {
		return nil, models.ErrThreadNotFound
	}

	var queryRows *pgx.Rows
	var error error

	if since != nil {
		if bytes.Equal([]byte("true"), desc) {
			switch string(sort) {
			case "tree":
				queryRows, error = database.App.DB.Query(ops.SQLSelectPostsSinceDescLimitTree, thread.ID, since, limit)
			case "parent_tree":
				queryRows, error = database.App.DB.Query(ops.SQLSelectPostsSinceDescLimitParentTree, thread.ID, since, limit)
			default:
				queryRows, error = database.App.DB.Query(ops.SQLSelectPostsSinceDescLimitFlat, thread.ID, since, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				queryRows, error = database.App.DB.Query(ops.SQLSelectPostsSinceAscLimitTree, thread.ID, since, limit)
			case "parent_tree":
				queryRows, error = database.App.DB.Query(ops.SQLSelectPostsSinceAscLimitParentTree, thread.ID, since, limit)
			default:
				queryRows, error = database.App.DB.Query(ops.SQLSelectPostsSinceAscLimitFlat, thread.ID, since, limit)
			}
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			switch string(sort) {
			case "tree":
				queryRows, error = database.App.DB.Query(ops.SQLSelectPostsDescLimitTree, thread.ID, limit)
			case "parent_tree":
				queryRows, error = database.App.DB.Query(ops.SQLSelectPostsDescLimitParentTree, thread.ID, limit)
			default:
				queryRows, error = database.App.DB.Query(ops.SQLSelectPostsDescLimitFlat, thread.ID, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				queryRows, error = database.App.DB.Query(ops.SQLSelectPostsAscLimitTree, thread.ID, limit)
			case "parent_tree":
				queryRows, error = database.App.DB.Query(ops.SQLSelectPostsAscLimitParentTree, thread.ID, limit)
			default:
				queryRows, error = database.App.DB.Query(ops.SQLSelectPostsAscLimitFlat, thread.ID, limit)
			}
		}
	}
	defer queryRows.Close()

	if error != nil {
		return nil, models.ErrGlobal
	}

	posts := models.Posts{}
	for queryRows.Next() {
		post := models.Post{}

		if error = queryRows.Scan(
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
