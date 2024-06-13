package routes

import (
	"proyectoqueso/controllers"
	// "proyectoqueso/middleware"

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
    // apiUsersGroup.GET("/:id", userController.GetUser)
    // apiUsersGroup.PUT("/:id", userController.GetUser)
    // apiUsersGroup.DELETE("/:id", userController.GetUser)
}

