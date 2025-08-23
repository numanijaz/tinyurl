package models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Email          string
	Name           string
	HashedPassword string
	OAuthUser      bool

	ShortURLs []UrlModel `gorm:"foreignKey:UserID;references:ID"`
}
