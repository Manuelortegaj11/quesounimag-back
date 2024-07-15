package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusProcessing OrderStatus = "processing"
	StatusCompleted  OrderStatus = "completed"
)

type Order struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	UserID      uuid.UUID `gorm:"type:char(36);index"`
	User        User      `gorm:"foreignKey:UserID"`
	TotalAmount float64
	Status      OrderStatus `gorm:"type:enum('pending', 'processing', 'completed');default:pending"`
}

type OrderDetail struct {
	gorm.Model
	ID        int
	OrderID   uuid.UUID `gorm:"type:char(36);index"`
	ProductID uint
	Quantity  int
	Order     Order   `gorm:"foreignKey:OrderID"`
	Product   Product `gorm:"foreignKey:ProductID"`
}
