package database

import (
	"fmt"
	"log"
	"os"

	"github.com/numanijaz/tinyurl/config"
	"github.com/numanijaz/tinyurl/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func getPostgresDataSourceName() string {
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDBName := os.Getenv("POSTGRES_DATABASE")
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		pgHost, pgUser, pgPassword, pgDBName, pgPort,
	)
}

func InitAndMigrateDB() {
	var db *gorm.DB
	var err error
	GO_ENV := config.GetConfig().GO_ENV
	if GO_ENV == "development" {
		db, err = gorm.Open(sqlite.Open("data.db"))
	} else {
		dsn := getPostgresDataSourceName()
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}
	if err != nil {
		log.Fatalf("Failed to connect to db. Error: %v", err)
	}
	log.Println("Connected to", GO_ENV, " database!")

	// global reference
	DB = db

	// migrate
	db.AutoMigrate(&models.UserModel{}, &models.UrlModel{})
}
