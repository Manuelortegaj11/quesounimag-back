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
	apiUsersGroup.GET("", userController.GetAllUsers)
	apiUsersGroup.GET("/me", userController.GetUserById)
	apiUsersGroup.POST("", userController.CreateUser)

}
