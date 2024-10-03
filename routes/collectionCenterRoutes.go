package routes

import (
	"proyectoqueso/controllers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupCollectionCenterRoutes(e *echo.Echo, db *gorm.DB) {
	collectionCenterController := controllers.NewCollectionCenterController(db)

	apiCollectionCenterGroup := e.Group("/v1/collectioncenter")

	// Rutas para los centros de acopio
	apiCollectionCenterGroup.GET("", collectionCenterController.GetAllCollectionCenters)       // Obtener todos los centros de acopio
	apiCollectionCenterGroup.POST("", collectionCenterController.CreateCollectionCenter)       // Crear un nuevo centro de acopio
	apiCollectionCenterGroup.GET("/:id", collectionCenterController.GetCollectionCenter)       // Obtener un centro de acopio por ID
	apiCollectionCenterGroup.PUT("/:id", collectionCenterController.UpdateCollectionCenter)    // Actualizar un centro de acopio
	apiCollectionCenterGroup.DELETE("/:id", collectionCenterController.DeleteCollectionCenter) // Eliminar un centro de acopio

	// Rutas para gestionar inventarios del centro de acopio
	apiCollectionCenterGroup.POST("/:id/inventory", collectionCenterController.CreateProductInInventory)                 // Agregar un producto al inventario del centro de acopio
	apiCollectionCenterGroup.GET("/:id/inventory", collectionCenterController.GetInventory)                              // Obtener inventario del centro de acopio
	apiCollectionCenterGroup.GET("/:id/inventory/total/:product_id", collectionCenterController.GetTotalProductQuantity) // Obtener cantidad total de un producto en el inventario del centro de acopio
}
