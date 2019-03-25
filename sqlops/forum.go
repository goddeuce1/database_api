package middlewares

//FCMInsertValues - used for ForumCreateMiddleware as request text
const FCMInsertValues = `
	INSERT INTO forums("slug", "title", "user") 
	VALUES($1, $2, $3)
	`

//FSDGetValues - used for ForumSlugDetails as request text
const FSDGetValues = `
	SELECT "posts", "slug", "threads", "title", "user"
	FROM forums
	WHERE "slug" = $1
	`

//FSCMSelectForumBySlug - used for ForumSlugCreateMiddleware as request text
const FSCMSelectForumBySlug = `
	SELECT "slug"
	FROM forums
	WHERE lower("slug") = $1
	`

//FSCMSelectUserAndForum - used for ForumSlugCreateMiddleware as request text
const FSCMSelectUserAndForum = `
	SELECT "u.nickname", "f.slug" 
	FROM users AS u, forums AS f
	WHERE lower("u.nickname") = $1 AND lower("f.slug") = $2
	`

//FSCMSelectThreadBySlug - used for ForumCreateMiddleware as request text
const FSCMSelectThreadBySlug = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title"
	FROM threads
	WHERE lower("slug") = $1
	`

//FSCMInsertValues - used for ForumCreateMiddleware as request text
const FSCMInsertValues = `
	INSERT INTO threads("author", "forum", "message", "title", "slug", "created") 
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

//FSTSelectThreadsLSD - used for ForumSlugThreads as request text
const FSTSelectThreadsLSD = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1 AND $2 >= "created"
	ORDER BY "created" DESC
	LIMIT $3
	`

//FSTSelectThreadsLD - used for ForumSlugThreads as request text
const FSTSelectThreadsLD = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1
	ORDER BY "created" DESC
	LIMIT $2
	`

//FSTSelectThreadsLS - used for ForumSlugThreads as request text
const FSTSelectThreadsLS = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1 AND $2 <= "created"
	ORDER BY "created"
	LIMIT $3
	`

//FSTSelectThreadsL - used for ForumSlugThreads as request text
const FSTSelectThreadsL = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1
	ORDER BY "created"
	LIMIT $2
	`

//FSTSelectThreadsSD - used for ForumSlugThreads as request text
const FSTSelectThreadsSD = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1 AND $2 >= "created"
	ORDER BY "created" DESC
	`

//FSTSelectThreadsD - used for ForumSlugThreads as request text
const FSTSelectThreadsD = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1
	ORDER BY created DESC
	`

//FSTSelectThreadsS - used for ForumSlugThreads as request text
const FSTSelectThreadsS = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1 AND $2 <= "created"
	ORDER BY "created"
	`

//FSTSelectThreads - used for ForumSlugThreads as request text
const FSTSelectThreads = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1
	ORDER BY "created"
	`

//TCMUpdateForumThreadsCount - update forum thread count
const TCMUpdateForumThreadsCount = `
	UPDATE forums
	SET "threads" = "threads" + 1
	WHERE "slug" = $1
	`

//FSTSelectUsersL - asc limit
const FSTSelectUsersL = `
	SELECT "about", "email", "fullname", "nickname"
	FROM users
	WHERE "nickname" IN (
		SELECT "author" FROM threads WHERE "forum" = $1
		UNION
		SELECT "author" FROM posts WHERE "forum" = $1
	)
	ORDER BY lower("nickname")::bytea
	LIMIT $2
	`

//FSTSelectUsersLD - desc limit
const FSTSelectUsersLD = `
	SELECT "about", "email", "fullname", "nickname"
	FROM users
	WHERE "nickname" IN (
		SELECT "author" FROM threads WHERE "forum" = $1
		UNION
		SELECT "author" FROM posts WHERE "forum" = $1
	)
	ORDER BY lower("nickname")::bytea DESC
	LIMIT $2
	`

//FSTSelectUsers - asc
const FSTSelectUsers = `
	SELECT "about", "email", "fullname", "nickname"
	FROM users
	WHERE "nickname" IN (
		SELECT "author" FROM threads WHERE "forum" = $1
		UNION
		SELECT "author" FROM posts WHERE "forum" = $1
	)
	ORDER BY lower("nickname")::bytea
	`

//FSTSelectUsersD - desc
const FSTSelectUsersD = `
	SELECT "about", "email", "fullname", "nickname"
	FROM users
	WHERE "nickname" IN (
		SELECT "author" FROM threads WHERE "forum" = $1
		UNION
		SELECT "author" FROM posts WHERE "forum" = $1
	)
	ORDER BY lower("nickname")::bytea DESC
	`

//FSTSelectUsersS - asc since
const FSTSelectUsersS = `
	SELECT "about", "email", "fullname", "nickname"
	FROM users
	WHERE lower("nickname")::bytea > lower($1)::bytea AND "nickname" IN (
		SELECT "author" FROM threads WHERE "forum" = $2
		UNION
		SELECT "author" FROM posts WHERE "forum" = $2
	)
	ORDER BY lower("nickname")::bytea
	`

//FSTSelectUsersSD - desc since
const FSTSelectUsersSD = `
	SELECT "about", "email", "fullname", "nickname"
	FROM users
	WHERE lower("nickname")::bytea < lower($1)::bytea AND "nickname" IN (
		SELECT "author" FROM threads WHERE "forum" = $2
		UNION
		SELECT "author" FROM posts WHERE "forum" = $2
	)
	ORDER BY lower("nickname")::bytea DESC
	`

//FSTSelectUsersLS - asc limit since
const FSTSelectUsersLS = `
	SELECT "about", "email", "fullname", "nickname"
	FROM users
	WHERE lower("nickname")::bytea > lower($1)::bytea AND "nickname" IN (
		SELECT "author" FROM threads WHERE "forum" = $2
		UNION
		SELECT "author" FROM posts WHERE "forum" = $2
	)
	ORDER BY lower("nickname")::bytea
	LIMIT $3
`

//FSTSelectUsersLSD - limit since desc
const FSTSelectUsersLSD = `
	SELECT "about", "email", "fullname", "nickname"
	FROM users
	WHERE lower("nickname")::bytea < lower($1)::bytea AND "nickname" IN (
		SELECT "author" FROM threads WHERE "forum" = $2
		UNION
		SELECT "author" FROM posts WHERE "forum" = $2
	)
	ORDER BY lower("nickname")::bytea DESC
	LIMIT $3
	`
