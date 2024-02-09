package handlers

import (
	"database/sql"
	"main/auth"
	"main/cards"
	"main/users"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, conn *sql.DB) {
	UserRoutes(r, conn)
	AuthRoutes(r, conn)
	CardsRoutes(r, conn)
}

func UserRoutes(r *gin.Engine, conn *sql.DB) {
	r.GET("users/:id", users.GetUser(conn))
}

func AuthRoutes(r *gin.Engine, conn *sql.DB) {
	r.POST("auth/login", auth.Login("email@example.com", "examplepassword"))
}

func CardsRoutes(r *gin.Engine, conn *sql.DB) {
	r.GET("cards/upgrades", cards.GetCards(conn))
	r.POST("cards/upload-card-collection", cards.UploadCardCollection(conn))
}
