package routers

import (
	"commandup/config"
	"commandup/middleware"
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
		AllowOrigins:  config.AppConfig.AllowedOrigins,
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Set-Cookie"},
	}))

	api := r.Group("/api/")
	api.Use(middleware.AuthMiddleware())
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
		api.POST("auth/login", routers.Login)
		api.POST("auth/signup", routers.Signup)
		api.POST("auth/refresh", routers.RefreshToken)
		api.POST("auth/logout", routers.Logout)

		// Card routes - upgrades doesn't require auth, but middleware sets context
		api.POST("cards/upgrades", routers.GetCardUpgrades)
		// Upload requires authentication
		api.POST("cards/upload-card-collection", middleware.RequireAuth(), routers.UploadCardCollection)

		// User routes
		api.GET("user/:id", routers.GetUser(nil))
		api.POST("user/create", routers.CreateUser)
	}

	return r
}
