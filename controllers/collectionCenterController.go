package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"proyectoqueso/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

	if err := uc.DB.Preload("CollectionCenterInventory").Find(&centers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve collection centers"})
	}

	// Crear una estructura para la respuesta personalizada
	type CollectionCenterResponse struct {
		ID                        uint                               `json:"ID"`
		Name                      string                             `json:"Name"`
		Location                  string                             `json:"Location"`
		UserID                    *uuid.UUID                         `json:"user_id"`
		CollectionCenterInventory []models.CollectionCenterInventory `json:"collection_center_inventory"`
	}

	var response []CollectionCenterResponse
	for _, center := range centers {
		response = append(response, CollectionCenterResponse{
			ID:                        center.ID,
			Name:                      center.Name,
			Location:                  center.Location,
			UserID:                    center.UserID,
			CollectionCenterInventory: center.CollectionCenterInventory,
		})
	}

	return c.JSON(http.StatusOK, response)
}

// CreateCollectionCenter creates a new collection center
func (uc *CollectionCenterController) CreateCollectionCenter(c echo.Context) error {
	var center models.CollectionCenter

	if err := c.Bind(&center); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Asegurarse de que UserID sea nulo si no fue proporcionado
	if center.UserID != nil && *center.UserID == uuid.Nil {
		center.UserID = nil
	}

	if err := uc.DB.Create(&center).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create collection center"})
	}

	return c.JSON(http.StatusCreated, center)
}

func (uc *CollectionCenterController) DeleteCollectionCenter(c echo.Context) error {
	id := c.Param("id")

	// Remover cualquier slash al inicio del ID
	id = strings.TrimPrefix(id, "/")

	var center models.CollectionCenter

	// Buscar el centro de acopio por su ID
	if err := uc.DB.First(&center, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Collection center not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to find collection center"})
	}

	// Eliminar el centro de acopio
	if err := uc.DB.Delete(&center).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete collection center"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Collection center deleted successfully"})
}

func (uc *CollectionCenterController) UpdateCollectionCenter(c echo.Context) error {
    id := c.Param("id")

    // Eliminar cualquier barra inclinada del inicio y del final
    id = strings.Trim(id, "/")

    // Comprobar si id no está vacío antes de convertir
    if id == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
    }

    // Convertir el ID a un entero
    collectionCenterID, err := strconv.Atoi(id)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
    }

    var center models.CollectionCenter

    // Buscar el centro de acopio por su ID
    if err := uc.DB.First(&center, collectionCenterID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.JSON(http.StatusNotFound, map[string]string{"error": "Collection center not found"})
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to find collection center"})
    }

    // Crear una estructura para los datos del centro de acopio
    var updatedCenter models.CollectionCenter

    // Parsear el cuerpo de la solicitud directamente en la estructura
    if err := c.Bind(&updatedCenter); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
    }

    // Actualizar solo los campos permitidos en el modelo
    center.Name = updatedCenter.Name
    center.Location = updatedCenter.Location
    // Agrega otros campos según sea necesario

    // Guardar los cambios
    if err := uc.DB.Save(&center).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update collection center"})
    }

    // Retornar el centro de acopio actualizado
    return c.JSON(http.StatusOK, center)
}
