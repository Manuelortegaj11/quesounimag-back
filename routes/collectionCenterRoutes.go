package routes

import (
	"proyectoqueso/controllers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupCollectionCenterRoutes(e *echo.Echo, db *gorm.DB) {

	collectionCenterController := controllers.NewCollectionCenterController(db)

	apiCollectionCenterGroup := e.Group("/v1/collectioncenter")
	apiCollectionCenterGroup.GET("", collectionCenterController.GetAllCollectionCenter)
	apiCollectionCenterGroup.POST("", collectionCenterController.CreateCollectionCenter)
  apiCollectionCenterGroup.DELETE(":id", collectionCenterController.DeleteCollectionCenter)
  apiCollectionCenterGroup.PUT(":id", collectionCenterController.UpdateCollectionCenter)

  // Inventarios centro de acopio
	apiCollectionCenterGroup.POST("/inventory", collectionCenterController.CreateProductInInventory)
	apiCollectionCenterGroup.GET("/inventory/total/:product_id", collectionCenterController.GetTotalProductQuantity)

}
