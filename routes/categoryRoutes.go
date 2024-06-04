package routes

import (
	"proyectoqueso/controllers"
	// "proyectoqueso/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupDomainRoutes configura las rutas relacionadas con los dominios
func SetupCategoryRoutes(e *echo.Echo, db *gorm.DB) {

    categoryController := controllers.NewCategoryController(db)
    apiCategoryGroup := e.Group("/v1/category")
    apiCategoryGroup.POST("/", categoryController.CreateCategory)
}

