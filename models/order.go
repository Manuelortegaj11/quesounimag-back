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
	UserID      uint
	User        User `gorm:"foreignKey:UserID"`
	TotalAmount float64
	Status      OrderStatus `gorm:"type:enum('pending', 'processing', 'completed');default:pending"`
}

type OrdersDetail struct {
	gorm.Model
	ID        int
	Quantity  int
	OrderID   uint
	ProductID uint
	Order     Order   `gorm:"foreignKey:OrderID"`
	Product   Product `gorm:"foreignKey:ProductID"`
}
