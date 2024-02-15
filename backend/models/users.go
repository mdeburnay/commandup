package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int
	Email     string
	Password  string
	Username  string
	CreatedAt string
	UpdatedAt string
}

func GetUser(email string) (user User, err error) {
	err = db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password)
	return user, err
}

func CreateUser(email string, password string) (err error) {
	_, err = db.Exec("INSERT INTO users (id, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", uuid.New(), email, password, time.Now(), time.Now())
	return err
}
