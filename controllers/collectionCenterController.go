package controllers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"proyectoqueso/models"
	"strconv"
)

// CollectionCenterController
type CollectionCenterController struct {
	DB *gorm.DB
}

// NewCollectionCenterController es el constructor para CollectionCenterController
func NewCollectionCenterController(db *gorm.DB) *CollectionCenterController {
	return &CollectionCenterController{
		DB: db,
	}
}

// Crear Centro de Acopio
func (c *CollectionCenterController) CreateCollectionCenter(ctx echo.Context) error {
	var collectionCenter models.CollectionCenter
	if err := ctx.Bind(&collectionCenter); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "Error al procesar la solicitud"})
	}

	// Generar un nuevo UUID para el UserID si no existe
	if collectionCenter.UserID == nil {
		collectionCenter.UserID = new(uuid.UUID)
		*collectionCenter.UserID = uuid.New()
	}

	if err := c.DB.Create(&collectionCenter).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Error al crear el centro de acopio"})
	}

	return ctx.JSON(http.StatusCreated, collectionCenter)
}

// Obtener todos los Centros de Acopio
func (c *CollectionCenterController) GetAllCollectionCenters(ctx echo.Context) error {
	var collectionCenters []models.CollectionCenter
	if err := c.DB.Preload("CollectionCenterInventory").Find(&collectionCenters).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Error al obtener los centros de acopio"})
	}

	return ctx.JSON(http.StatusOK, collectionCenters)
}

// Obtener un Centro de Acopio por ID
func (c *CollectionCenterController) GetCollectionCenter(ctx echo.Context) error {
	id := ctx.Param("id")
	var collectionCenter models.CollectionCenter

	// Preload para cargar la relación con el inventario
	if err := c.DB.Preload("CollectionCenterInventory").Preload("CollectionCenterInventory.Product").
		First(&collectionCenter, id).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"message": "Centro de acopio no encontrado"})
	}

	return ctx.JSON(http.StatusOK, collectionCenter)
}

// Actualizar un Centro de Acopio
func (c *CollectionCenterController) UpdateCollectionCenter(ctx echo.Context) error {
	id := ctx.Param("id")
	var collectionCenter models.CollectionCenter

	if err := c.DB.First(&collectionCenter, id).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"message": "Centro de acopio no encontrado"})
	}

	if err := ctx.Bind(&collectionCenter); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "Error al procesar la solicitud"})
	}

	if err := c.DB.Save(&collectionCenter).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Error al actualizar el centro de acopio"})
	}

	return ctx.JSON(http.StatusOK, collectionCenter)
}

// Eliminar un Centro de Acopio
func (c *CollectionCenterController) DeleteCollectionCenter(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.CollectionCenter{}, id).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"message": "Centro de acopio no encontrado"})
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// Obtener inventario de un Centro de Acopio
func (c *CollectionCenterController) GetInventory(ctx echo.Context) error {
	id := ctx.Param("id")
	var inventory []models.CollectionCenterInventory

	// Cargar productos en el inventario
	if err := c.DB.Where("collection_center_id = ?", id).Preload("Product").Find(&inventory).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Error al obtener el inventario"})
	}

	return ctx.JSON(http.StatusOK, inventory)
}

// Obtener cantidad total de un producto en el inventario del centro de acopio
func (c *CollectionCenterController) GetTotalProductQuantity(ctx echo.Context) error {
	centerID := ctx.Param("id")
	productID := ctx.Param("product_id")

	var totalQuantity int

	// Realizar la consulta para obtener la cantidad total del producto en el inventario
	if err := c.DB.Model(&models.CollectionCenterInventory{}).
		Where("collection_center_id = ? AND product_id = ?", centerID, productID).
		Select("SUM(quantity)").Scan(&totalQuantity).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Error al obtener la cantidad total del producto"})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"total_quantity": totalQuantity})
}

// Crear un producto en el inventario del centro de acopio
func (c *CollectionCenterController) CreateProductInInventory(ctx echo.Context) error {
	centerID := ctx.Param("id")
	var inventory models.CollectionCenterInventory

	// Convertir centerID a uint
	centerIDUint, err := strconv.ParseUint(centerID, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "ID de centro de acopio no válido"})
	}

	// Verificar si el producto existe antes de agregarlo al inventario
	productID := ctx.FormValue("product_id")
	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "ID de producto no válido"})
	}

	// Bind de los datos del inventario
	if err := ctx.Bind(&inventory); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "Error al procesar la solicitud"})
	}

	// Asignar los valores convertidos
	inventory.CollectionCenterID = uint(centerIDUint)
	inventory.ProductID = uint(productIDUint)

	// Crear el inventario
	if err := c.DB.Create(&inventory).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Error al agregar el producto al inventario"})
	}

	return ctx.JSON(http.StatusCreated, inventory)
}
