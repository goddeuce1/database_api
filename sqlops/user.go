package middlewares

//UCMGetByNickOrMail - used for UserCreateMiddleware as request text
const UCMGetByNickOrMail = `
	SELECT "about", "email", "fullname", "nickname"
	FROM users
	WHERE "nickname" = $1 OR "email" = $2
	`

//UCMGetByNick - used for UserProfileGetMiddleware as request text
const UCMGetByNick = `
	SELECT "about", "email", "fullname", "nickname"
	FROM users
	WHERE "nickname" = $1
	`

//UCMInsertValues - used for UserCreateMiddleware as request text
const UCMInsertValues = `
	INSERT INTO users("about", "email", "fullname", "nickname") 
	VALUES($1, $2, $3, $4)
	`

//UPPUpdateSettings - used for UserProfilePostMiddleware as request text
const UPPUpdateSettings = `
	UPDATE users
	SET 
	"about" = coalesce(nullif($1, ''), "about"), 
	"email" = coalesce(nullif($2, ''), "email"), 
	"fullname" = coalesce(nullif($3, ''), "fullname")
	WHERE nickname = $4
	RETURNING "fullname", "about", "email"
	`
