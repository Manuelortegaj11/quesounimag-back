package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusProcessing OrderStatus = "rejected"
	StatusCompleted  OrderStatus = "accepted"
)

type Order struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:char(36);primary_key"`
	UserID       uuid.UUID `gorm:"type:char(36);index"`
	User         User      `gorm:"foreignKey:UserID"`
	TotalAmount  float64
	Status       OrderStatus   `gorm:"type:enum('pending', 'rejected', 'accepted');default:pending"`
	OrderDetails []OrderDetail `gorm:"foreignKey:OrderID"`
	OrderAddress OrderAddress  `gorm:"foreignKey:OrderID"`
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

type OrderAddress struct {
	gorm.Model
	OrderID       uuid.UUID   `gorm:"type:char(36);index"`
	UserAddressID int         `gorm:"index"`
	UserAddress   UserAddress `gorm:"foreignKey:UserAddressID"`
}
