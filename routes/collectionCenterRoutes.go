package routes

import (
	"proyectoqueso/controllers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupDomainRoutes configura las rutas relacionadas con los dominios
func SetupCollectionCenterRoutes(e *echo.Echo, db *gorm.DB) {

	collectionCenterController := controllers.NewCollectionCenterController(db)

	apiCollectionCenterGroup := e.Group("/v1/collectioncenter")
	apiCollectionCenterGroup.GET("", collectionCenterController.GetAllCollectionCenter)
	apiCollectionCenterGroup.POST("", collectionCenterController.CreateCollectionCenter)

}
