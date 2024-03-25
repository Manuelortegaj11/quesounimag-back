package models

import (
	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type Domain struct {
	gorm.Model
	UserID uint
	Name   string    `gorm:"unique"`
	User   User      `gorm:"foreignKey:UserID"`
}
