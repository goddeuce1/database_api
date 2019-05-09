package middlewares

//FCMSelectNick - used for ForumCreateMiddleware as request text
const FCMSelectNick = `
	SELECT "nickname"
	FROM users
	WHERE "nickname" = $1
	`

//FCMInsertValues - used for ForumCreateMiddleware as request text
const FCMInsertValues = `
	INSERT INTO forums("slug", "title", "user") 
	VALUES($1, $2, $3)
	RETURNING "user", "threads", "posts"
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
	WHERE "slug" = $1
	`

//FSCMSelectThreadBySlug - used for ForumCreateMiddleware as request text
const FSCMSelectThreadBySlug = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "slug" = $1
	`

//FSCMInsertValues - used for ForumCreateMiddleware as request text
const FSCMInsertValues = `
	INSERT INTO threads("author", "forum", "message", "title", "slug", "created") 
	VALUES($1, $2, $3, $4, nullif($5, ''), $6)
	RETURNING "id", "votes"
	`

//FSTSelectThreadsLSD - used for ForumSlugThreads as request text
const FSTSelectThreadsLSD = `
	SELECT "author", "created", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1 AND $2 >= "created"
	ORDER BY "created" DESC
	LIMIT $3
	`

//FSTSelectThreadsLD - used for ForumSlugThreads as request text
const FSTSelectThreadsLD = `
	SELECT "author", "created", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1
	ORDER BY "created" DESC
	LIMIT $2
	`

//FSTSelectThreadsLS - used for ForumSlugThreads as request text
const FSTSelectThreadsLS = `
	SELECT "author", "created", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1 AND $2 <= "created"
	ORDER BY "created"
	LIMIT $3
	`

//FSTSelectThreadsL - used for ForumSlugThreads as request text
const FSTSelectThreadsL = `
	SELECT "author", "created", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1
	ORDER BY "created"
	LIMIT $2
	`

//FSTSelectThreadsSD - used for ForumSlugThreads as request text
const FSTSelectThreadsSD = `
	SELECT "author", "created", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1 AND $2 >= "created"
	ORDER BY "created" DESC
	`

//FSTSelectThreadsD - used for ForumSlugThreads as request text
const FSTSelectThreadsD = `
	SELECT "author", "created", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1
	ORDER BY created DESC
	`

//FSTSelectThreadsS - used for ForumSlugThreads as request text
const FSTSelectThreadsS = `
	SELECT "author", "created", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "forum" = $1 AND $2 <= "created"
	ORDER BY "created"
	`

//FSTSelectThreads - used for ForumSlugThreads as request text
const FSTSelectThreads = `
	SELECT "author", "created", "id", "message", "slug", "title", "votes"
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

//TCMInsertToNewTable - insert to fu_table
const TCMInsertToNewTable = `
	INSERT INTO fu_table(nickname, forum)
	VALUES($1, $2)
	ON CONFLICT ON CONSTRAINT fu_table_constraint DO 
	NOTHING
	`

//FSTSelectUsersL - asc limit
const FSTSelectUsersL = `
	SELECT a.about, a.email, a.fullname, a.nickname
	FROM fu_table b
	JOIN users a ON a.nickname = b.nickname
	WHERE b.forum = $1
	ORDER BY a.nickname
	LIMIT $2
	`

//FSTSelectUsersLD - desc limit
const FSTSelectUsersLD = `
	SELECT a.about, a.email, a.fullname, a.nickname
	FROM fu_table b
	JOIN users a ON a.nickname = b.nickname
	WHERE b.forum = $1
	ORDER BY a.nickname DESC
	LIMIT $2
	`

//FSTSelectUsers - asc
const FSTSelectUsers = `
	SELECT a.about, a.email, a.fullname, a.nickname
	FROM fu_table b
	JOIN users a ON a.nickname = b.nickname
	WHERE b.forum = $1
	ORDER BY a.nickname
	`

//FSTSelectUsersD - desc
const FSTSelectUsersD = `
	SELECT a.about, a.email, a.fullname, a.nickname
	FROM fu_table b
	JOIN users a ON a.nickname = b.nickname
	WHERE b.forum = $1
	ORDER BY a.nickname DESC
	`

//FSTSelectUsersS - asc since
const FSTSelectUsersS = `
	SELECT a.about, a.email, a.fullname, a.nickname
	FROM fu_table b
	JOIN users a ON a.nickname = b.nickname
	WHERE a.nickname > $1 AND b.forum = $2
	ORDER BY a.nickname
	`

//FSTSelectUsersSD - desc since
const FSTSelectUsersSD = `
	SELECT a.about, a.email, a.fullname, a.nickname
	FROM fu_table b
	JOIN users a ON a.nickname = b.nickname
	WHERE a.nickname < $1 AND b.forum = $2
	ORDER BY a.nickname DESC
	`

//FSTSelectUsersLS - asc limit since
const FSTSelectUsersLS = `
	SELECT a.about, a.email, a.fullname, a.nickname
	FROM fu_table b
	JOIN users a ON a.nickname = b.nickname
	WHERE a.nickname > $1 AND b.forum = $2
	ORDER BY a.nickname
	LIMIT $3
`

//FSTSelectUsersLSD - limit since desc
const FSTSelectUsersLSD = `
	SELECT a.about, a.email, a.fullname, a.nickname
	FROM fu_table b
	JOIN users a ON a.nickname = b.nickname
	WHERE a.nickname < $1 AND b.forum = $2
	ORDER BY a.nickname DESC
	LIMIT $3
	`
