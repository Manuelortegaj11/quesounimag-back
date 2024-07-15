package routes

import (
	"proyectoqueso/controllers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupDomainRoutes configura las rutas relacionadas con los dominios
func SetupOrderRoutes(e *echo.Echo, db *gorm.DB) {

	orderController := controllers.NewOrderController(db)

	apiOrderGroup := e.Group("/v1/order")
	apiOrderGroup.POST("", orderController.CreateOrder)
	apiOrderGroup.GET("", orderController.GetOrdersByUserID)

}
