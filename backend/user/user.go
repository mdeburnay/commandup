package user

import (
	"database/sql"
)

type User struct {
	ID        int
	Email     string
	Password  string
	Username  string
	CreatedAt string
	UpdatedAt string
}

func GetUser(db *sql.DB, id int) (User, error) {
	var user User
	err := db.QueryRow("SELECT username, email FROM users WHERE id = $1", id).Scan(&user.Username, &user.Email)
	if err != nil {
		return user, err
	}

	return user, nil
}
