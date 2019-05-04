package middlewares

import (
	"park_base/park_db/models"
	ops "park_base/park_db/sqlops"
	"strconv"

	"park_base/park_db/database"

	"github.com/jackc/pgx"
)

//ForumCreateMiddleware - creates forum
func ForumCreateMiddleware(forum *models.Forum) (*models.Forum, *models.Error) {
	err := database.App.DB.QueryRow(ops.FCMInsertValues, forum.Slug, forum.Title, forum.User).Scan(&forum.User, &forum.Threads, &forum.Posts)

	if err != nil {
		switch err.(pgx.PgError).Code {
		case "23502":
			return nil, models.ErrForumOwnerNotFound
		case "23503":
			return nil, models.ErrForumOwnerNotFound
		case "23505":
			dublicateForum, _ := ForumSlugDetailsMiddleware(forum.Slug)
			return dublicateForum, models.ErrForumAlreadyExists
		}
	}

	return forum, nil
}

//ForumSlugCreateMiddleware - create thread
func ForumSlugCreateMiddleware(thread *models.Thread, forum string) (models.Threads, *models.Error) {
	err := database.App.DB.QueryRow(
		ops.FSCMInsertValues,
		thread.Author,
		forum,
		thread.Message,
		thread.Title,
		thread.Slug,
		thread.Created,
	).Scan(&thread.ID, &thread.Forum, &thread.Votes)

	if err == nil {
		_, _ = database.App.DB.Exec(
			ops.TCMUpdateForumThreadsCount,
			thread.Forum,
		)

		threads := models.Threads{}
		threads = append(threads, thread)

		return threads, nil
	}

	if err != nil {

		switch err.(pgx.PgError).Code {
		case "23502":
			return nil, models.ErrForumOrAuthorNotFound
		case "23503":
			return nil, models.ErrForumOrAuthorNotFound
		case "23505":
			conflictThreads := models.Threads{}
			conflictThread := models.Thread{}

			row := database.App.DB.QueryRow(ops.FSCMSelectThreadBySlug, thread.Slug)
			err = row.Scan(
				&conflictThread.Author,
				&conflictThread.Created,
				&conflictThread.Forum,
				&conflictThread.ID,
				&conflictThread.Message,
				&conflictThread.Slug,
				&conflictThread.Title,
				&conflictThread.Votes,
			)

			if err == nil {
				conflictThreads = append(conflictThreads, &conflictThread)
				return conflictThreads, models.ErrThreadAlreadyExists
			}
		}
	}

	return nil, models.ErrGlobal
}

//ForumSlugDetailsMiddleware - returns forum details by slug
func ForumSlugDetailsMiddleware(slug string) (*models.Forum, *models.Error) {
	forum := models.Forum{}
	rows := database.App.DB.QueryRow(ops.FSDGetValues, slug)

	err := rows.Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)

	if err != nil {
		return nil, models.ErrForumNotFound
	}

	return &forum, nil
}

//ForumSlugThreadsMiddleware - returns threads from forum by slug
func ForumSlugThreadsMiddleware(limit, since, desc, slug string) (models.Threads, *models.Error) {
	var rows *pgx.Rows

	if _, error := strconv.Atoi(limit); error == nil {

		if desc == "true" {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query(ops.FSTSelectThreadsLSD, slug, since, limit)
			} else {
				rows, _ = database.App.DB.Query(ops.FSTSelectThreadsLD, slug, limit)
			}

		} else {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query(ops.FSTSelectThreadsLS, slug, since, limit)
			} else {
				rows, _ = database.App.DB.Query(ops.FSTSelectThreadsL, slug, limit)
			}

		}

	} else {

		if desc == "true" {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query(ops.FSTSelectThreadsSD, slug, since)
			} else {
				rows, _ = database.App.DB.Query(ops.FSTSelectThreadsD, slug)
			}

		} else {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query(ops.FSTSelectThreadsS, slug, since)

			} else {
				rows, _ = database.App.DB.Query(ops.FSTSelectThreads, slug)
			}

		}

	}

	defer rows.Close()

	threads := models.Threads{}

	for rows.Next() {
		thread := models.Thread{}
		_ = rows.Scan(
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.ID,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes,
		)
		threads = append(threads, &thread)
	}

	if len(threads) == 0 {
		var slugString string
		_ = database.App.DB.QueryRow(ops.FSCMSelectForumBySlug, slug).Scan(&slugString)

		if slugString == "" {
			return nil, models.ErrForumNotFound
		}

	}

	return threads, nil
}

//ForumSlugUsersMiddleware - returns users
func ForumSlugUsersMiddleware(limit, since, desc, slug string) (models.Users, *models.Error) {
	var rows *pgx.Rows

	if _, error := strconv.Atoi(limit); error == nil {

		if desc == "true" {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query(ops.FSTSelectUsersLSD, since, slug, limit)
			} else {
				rows, _ = database.App.DB.Query(ops.FSTSelectUsersLD, slug, limit)
			}

		} else {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query(ops.FSTSelectUsersLS, since, slug, limit)
			} else {
				rows, _ = database.App.DB.Query(ops.FSTSelectUsersL, slug, limit)
			}

		}

	} else {

		if desc == "true" {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query(ops.FSTSelectUsersSD, since, slug)
			} else {
				rows, _ = database.App.DB.Query(ops.FSTSelectUsersD, slug)
			}

		} else {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query(ops.FSTSelectUsersS, since, slug)

			} else {
				rows, _ = database.App.DB.Query(ops.FSTSelectUsers, slug)
			}

		}

	}

	defer rows.Close()

	users := models.Users{}

	for rows.Next() {
		user := models.User{}
		_ = rows.Scan(
			&user.About,
			&user.Email,
			&user.Fullname,
			&user.Nickname,
		)

		users = append(users, &user)
	}

	if len(users) == 0 {
		var slugString string
		_ = database.App.DB.QueryRow(ops.FSCMSelectForumBySlug, slug).Scan(&slugString)

		if slugString == "" {
			return nil, models.ErrForumNotFound
		}

	}

	return users, nil
}
