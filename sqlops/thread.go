package middlewares

//SQLInsertVote -
const SQLInsertVote = `
	INSERT INTO votes (thread, nickname, voice) 
	VALUES ($1, $2, $3)
	`

//SQLUpdateVote -
const SQLUpdateVote = `
	UPDATE votes SET 
	voice = $3
	WHERE thread = $1 
	AND nickname = $2
	`

//SQLSelectThreadAndVoteByID -
const SQLSelectThreadAndVoteByID = `
	SELECT votes.voice, threads.id, threads.votes, u.nickname
	FROM (SELECT 1) s
	LEFT JOIN threads ON threads.id = $1
	LEFT JOIN "users" u ON u.nickname = $2
	LEFT JOIN votes ON threads.id = votes.thread AND u.nickname = votes.nickname
	`

//SQLSelectThreadAndVoteBySlug -
const SQLSelectThreadAndVoteBySlug = `
	SELECT votes.voice, threads.id, threads.votes, u.nickname
	FROM (SELECT 1) s
	LEFT JOIN threads ON threads.slug = $1
	LEFT JOIN users as u ON u.nickname = $2
	LEFT JOIN votes ON threads.id = votes.thread AND u.nickname = votes.nickname
	`

//SQLUpdateThreadVotes -
const SQLUpdateThreadVotes = `
	UPDATE threads SET
	votes = $1
	WHERE id = $2
	RETURNING author, created, forum, "message" , slug, title, id, votes
	`

//TCMInsertValues - used for ThreadCreateMiddleware as request text
const TCMInsertValues = `
	INSERT INTO posts("author", "created", "forum", "message", "parent", "thread", "path") 
	VALUES($1, $2, $3, $4, $5, $6, (SELECT path FROM posts WHERE id = $5) || (select currval(pg_get_serial_sequence('posts', 'id'))))
	RETURNING "id", "created"
	`

//TCMFindForumByThread - used for ThreadCreateMiddleware as request text
const TCMFindForumByThread = `
	SELECT "forum", "id"
	FROM threads
	WHERE "id" = $1
	`

//TCMFindForumBySlug - used for ThreadCreateMiddleware as request text
const TCMFindForumBySlug = `
	SELECT "forum", "id"
	FROM threads
	WHERE "slug" = $1
	`

//TFBySlug - thread find by slug
const TFBySlug = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "slug" = $1
	`

//TFByID - thread find by ID
const TFByID = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "id" = $1
	`

//TCMFindPostByParent - finds parent by value
const TCMFindPostByParent = `
	SELECT "id"
	FROM posts
	WHERE "id" = $1 AND "thread" = $2
	`

//TCMUpdateForumPostsCount - updates posts count (forum)
const TCMUpdateForumPostsCount = `
	UPDATE forums
	SET "posts" = "posts" + 1
	WHERE "slug" = $1
	`

//TDPUpdateMessageID - updates thread message
const TDPUpdateMessageID = `
	UPDATE threads
	SET "message" = $1
	WHERE "id" = $2
	`

//TDPUpdateMessageSlug - updates thread message
const TDPUpdateMessageSlug = `
	UPDATE threads
	SET "message" = $1
	WHERE "slug" = $2
	`

//TDPUpdateTitleID - updates thread message
const TDPUpdateTitleID = `
	UPDATE threads
	SET "title" = $1
	WHERE "id" = $2
	`

//TDPUpdateTitleSlug - updates thread message
const TDPUpdateTitleSlug = `
	UPDATE threads
	SET "title" = $1
	WHERE "slug" = $2
	`

//TCMUpdatePath - updates post path
const TCMUpdatePath = `
	UPDATE posts
	SET "path" = concat("path", $1, ".")
	WHERE "id" = $2
	`

//SQLSelectPostsSinceDescLimitTree -
const SQLSelectPostsSinceDescLimitTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 AND (path < (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
	ORDER BY path DESC
	LIMIT $3::TEXT::INTEGER
	`

//SQLSelectPostsSinceDescLimitParentTree -
const SQLSelectPostsSinceDescLimitParentTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE path[1] IN (
		SELECT id
		FROM posts
		WHERE thread = $1 AND parent = 0 AND id < (SELECT path[1] FROM posts WHERE id = $2::TEXT::INTEGER)
		ORDER BY id DESC
		LIMIT $3::TEXT::INTEGER
	)
	ORDER BY path
	`

//SQLSelectPostsSinceDescLimitFlat -
const SQLSelectPostsSinceDescLimitFlat = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 AND id < $2::TEXT::INTEGER
	ORDER BY id DESC
	LIMIT $3::TEXT::INTEGER
	`

//SQLSelectPostsSinceAscLimitTree -
const SQLSelectPostsSinceAscLimitTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 AND (path > (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
	ORDER BY path
	LIMIT $3::TEXT::INTEGER
	`

//SQLSelectPostsSinceAscLimitParentTree -
const SQLSelectPostsSinceAscLimitParentTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE path[1] IN (
		SELECT id
		FROM posts
		WHERE thread = $1 AND parent = 0 AND id > (SELECT path[1] FROM posts WHERE id = $2::TEXT::INTEGER)
		ORDER BY id LIMIT $3::TEXT::INTEGER
	)
	ORDER BY path
	`

//SQLSelectPostsSinceAscLimitFlat -
const SQLSelectPostsSinceAscLimitFlat = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 AND id > $2::TEXT::INTEGER
	ORDER BY id
	LIMIT $3::TEXT::INTEGER
	`

//SQLSelectPostsDescLimitTree -
const SQLSelectPostsDescLimitTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 
	ORDER BY path DESC
	LIMIT $2::TEXT::INTEGER
	`

//SQLSelectPostsDescLimitParentTree -
const SQLSelectPostsDescLimitParentTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 AND path[1] IN (
		SELECT path[1]
		FROM posts
		WHERE thread = $1
		GROUP BY path[1]
		ORDER BY path[1] DESC
		LIMIT $2::TEXT::INTEGER
	)
	ORDER BY path[1] DESC, path
	`

//SQLSelectPostsDescLimitFlat -
const SQLSelectPostsDescLimitFlat = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1
	ORDER BY id DESC
	LIMIT $2::TEXT::INTEGER
	`

//SQLSelectPostsAscLimitTree -
const SQLSelectPostsAscLimitTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 
	ORDER BY path
	LIMIT $2::TEXT::INTEGER
	`

//SQLSelectPostsAscLimitParentTree -
const SQLSelectPostsAscLimitParentTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 AND path[1] IN (
		SELECT path[1] 
		FROM posts 
		WHERE thread = $1 
		GROUP BY path[1]
		ORDER BY path[1]
		LIMIT $2::TEXT::INTEGER
	)
	ORDER BY path
	`

//SQLSelectPostsAscLimitFlat -
const SQLSelectPostsAscLimitFlat = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 
	ORDER BY id
	LIMIT $2::TEXT::INTEGER
	`
