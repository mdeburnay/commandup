package users

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID        int
	Email     string
	Password  string
	Username  string
	CreatedAt string
	UpdatedAt string
}

func GetUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var user User

		err := db.QueryRow("SELECT username, email FROM users WHERE id = $1", id).Scan(&user.Username, &user.Email)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON((http.StatusOK), user)
	}
}
