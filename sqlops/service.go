package middlewares

//ServiceStatus - used for ServiceStatus as request text
const ServiceStatus = `
	SELECT (
		SELECT COUNT(*) FROM forums
	) AS value1,
	(
		SELECT COUNT(*) FROM posts
	) AS value2,
	(
		SELECT COUNT(*) FROM threads
	) AS value3,
	(
		SELECT COUNT(*) FROM users
	) AS value4;
	`
