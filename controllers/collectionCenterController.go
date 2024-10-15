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
	// Capturar parámetros de paginación y búsqueda
	pageParam := c.QueryParam("page")
	pageSizeParam := c.QueryParam("pageSize")
	search := c.QueryParam("name")

	// Definir valores predeterminados de paginación
	page := 1
	pageSize := 10

	// Convertir los parámetros de paginación a enteros si existen
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil {
			page = p
		}
	}
	if pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil {
			pageSize = ps
		}
	}

	// Cálculo del offset
	offset := (page - 1) * pageSize

	// Inicializar la lista de centros de colección
	var centers []models.CollectionCenter
	var totalItems int64

	// Construcción de la consulta con filtros de búsqueda
	query := uc.DB.Preload("CollectionCenterInventory").Preload("User")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Contar el número total de registros (sin paginación)
	if err := query.Model(&models.CollectionCenter{}).Count(&totalItems).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to count collection centers"})
	}

	// Ejecutar la consulta con paginación
	if err := query.Offset(offset).Limit(pageSize).Find(&centers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve collection centers"})
	}

	// Crear una estructura para la respuesta personalizada
	type CollectionCenterResponse struct {
		ID                        uint                               `json:"ID"`
		Name                      string                             `json:"Name"`
		Location                  string                             `json:"Location"`
		UserID                    string                             `json:"user_id"` // Se mantiene user_id pero con el nombre del usuario
		CollectionCenterInventory []models.CollectionCenterInventory `json:"collection_center_inventory"`
	}

	// Generar la respuesta con los datos obtenidos
	var response []CollectionCenterResponse
	for _, center := range centers {
		response = append(response, CollectionCenterResponse{
			ID:                        center.ID,
			Name:                      center.Name,
			Location:                  center.Location,
			UserID:                    center.User.FirstName, // Se asigna el nombre del usuario
			CollectionCenterInventory: center.CollectionCenterInventory,
		})
	}

	// Enviar respuesta con paginación y datos
	return c.JSON(http.StatusOK, map[string]interface{}{
		"totalItems": totalItems,
		"data":       response,
	})
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

	return c.JSON(http.StatusOK, center)
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

// Input para el inventario
type InventoryInput struct {
	CollectionCenterID uint `json:"collection_center_id"`
	ProductID          uint `json:"product_id"`
	Quantity           uint `json:"quantity"`
}

// Crear un producto en el inventario de un centro de acopio
func (ctrl *CollectionCenterController) CreateProductInInventory(c echo.Context) error {
	var input InventoryInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Verificar si el producto existe
	var product models.Product
	if err := ctrl.DB.First(&product, input.ProductID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	// Verificar si el centro de acopio existe
	var collectionCenter models.CollectionCenter
	if err := ctrl.DB.First(&collectionCenter, input.CollectionCenterID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Collection center not found"})
	}

	// Crear inventario para el centro de acopio
	inventory := models.CollectionCenterInventory{
		CollectionCenterID: input.CollectionCenterID,
		ProductID:          input.ProductID,
		Quantity:           input.Quantity,
	}

	if err := ctrl.DB.Create(&inventory).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create inventory record"})
	}

	// Preload del producto y del centro de acopio después de crear el registro
	if err := ctrl.DB.Preload("Product").Preload("CollectionCenter").First(&inventory, inventory.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load inventory with relations"})
	}

	return c.JSON(http.StatusOK, inventory)
}

// Obtener la cantidad total de un producto en todos los centros de acopio
func (ctrl *CollectionCenterController) GetTotalProductQuantity(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	var totalQuantity uint
	err = ctrl.DB.Model(&models.CollectionCenterInventory{}).
		Where("product_id = ?", productID).
		Select("SUM(quantity)").
		Scan(&totalQuantity).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to calculate total quantity"})
	}

	return c.JSON(http.StatusOK, map[string]uint{"total_quantity": totalQuantity})
}
