package controllers

import (
	"net/http"

	"gorm.io/gorm"
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	"proyectoqueso/models"
)

type CollectionCenterController struct {
	DB *gorm.DB
}

func NewCollectionCenterController(db *gorm.DB) *CollectionCenterController {
	return &CollectionCenterController{DB: db}
}

// GetAllCollectionCenter returns all collection centers from the database
func (uc *CollectionCenterController) GetAllCollectionCenter(c echo.Context) error {
	var centers []models.CollectionCenter

	if err := uc.DB.Preload("User").Preload("CollectionCenterInventory").Find(&centers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve collection centers"})
	}

	return c.JSON(http.StatusOK, centers)
}

// CreateCollectionCenter creates a new collection center
func (uc *CollectionCenterController) CreateCollectionCenter(c echo.Context) error {
	var center models.CollectionCenter

	if err := c.Bind(&center); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Si UserID es nil, no se verifica la existencia del usuario
	if center.UserID != nil && *center.UserID != uuid.Nil {
		var user models.User
		if err := uc.DB.First(&user, "id = ?", center.UserID).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "User not found"})
		}
		center.User = &user // Asignar la dirección del valor a User
	} else {
		center.UserID = nil // Asegurar que UserID sea nulo si no fue proporcionado
		center.User = nil    // Asegurar que User sea nulo si no se proporciona UserID
	}

	if err := uc.DB.Create(&center).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create collection center"})
	}

	return c.JSON(http.StatusCreated, center)
}
