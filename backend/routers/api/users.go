package routers

import (
	"database/sql"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var user models.User

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
