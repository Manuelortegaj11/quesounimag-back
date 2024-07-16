package routes

import (
	"proyectoqueso/controllers"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupDomainRoutes configura las rutas relacionadas con los dominios
func SetupCategoryRoutes(e *echo.Echo, db *gorm.DB) {

	categoryController := controllers.NewCategoryController(db)
	apiCategoryGroup := e.Group("/v1/category")
	apiCategoryGroup.POST("", categoryController.CreateCategory)
  apiCategoryGroup.PUT(":id", categoryController.UpdateCategory)
  apiCategoryGroup.DELETE(":id", categoryController.DeleteCategory)
	apiCategoryGroup.GET("", categoryController.GetCategories)
	apiCategoryGroup.GET(":id", categoryController.GetCategories)
	apiCategoryGroup.GET(":slug", categoryController.GetCategories)
}
