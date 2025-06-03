package g

import (
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func SetDB(database *gorm.DB) {
	db = database
}

func DB() *gorm.DB {
	return db
}
