package routes

import (
	"proyectoqueso/controllers"
	// "proyectoqueso/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupDomainRoutes configura las rutas relacionadas con los dominios
func SetupProductRoutes(e *echo.Echo, db *gorm.DB) {

    productController := controllers.NewProductController(db)
    apiProductGroup := e.Group("/v1/product")
    apiProductGroup.POST("", productController.CreateProduct)
    apiProductGroup.PUT("", productController.UpdateProduct)
    apiProductGroup.DELETE("", productController.DeleteProduct)
    apiProductGroup.GET("", productController.GetProducts)
    apiProductGroup.GET(":id", productController.GetProducts)
    apiProductGroup.GET(":slug", productController.GetProducts)
}

