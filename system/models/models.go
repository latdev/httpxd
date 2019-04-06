package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type Session struct {
	Session   string `gorm:"type:char(20);unique_index:ix_session"`
	LastWrite time.Time
	Value     string `gorm:"type:LONGTEXT"`
}

type __Session struct {
	Session   string    `gorm:"type:char(20);unique_index:ix_session"`
	LastWrite time.Time `gorm:"unique_index:ix_session"`
	Key       string    `gorm:"type:VARCHAR(100);index"`
	Value     string    `gorm:"type:LONGTEXT"`
}

func ForceModelsCreation(db *gorm.DB) (error) {
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Session{})
	return db.Error
}