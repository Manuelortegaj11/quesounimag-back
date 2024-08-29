package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectionCenter struct {
	gorm.Model
	ID                       uint                        `gorm:"primaryKey" json:"Id"`
	Name                     string                      `gorm:"not null" json:"Name"`
	Location                 string                      `gorm:"not null" json:"Location"`
	UserID                   *uuid.UUID                 `gorm:"type:char(36);index;null" json:"user_id"`  // Mantener el campo UserID
	User                     *User                       `gorm:"foreignKey:UserID" json:"-"`  // Omite este campo
	CollectionCenterInventory []CollectionCenterInventory `gorm:"foreignKey:CollectionCenterID" json:"collection_center_inventory"`
}
