package models

type User struct {
	ID        int
	Email     string
	Password  string
	Username  string
	CreatedAt string
	UpdatedAt string
}

func GetUser(email string) (user User, err error) {
	err = db.QueryRow("SELECT id, email, password, username FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password, &user.Username)
	return user, err
}

func CreateUser(email string, password string, username string) (err error) {
	_, err = db.Exec("INSERT INTO users (email, password, username) VALUES ($1, $2, $3)", email, password, username)
	return err
}
