package models

import (
	// "github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	gorm.Model
	UserID        uint
	ProductID     uint
	paymentDate   time.Time
	paymentAmount uint32
	User          User    `gorm:"foreignKey:UserID"`
	Product       Product `gorm:"foreignKey:ProductID"`
}

// type Domain struct {
// 	gorm.Model
// 	UserID uuid.UUID `gorm:"type:char(36)"`
// 	Name   string    `gorm:"unique"`
// 	User   User      `gorm:"foreignKey:UserID"`
// }
