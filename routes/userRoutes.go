package routes

import (
	"proyectoqueso/controllers"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupDomainRoutes configura las rutas relacionadas con los dominios
func SetupUserRoutes(e *echo.Echo, db *gorm.DB) {

	userController := controllers.NewUserController(db)

	apiUsersGroup := e.Group("/v1/user")
	// apiUsersGroup.Use(middleware.JwtMiddleware)
	apiUsersGroup.GET("/me", userController.GetUserById)
	apiUsersGroup.POST("", userController.CreateUser)

}

func SetupAddressRoutes(e *echo.Echo, db *gorm.DB) {

	userController := controllers.NewUserController(db)
	apiAddressGroup := e.Group("/v1/address")
	apiAddressGroup.GET("", userController.GetAllAddress)
	apiAddressGroup.POST("", userController.CreateAddress)
  apiAddressGroup.PUT(":id", userController.UpdateAddress)
	apiAddressGroup.DELETE(":id", userController.DeleteAddress)
}
