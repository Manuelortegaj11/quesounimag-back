package models

import (
	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectionCenterInventory struct {
    gorm.Model
    CollectionCenterID uint
    CollectionCenter   CollectionCenter
    ProductID          uint
    Product            Product
    Quantity           uint
}
