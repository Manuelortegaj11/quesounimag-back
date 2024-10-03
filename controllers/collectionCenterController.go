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

// Constructor
func NewCollectionCenterController(db *gorm.DB) *CollectionCenterController {
	return &CollectionCenterController{
		DB: db,
	}
}

// Crear Centro de Acopio
func (c *CollectionCenterController) CreateCollectionCenter(ctx echo.Context) error {
	var collectionCenter models.CollectionCenter

	// Validar si el JSON está bien formateado
	if err := ctx.Bind(&collectionCenter); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "Error al procesar la solicitud"})
	}

	// Si UserID es nulo, generar un nuevo UUID
	if collectionCenter.UserID == nil {
		newUUID := uuid.New()
		collectionCenter.UserID = &newUUID
	}

	// Intentar crear el nuevo centro de acopio en la base de datos
	if err := c.DB.Create(&collectionCenter).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Error al crear el centro de acopio"})
	}

	return ctx.JSON(http.StatusCreated, collectionCenter)
}

// Obtener todos los Centros de Acopio
func (c *CollectionCenterController) GetAllCollectionCenters(ctx echo.Context) error {
	var collectionCenters []models.CollectionCenter

	// Preload para cargar inventario relacionado
	if err := c.DB.Preload("CollectionCenterInventory").Find(&collectionCenters).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Error al obtener los centros de acopio"})
	}

	return ctx.JSON(http.StatusOK, collectionCenters)
}

// Obtener un Centro de Acopio por ID
func (c *CollectionCenterController) GetCollectionCenter(ctx echo.Context) error {
	id := ctx.Param("id")
	var collectionCenter models.CollectionCenter

	// Preload para cargar la relación con inventario y productos
	if err := c.DB.Preload("CollectionCenterInventory.Product").
		First(&collectionCenter, id).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"message": "Centro de acopio no encontrado"})
	}

	return ctx.JSON(http.StatusOK, collectionCenter)
}

// Actualizar un Centro de Acopio
func (c *CollectionCenterController) UpdateCollectionCenter(ctx echo.Context) error {
	id := ctx.Param("id")
	var collectionCenter models.CollectionCenter

	// Verificar si el centro de acopio existe
	if err := c.DB.First(&collectionCenter, id).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"message": "Centro de acopio no encontrado"})
	}

	// Actualizar los datos del centro de acopio
	if err := ctx.Bind(&collectionCenter); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "Error al procesar la solicitud"})
	}

	// Guardar cambios
	if err := c.DB.Save(&collectionCenter).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Error al actualizar el centro de acopio"})
	}

	return ctx.JSON(http.StatusOK, collectionCenter)
}

// Eliminar un Centro de Acopio
func (c *CollectionCenterController) DeleteCollectionCenter(ctx echo.Context) error {
	id := ctx.Param("id")

	// Eliminar centro de acopio por ID
	if err := c.DB.Delete(&models.CollectionCenter{}, id).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"message": "Centro de acopio no encontrado"})
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// Obtener inventario de un Centro de Acopio
func (c *CollectionCenterController) GetInventory(ctx echo.Context) error {
	id := ctx.Param("id")
	var inventory []models.CollectionCenterInventory

	// Cargar inventario del centro
	if err := c.DB.Where("collection_center_id = ?", id).Preload("Product").Find(&inventory).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Error al obtener el inventario"})
	}

	return ctx.JSON(http.StatusOK, inventory)
}

// Obtener cantidad total de un producto
func (c *CollectionCenterController) GetTotalProductQuantity(ctx echo.Context) error {
	centerID := ctx.Param("id")
	productID := ctx.Param("product_id")

	var totalQuantity int

	// Consultar cantidad total del producto
	if err := c.DB.Model(&models.CollectionCenterInventory{}).
		Where("collection_center_id = ? AND product_id = ?", centerID, productID).
		Select("SUM(quantity)").Scan(&totalQuantity).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Error al obtener la cantidad total del producto"})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"total_quantity": totalQuantity})
}

// Crear producto en inventario del centro de acopio
func (c *CollectionCenterController) CreateProductInInventory(ctx echo.Context) error {
	centerID := ctx.Param("id")
	var inventory models.CollectionCenterInventory

	// Convertir ID de centro a uint
	centerIDUint, err := strconv.ParseUint(centerID, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "ID de centro de acopio no válido"})
	}

	// Convertir ID de producto a uint
	productID := ctx.FormValue("product_id")
	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "ID de producto no válido"})
	}

	// Bind datos del inventario
	if err := ctx.Bind(&inventory); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "Error al procesar la solicitud"})
	}

	// Asignar valores a las claves foráneas
	inventory.CollectionCenterID = uint(centerIDUint)
	inventory.ProductID = uint(productIDUint)

	// Crear inventario en la base de datos
	if err := c.DB.Create(&inventory).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Error al agregar el producto al inventario"})
	}

	return ctx.JSON(http.StatusCreated, inventory)
}
