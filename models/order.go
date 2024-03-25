package models

import (
	// "github.com/google/uuid"
	"gorm.io/gorm"
  "time"
)

type Order struct {
	gorm.Model
	UserID       uint
	ProductID    uint
	dateStarted  time.Time
	dateFinished time.Time
	User         User    `gorm:"foreignKey:UserID"`
	Product      Product `gorm:"foreignKey:ProductID"`
}
