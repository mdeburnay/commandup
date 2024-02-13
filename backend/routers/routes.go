package routers

import (
	routers "commandup/routers/api"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use((gin.Recovery()))

	r.StaticFS("/static", http.Dir(filepath.Join(".", "frontend", "build", "static")))

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	api := r.Group("/api/")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello World",
			})
		})

		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// Auth routes
		api.GET("auth/login", routers.Login(nil))

		// Card routes
		api.GET("cards/upgrades", routers.GetCardUpgrades)
		api.POST("cards/upload-card-collection", routers.UploadCardCollection)

		// User routes
		api.GET("users/:id", routers.GetUser(nil))
	}

	return r
}
