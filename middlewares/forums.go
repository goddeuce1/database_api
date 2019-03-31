package middlewares

import (
	"strconv"
	"strings"

	"../database"
	"../models"
	ops "../sqlops"
	"github.com/jackc/pgx"
)

//ForumCreateMiddleware - creates forum
func ForumCreateMiddleware(forum *models.Forum) (*models.Forum, *models.Error) {
	user, error := UserProfileGetMiddleware(forum.User)

	if error != nil {
		return nil, models.ErrForumOwnerNotFound
	}

	forum.User = user.Nickname

	rows, _ := database.App.DB.Exec(ops.FCMInsertValues, forum.Slug, forum.Title, forum.User)

	if rows.RowsAffected() == 0 {
		dublicateForum, _ := ForumSlugDetailsMiddleware(forum.Slug)
		return dublicateForum, models.ErrForumAlreadyExists
	}

	return forum, nil
}

//ForumSlugCreateMiddleware - create forum slug
func ForumSlugCreateMiddleware(thread *models.Thread, forum string) (models.Threads, *models.Error) {
	_, error := UserProfileGetMiddleware(thread.Author)

	if error != nil {
		return nil, models.ErrForumOrAuthorNotFound
	}

	rowForum := database.App.DB.QueryRow(ops.FSCMSelectForumBySlug, strings.ToLower(forum))

	var gotForum string
	_ = rowForum.Scan(&gotForum)

	if gotForum == "" {
		return nil, models.ErrForumOrAuthorNotFound
	}

	thread.Forum = gotForum

	conflictThreads := models.Threads{}
	conflictThread := models.Thread{}

	row := database.App.DB.QueryRow(ops.FSCMSelectThreadBySlug, strings.ToLower(thread.Slug))
	_ = row.Scan(
		&conflictThread.Author,
		&conflictThread.Created,
		&conflictThread.Forum,
		&conflictThread.ID,
		&conflictThread.Message,
		&conflictThread.Slug,
		&conflictThread.Title,
	)

	if conflictThread.Slug != "" {
		conflictThreads = append(conflictThreads, &conflictThread)
		return conflictThreads, models.ErrThreadAlreadyExists
	}

	var id int
	err := database.App.DB.QueryRow(
		ops.FSCMInsertValues,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Title,
		thread.Slug,
		thread.Created,
	).Scan(&id)

	if err != nil {
		return nil, models.ErrGlobal
	}

	_, err = database.App.DB.Exec(
		ops.TCMUpdateForumThreadsCount,
		thread.Forum,
	)

	if err != nil {
		return nil, models.ErrGlobal
	}

	threads := models.Threads{}
	thread.ID = int64(id)

	threads = append(threads, thread)

	return threads, nil
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
	var err error

	if _, error := strconv.Atoi(limit); error == nil && len(limit) != 0 {

		if desc == "true" {

			if len(since) != 0 {
				rows, err = database.App.DB.Query(ops.FSTSelectThreadsLSD, slug, since, limit)
			} else {
				test, _ := strconv.Atoi(limit)
				rows, err = database.App.DB.Query(ops.FSTSelectThreadsLD, slug, test)
			}

		} else {

			if len(since) != 0 {
				rows, err = database.App.DB.Query(ops.FSTSelectThreadsLS, slug, since, limit)
			} else {
				rows, err = database.App.DB.Query(ops.FSTSelectThreadsL, slug, limit)
			}

		}

	} else {

		if desc == "true" {

			if len(since) != 0 {
				rows, err = database.App.DB.Query(ops.FSTSelectThreadsSD, slug, since)
			} else {
				rows, err = database.App.DB.Query(ops.FSTSelectThreadsD, slug)
			}

		} else {

			if len(since) != 0 {
				rows, err = database.App.DB.Query(ops.FSTSelectThreadsS, slug, since)

			} else {
				rows, err = database.App.DB.Query(ops.FSTSelectThreads, slug)
			}

		}

	}

	if err != nil {
		return nil, models.ErrGlobal
	}

	defer rows.Close()

	threads := models.Threads{}

	for rows.Next() {
		thread := models.Thread{}
		err := rows.Scan(
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
			return nil, models.ErrGlobal
		}

		threads = append(threads, &thread)
	}

	if len(threads) == 0 {
		_, error := ForumSlugDetailsMiddleware(slug)

		if error != nil {
			return nil, models.ErrForumNotFound
		}

	}

	return threads, nil

}

//ForumSlugUsersMiddleware - returns users
func ForumSlugUsersMiddleware(limit, since, desc, slug string) (models.Users, *models.Error) {
	var rows *pgx.Rows
	var err error

	if _, error := strconv.Atoi(limit); error == nil && len(limit) != 0 {

		if desc == "true" {

			if len(since) != 0 {
				rows, err = database.App.DB.Query(ops.FSTSelectUsersLSD, since, slug, limit)
			} else {
				rows, err = database.App.DB.Query(ops.FSTSelectUsersLD, slug, limit)
			}

		} else {

			if len(since) != 0 {
				rows, err = database.App.DB.Query(ops.FSTSelectUsersLS, since, slug, limit)
			} else {
				rows, err = database.App.DB.Query(ops.FSTSelectUsersL, slug, limit)
			}

		}

	} else {

		if desc == "true" {

			if len(since) != 0 {
				rows, err = database.App.DB.Query(ops.FSTSelectUsersSD, since, slug)
			} else {
				rows, err = database.App.DB.Query(ops.FSTSelectUsersD, slug)
			}

		} else {

			if len(since) != 0 {
				rows, err = database.App.DB.Query(ops.FSTSelectUsersS, since, slug)

			} else {
				rows, err = database.App.DB.Query(ops.FSTSelectUsers, slug)
			}

		}

	}

	if err != nil {
		return nil, models.ErrGlobal
	}

	defer rows.Close()

	users := models.Users{}

	for rows.Next() {
		user := models.User{}
		err := rows.Scan(
			&user.About,
			&user.Email,
			&user.Fullname,
			&user.Nickname,
		)

		if err != nil {
			return nil, models.ErrGlobal
		}

		users = append(users, &user)
	}

	if len(users) == 0 {
		_, error := ForumSlugDetailsMiddleware(slug)

		if error != nil {
			return nil, models.ErrForumNotFound
		}

	}

	return users, nil
}
