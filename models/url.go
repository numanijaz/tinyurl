package models

import "gorm.io/gorm"

type UrlModel struct {
	gorm.Model
	ShortUrl    string
	OriginalUrl string
	VisitCount  int
}
