package routers

import (
	"commandup/config"
	"commandup/models"
	"commandup/utils"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// setAuthCookies sets httpOnly cookies for access and refresh tokens
func setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	c.SetCookie(
		"access_token",
		accessToken,
		15*60, // 15 min
		"/",
		"",
		config.AppConfig.CookieSecure,
		true, // httpOnly
	)

	// Refresh token cookie (7 days)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		7*24*60*60,
		"/",
		"",
		config.AppConfig.CookieSecure,
		true, // httpOnly
	)
}

// clearAuthCookies clears the auth cookies
func clearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", config.AppConfig.CookieSecure, true)
	c.SetCookie("refresh_token", "", -1, "/", "", config.AppConfig.CookieSecure, true)
}

func Login(c *gin.Context) {

	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// add email validation to prevent sql injection
	if !utils.ValidateEmail(loginReq.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	var user models.User
	user, err := models.GetUser(loginReq.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			log.Printf("Database error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if err := utils.CheckPassword(loginReq.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Username)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Set httpOnly cookies
	setAuthCookies(c, accessToken, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

func Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if !utils.ValidateEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	if err := models.CreateUser(user.Email, hashedPassword, user.Username); err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get the newly created user to get their ID
	createdUser, err := models.GetUser(user.Email)
	if err != nil {
		log.Printf("Error fetching created user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
		return
	}

	// Generate tokens and auto-login
	accessToken, err := utils.GenerateAccessToken(createdUser.ID, createdUser.Email, createdUser.Username)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(createdUser.ID)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Set httpOnly cookies
	setAuthCookies(c, accessToken, refreshToken)

	log.Default().Printf("User created for %s", user.Email)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
		"user": gin.H{
			"id":       createdUser.ID,
			"email":    createdUser.Email,
			"username": createdUser.Username,
		},
	})
}

func RefreshToken(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return
	}

	// Validate refresh token
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Get user to get email and username for new access token
	user, err := models.GetUserByID(claims.UserID)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
		return
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Username)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Set new access token cookie
	c.SetCookie(
		"access_token",
		accessToken,
		15*60, // 15 minutes
		"/",
		"",
		config.AppConfig.CookieSecure,
		true, // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed",
	})
}

func Logout(c *gin.Context) {
	clearAuthCookies(c)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
