package models

import (
	"gorm.io/gorm"
)

type CollectionCenterInventory struct {
	gorm.Model
	CollectionCenterID uint             `gorm:"not null"`
	CollectionCenter   CollectionCenter `gorm:"foreignKey:CollectionCenterID"`
	ProductID          uint             `gorm:"not null"`
	Product            Product          `gorm:"foreignKey:ProductID"`
	Quantity           int              `gorm:"not null"`
}
