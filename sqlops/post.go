package middlewares

//PIDUGetPostByID - select post by id
const PIDUGetPostByID = `
	SELECT "author", "forum", "id", "isedited", "message", "thread", "created", "parent"
	FROM posts
	WHERE "id" = $1
	`

//PIDUUpdateMessage - update post message
const PIDUUpdateMessage = `
	UPDATE posts 
	SET 
	"message" = CASE WHEN $1 = '' OR $1 = "message" THEN "message" ELSE $1 END,
	"isedited" = CASE WHEN $1 = '' OR $1 = "message" THEN FALSE ELSE TRUE END
	WHERE "id" = $2
	RETURNING "author", "forum", "id", "isedited", "message", "thread", "created", "parent"
	`

//PIDUGetUserByName - gets user by his nickname
const PIDUGetUserByName = `
	SELECT "about", "email", "fullname", "nickname"
	FROM users
	WHERE "nickname" = $1
	`

//PIDUGetThreadByID - gets thread info by its id
const PIDUGetThreadByID = `
	SELECT "author", "created", "forum", "id", "message", "slug", "title", "votes"
	FROM threads
	WHERE "id" = $1
	`

//PIDUGetForumByName - gets forum info by its name
const PIDUGetForumByName = `
	SELECT "posts", "slug", "threads", "title", "user"
	FROM forums
	WHERE "slug" = $1
	`
