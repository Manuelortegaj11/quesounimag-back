package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectionCenter struct {
	gorm.Model
	ID        uint        `gorm:"primaryKey"`
	Name      string      `gorm:"not null"`
	Location  string      `gorm:"not null"`
	UserID    *uuid.UUID  `gorm:"type:char(36);index;null"`  // Permitir nulos
	User      *User       `gorm:"foreignKey:UserID"`
	CollectionCenterInventory []CollectionCenterInventory `gorm:"foreignKey:CollectionCenterID"`
}
