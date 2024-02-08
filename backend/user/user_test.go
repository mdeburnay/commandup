package user

import (
	"database/sql"
)

func TestGetUserByID(db *sql.DB, id int) (User, error) {
	user, err := GetUserByID(db, id)

	if err != nil {
		return user, err
	}

	return user, nil
}
