package models

import (
	// "github.com/google/uuid"
  "github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                 uuid.UUID `gorm:"type:char(36);primary_key"`
	Email              string `gorm:"unique"`
	Password           string
	FirstName          string
	LastName           string
	Birthday           string
	PhoneNumber        string
	Country            string
	State              string
	City               string
	StreetAddress      string
	PostalCode         string
	LanguagePreference string
	CurrencyPreference string
  Roles []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	gorm.Model
	Name string `gorm:"unique"`
}

