package middlewares

import (
	"park_base/park_db/database"
	"park_base/park_db/models"
	ops "park_base/park_db/sqlops"

	"github.com/jackc/pgx"
)

//UserCreateMiddleware - create user
func UserCreateMiddleware(user *models.User) (models.Users, *models.Error) {
	_, err := database.App.DB.Exec(ops.UCMInsertValues, user.About, user.Email, user.Fullname, user.Nickname)

	if err == nil {
		return models.Users{}, nil
	}

	if err.(pgx.PgError).Code == "23505" {
		conflictUsers := models.Users{}
		conflictRows, _ := database.App.DB.Query(ops.UCMGetByNickOrMail, user.Nickname, user.Email)

		defer conflictRows.Close()

		for conflictRows.Next() {
			user := models.User{}
			_ = conflictRows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

			conflictUsers = append(conflictUsers, &user)
		}

		return conflictUsers, models.ErrUserAlreadyExists
	}
	return nil, models.ErrGlobal
}

//UserProfileGetMiddleware - returns desired user
func UserProfileGetMiddleware(nickname string) (*models.User, *models.Error) {
	user := models.User{}
	row := database.App.DB.QueryRow(ops.UCMGetByNickOrMail, nickname, "NULL")

	err := row.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

	if err != nil {
		return nil, models.ErrUserNotFound
	}

	return &user, nil
}

//UserProfilePostMiddleware - returns new user settings
func UserProfilePostMiddleware(user *models.User) (*models.User, *models.Error) {
	err := database.App.DB.QueryRow(ops.UPPUpdateSettings,
		user.About,
		user.Email,
		user.Fullname,
		user.Nickname,
	).Scan(&user.Fullname, &user.About, &user.Email, &user.Nickname)

	if err != nil {

		if err.Error() == "no rows in result set" {
			return nil, models.ErrUserNotFound
		}

		switch err.(pgx.PgError).Code {
		case "23505":
			return nil, models.ErrSettingsConflict
		}

	}

	return user, nil
}
