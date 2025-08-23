package models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Email          string
	Name           string
	HashedPassword string

	ShortURLs []UrlModel `gorm:"foreignKey:UserID;references:ID"`
}
