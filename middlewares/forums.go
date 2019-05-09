package middlewares

import (
	"park_base/park_db/models"
	"strconv"

	"park_base/park_db/database"

	"github.com/jackc/pgx"
)

//ForumCreateMiddleware - creates forum
//PREPARED
func ForumCreateMiddleware(forum *models.Forum) (*models.Forum, *models.Error) {
	err := database.App.DB.QueryRow("FCMSelectNick", forum.User).Scan(&forum.User)

	if err != nil {
		return nil, models.ErrForumOwnerNotFound
	}

	err = database.App.DB.QueryRow("FCMInsertValues", forum.Slug, forum.Title, forum.User).Scan(
		&forum.User,
		&forum.Threads,
		&forum.Posts)

	if err != nil {
		dublicateForum, _ := ForumSlugDetailsMiddleware(forum.Slug)
		return dublicateForum, models.ErrForumAlreadyExists
	}

	return forum, nil
}

//ForumSlugCreateMiddleware - create thread
//PREPARED
func ForumSlugCreateMiddleware(thread *models.Thread, forum string) (models.Threads, *models.Error) {
	err := database.App.DB.QueryRow("FSCMSelectForumBySlug", forum).Scan(&thread.Forum)

	if err != nil {
		return nil, models.ErrForumOrAuthorNotFound
	}

	err = database.App.DB.QueryRow("FCMSelectNick", thread.Author).Scan(&thread.Author)

	if err != nil {
		return nil, models.ErrForumOrAuthorNotFound
	}

	err = database.App.DB.QueryRow(
		"FSCMInsertValues",
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Title,
		thread.Slug,
		thread.Created,
	).Scan(&thread.ID, &thread.Votes)

	if err != nil {
		conflictThreads := models.Threads{}
		conflictThread := models.Thread{}

		row := database.App.DB.QueryRow("TFBySlug", thread.Slug)
		_ = row.Scan(
			&conflictThread.Author,
			&conflictThread.Created,
			&conflictThread.Forum,
			&conflictThread.ID,
			&conflictThread.Message,
			&conflictThread.Slug,
			&conflictThread.Title,
			&conflictThread.Votes,
		)

		conflictThreads = append(conflictThreads, &conflictThread)
		return conflictThreads, models.ErrThreadAlreadyExists
	}

	database.App.DB.Exec("TCMUpdateForumThreadsCount", thread.Forum)
	database.App.DB.Exec("TCMInsertToNewTable", thread.Author, thread.Forum)

	threads := models.Threads{}
	threads = append(threads, thread)

	return threads, nil
}

//ForumSlugDetailsMiddleware - returns forum details by slug
//PREPARED
func ForumSlugDetailsMiddleware(slug string) (*models.Forum, *models.Error) {
	forum := models.Forum{}
	rows := database.App.DB.QueryRow("FSDGetValues", slug)

	err := rows.Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)

	if err != nil {
		return nil, models.ErrForumNotFound
	}

	return &forum, nil
}

//ForumSlugThreadsMiddleware - returns threads from forum by slug
//PREPARED
func ForumSlugThreadsMiddleware(limit, since, desc, slug string) (models.Threads, *models.Error) {
	var forumSlug string

	err := database.App.DB.QueryRow("FSCMSelectForumBySlug", slug).Scan(&forumSlug)

	if err != nil {
		return nil, models.ErrForumNotFound
	}

	var rows *pgx.Rows

	if _, error := strconv.Atoi(limit); error == nil {

		if desc == "true" {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query("FSTSelectThreadsLSD", slug, since, limit)
			} else {
				rows, _ = database.App.DB.Query("FSTSelectThreadsLD", slug, limit)
			}

		} else {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query("FSTSelectThreadsLS", slug, since, limit)
			} else {
				rows, _ = database.App.DB.Query("FSTSelectThreadsL", slug, limit)
			}

		}

	} else {

		if desc == "true" {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query("FSTSelectThreadsSD", slug, since)
			} else {
				rows, _ = database.App.DB.Query("FSTSelectThreadsD", slug)
			}

		} else {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query("FSTSelectThreadsS", slug, since)

			} else {
				rows, _ = database.App.DB.Query("FSTSelectThreads", slug)
			}

		}

	}

	defer rows.Close()

	threads := models.Threads{}

	for rows.Next() {
		thread := models.Thread{}
		thread.Forum = forumSlug

		_ = rows.Scan(
			&thread.Author,
			&thread.Created,
			&thread.ID,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes,
		)

		threads = append(threads, &thread)
	}

	return threads, nil
}

//ForumSlugUsersMiddleware - returns users
//PREPARED
func ForumSlugUsersMiddleware(limit, since, desc, slug string) (models.Users, *models.Error) {
	var slugString string
	err := database.App.DB.QueryRow("FSCMSelectForumBySlug", slug).Scan(&slugString)

	if err != nil {
		return nil, models.ErrForumNotFound
	}

	var rows *pgx.Rows

	if _, error := strconv.Atoi(limit); error == nil {

		if desc == "true" {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query("FSTSelectUsersLSD", since, slug, limit)
			} else {
				rows, _ = database.App.DB.Query("FSTSelectUsersLD", slug, limit)
			}

		} else {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query("FSTSelectUsersLS", since, slug, limit)
			} else {
				rows, _ = database.App.DB.Query("FSTSelectUsersL", slug, limit)
			}

		}

	} else {

		if desc == "true" {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query("FSTSelectUsersSD", since, slug)
			} else {
				rows, _ = database.App.DB.Query("FSTSelectUsersD", slug)
			}

		} else {

			if len(since) != 0 {
				rows, _ = database.App.DB.Query("FSTSelectUsersS", since, slug)

			} else {
				rows, _ = database.App.DB.Query("FSTSelectUsers", slug)
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

	return users, nil
}
