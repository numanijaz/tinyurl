package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

type Config struct {
	GO_ENV               string
	HOST_NAME            string
	PORT                 string
	SECRET_KEY           string
	SESSIONS_SECRET      string
	GITHUB_CLIENT_ID     string
	GITHUB_CLIENT_SECRET string
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	BASE_URL             string
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

var cfg *Config
var CookieStore *sessions.CookieStore
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("Could not load environment from .env file")
		}

		host := env("HOST_NAME", "localhost")
		port := env("PORT", "8000")
		baseURl := env("BASE_URL", fmt.Sprintf("http://%s:%s", host, port))
		secretKey := env("SECRET_KEY", "")

		cfg = &Config{
			GO_ENV:               env("GO_ENV", "development"),
			PORT:                 env("PORT", "8000"),
			HOST_NAME:            host,
			SECRET_KEY:           secretKey,
			SESSIONS_SECRET:      env("SECRET_KEY", secretKey),
			GITHUB_CLIENT_ID:     env("GITHUB_CLIENT_ID", ""),
			GITHUB_CLIENT_SECRET: env("GITHUB_CLIENT_SECRET", ""),
			GOOGLE_CLIENT_ID:     env("GOOGLE_CLIENT_ID", ""),
			GOOGLE_CLIENT_SECRET: env("GOOGLE_CLIENT_SECRET", ""),
			BASE_URL:             baseURl,
		}

		CookieStore = sessions.NewCookieStore([]byte(cfg.SESSIONS_SECRET))
		CookieStore.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   int((24 * time.Hour).Seconds()),
			HttpOnly: true,
			// in productoin
			Secure: true,
		}
	})
	return cfg
}

func SetupOauthProviders() {
	cfg := GetConfig()
	// gothic
	goth.UseProviders(
		github.New(
			cfg.GITHUB_CLIENT_ID,
			cfg.GITHUB_CLIENT_SECRET,
			cfg.BASE_URL+"/api/auth/callback/github",
			"user:email",
		),
		google.New(
			cfg.GOOGLE_CLIENT_ID,
			cfg.GOOGLE_CLIENT_SECRET,
			cfg.BASE_URL+"/api/auth/callback/google",
			"email",
			"profile",
		),
	)
	// gothic.Store = config.CookieStore
}
