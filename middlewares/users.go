package middlewares

import (
	"park_base/park_db/database"
	"park_base/park_db/models"
)

//UserCreateMiddleware - create user
//PREPARED
func UserCreateMiddleware(user *models.User) (models.Users, *models.Error) {
	rows, _ := database.App.DB.Exec("UCMInsertValues", user.About, user.Email, user.Fullname, user.Nickname)

	if rows.RowsAffected() == 0 {
		conflictUsers := models.Users{}
		conflictRows, _ := database.App.DB.Query("UCMGetByNickOrMail", user.Nickname, user.Email)

		defer conflictRows.Close()

		for conflictRows.Next() {
			user := models.User{}
			_ = conflictRows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

			conflictUsers = append(conflictUsers, &user)
		}

		return conflictUsers, models.ErrUserAlreadyExists
	}

	return models.Users{}, nil
}

//UserProfileGetMiddleware - returns desired user
//PREPARED
func UserProfileGetMiddleware(nickname string) (*models.User, *models.Error) {
	user := models.User{}
	row := database.App.DB.QueryRow("UCMGetByNick", nickname)

	err := row.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

	if err != nil {
		return nil, models.ErrUserNotFound
	}

	return &user, nil
}

//UserProfilePostMiddleware - returns new user settings
//PREPARED
func UserProfilePostMiddleware(user *models.User) (*models.User, *models.Error) {
	profile, error := UserProfileGetMiddleware(user.Nickname)

	if error != nil {
		return nil, models.ErrUserNotFound
	}

	user.Nickname = profile.Nickname

	err := database.App.DB.QueryRow("UPPUpdateSettings", user.About, user.Email, user.Fullname, profile.Nickname).Scan(
		&user.Fullname,
		&user.About,
		&user.Email,
	)

	if err != nil {
		return nil, models.ErrSettingsConflict
	}

	return user, nil
}
