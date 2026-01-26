package config

import (
	"os"
	"strings"
)

type Config struct {
	AllowedOrigins []string
	CookieSecure   bool
}

var AppConfig *Config

func Init() {
	AppConfig = &Config{
		AllowedOrigins: getAllowedOrigins(),
		CookieSecure:   getCookieSecure(),
	}
}

func getAllowedOrigins() []string {
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	if originsEnv != "" {
		origins := strings.Split(originsEnv, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
		return origins
	}
	// Default to localhost:3000 for local development
	return []string{"http://localhost:3000"}
}

func getCookieSecure() bool {
	secureEnv := os.Getenv("COOKIE_SECURE")
	return secureEnv == "true" || secureEnv == "1"
}
