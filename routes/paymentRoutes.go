package routes

import (
	"proyectoqueso/controllers"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupDomainRoutes configura las rutas relacionadas con los dominios
func SetupPaymentRoutes(e *echo.Echo, db *gorm.DB) {

    paymentController := controllers.NewPaymentController(db)
    apiPaymentGroup := e.Group("/api/v1/payment")

    apiPaymentGroup.GET("", controllers.GetPayment)
    apiPaymentGroup.POST("", paymentController.CreatePayment)
    // apiPaymentGroup.GET("/:id", controllers.GetUser)
    // apiPaymentGroup.PUT("/:id", controllers.GetUser)
    // apiPaymentGroup.DELETE("/:id", controllers.GetUser)
}
