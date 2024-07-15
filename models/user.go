package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                 uuid.UUID `gorm:"type:char(36);primary_key"`
	Email              string    `gorm:"unique"`
	Password           string
	FirstName          string
	LastName           string
	Birthday           string
	LanguagePreference string
	CurrencyPreference string
	Roles              []Role        `gorm:"many2many:user_roles;"`
	Permissions        []Permission  `gorm:"many2many:user_permissions;"`
	Addresses          []UserAddress // One-to-many relationship
	ConfirmationCode   string
	IsConfirmed        bool `gorm:"default:false"`
}

type UserAddress struct {
	gorm.Model
	UserID        uuid.UUID `gorm:"type:char(36);index"`
	FullName      string
	PhoneNumber   string
	Country       string
	State         string
	City          string
	StreetAddress string
	PostalCode    string
	User          User `gorm:"foreignKey:UserID"`
}

type Role struct {
	gorm.Model
	Name        string       `gorm:"unique"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type UserRole struct {
	UserID uuid.UUID `gorm:"type:char(36)"`
	RoleID int
}

type RolePermission struct {
	RoleID       int
	PermissionID int
}

type UserPermission struct {
	UserID       uuid.UUID `gorm:"type:char(36)"`
	PermissionID int
}
