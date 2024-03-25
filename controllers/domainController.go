package controllers

import (
    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
    "net/http"
    "proyectoqueso/models"
)

type DomainController struct {
	DB *gorm.DB 
}

func NewDomainController(db *gorm.DB) *DomainController {
	return &DomainController{DB: db}
}

func GetDomains(c echo.Context) error {
    return c.JSON(http.StatusOK, "Lista de dominios")
}

func GetDomain(c echo.Context) error {
    id := c.Param("id")
    return c.String(http.StatusOK, id)
}

// Otros métodos para crear, actualizar y eliminar dominios
func (uc *DomainController) CreateDomain(c echo.Context) error {
	userID := uint(1)
	newDomain := models.Domain{
		Name:   "example.com",
		UserID: userID,
		User: models.User{
			ID: 1,
		},
	}

	// newDomain := new(models.Domain)
	if err := c.Bind(newDomain); err != nil {
		return err
	}

	if err := uc.DB.Create(&newDomain).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, newDomain)
}
