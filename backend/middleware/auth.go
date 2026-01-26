package middleware

import (
	"commandup/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Try to get token from cookie first (preferred for httpOnly cookies)
		if cookieToken, err := c.Cookie("access_token"); err == nil && cookieToken != "" {
			tokenString = cookieToken
		} else {
			// Fallback to Authorization header for backward compatibility
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		// If no token found, set unauthenticated and continue
		if tokenString == "" {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		// Store user info in context
		c.Set("authenticated", true)
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("username", claims.Username)

		c.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticated, exists := c.Get("authenticated")
		if !exists || !authenticated.(bool) {
			log.Printf("RequireAuth: Blocking unauthenticated request to %s", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		userID, _ := c.Get("user_id")
		log.Printf("RequireAuth: Allowing authenticated request for user_id: %v", userID)
		c.Next()
	}
}

// GetUserID returns the user ID from context if authenticated, otherwise returns 0 and false
func GetUserID(c *gin.Context) (int, bool) {
	authenticated, exists := c.Get("authenticated")
	if !exists || !authenticated.(bool) {
		return 0, false
	}

	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	return userID.(int), true
}

// IsAuthenticated checks if the user is authenticated
func IsAuthenticated(c *gin.Context) bool {
	authenticated, exists := c.Get("authenticated")
	return exists && authenticated.(bool)
}
