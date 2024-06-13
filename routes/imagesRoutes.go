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
    apiImagesGroup.PUT("", imageController.UploadImage)
    apiImagesGroup.PUT("", imageController.UploadImage)
}

