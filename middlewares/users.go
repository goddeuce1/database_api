package middlewares

import (
	"strings"

	"github.com/jackc/pgx"

	"../database"
	"../models"
	ops "../sqlops"
)

//UserCreateMiddleware - create user
func UserCreateMiddleware(user *models.User) (models.Users, *models.Error) {
	rows, _ := database.App.DB.Exec(ops.UCMInsertValues, user.About, user.Email, user.Fullname, user.Nickname)

	if rows.RowsAffected() == 0 {
		conflictUsers := models.Users{}
		conflictRows, err := database.App.DB.Query(ops.UCMGetByNickOrMail, user.Nickname, user.Email)

		if err != nil {
			return nil, models.ErrGlobal
		}

		defer conflictRows.Close()

		for conflictRows.Next() {
			user := models.User{}
			err := conflictRows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

			if err != nil {
				return nil, models.ErrGlobal
			}

			conflictUsers = append(conflictUsers, &user)
		}

		return conflictUsers, models.ErrUserAlreadyExists
	}

	return models.Users{}, nil
}

//UserProfileGetMiddleware - returns desired user
func UserProfileGetMiddleware(nickname string) (*models.User, *models.Error) {
	user := models.User{}
	row := database.App.DB.QueryRow(ops.UCMGetByNickOrMail, strings.ToLower(nickname), "NULL")

	err := row.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

	if err != nil {
		return nil, models.ErrUserNotFound
	}

	return &user, nil
}

//UserProfilePostMiddleware - returns new user settings
func UserProfilePostMiddleware(user *models.User) (*models.User, *models.Error) {
	rows := database.App.DB.QueryRow(ops.UPPUpdateSettings,
		user.About,
		user.Email,
		user.Fullname,
		user.Nickname,
	).Scan(&user.Fullname, &user.About, &user.Email, &user.Nickname)

	if rows != nil {
		_, ok := rows.(pgx.PgError)

		if ok {
			return nil, models.ErrSettingsConflict
		}
	}

	profile, error := UserProfileGetMiddleware(user.Nickname)

	if error != nil {
		return nil, models.ErrUserNotFound
	}

	return profile, nil
}
