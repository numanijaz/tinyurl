package models

import "gorm.io/gorm"

type UrlModel struct {
	gorm.Model
	UniqueHash  string
	OriginalUrl string
	VisitCount  int
	UserID      uint
	User        UserModel `gorm:"foreignKey:UserID;references:ID"`
}
