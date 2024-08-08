package routes

import (
	"proyectoqueso/controllers"
	// "proyectoqueso/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupDomainRoutes configura las rutas relacionadas con los dominios
func SetupImagesRoutes(e *echo.Echo, db *gorm.DB) {

	imageController := controllers.NewImageController(db)
	apiImagesGroup := e.Group("/v1/images")
	apiImagesGroup.PUT("/product/:id", imageController.UpdateImageProduct)
  apiImagesGroup.PUT("/category/:id", imageController.UpdateImageCategory)
	apiImagesGroup.PUT("/avatar/:id", imageController.UpdateImageAvatar)

  apiImagesGroup.GET("/:type/:id", imageController.GetImage)
}
