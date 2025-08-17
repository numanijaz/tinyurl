package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	SECRET_KEY string
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("Could not load environment from .env file")
		}

		cfg = &Config{
			SECRET_KEY: os.Getenv("SECRET_KEY"),
		}
	})
	return cfg
}
