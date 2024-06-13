package routes

import (
	"proyectoqueso/controllers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupAuthRoutes configura las rutas relacionadas con los dominios
func SetupAuthRoutes(e *echo.Echo, db *gorm.DB) {

	authController := controllers.NewAuthController(db)
	apiAuthGroup := e.Group("/v1/auth")
	apiAuthGroup.POST("/login", authController.LoginUser)
	apiAuthGroup.POST("/register", authController.RegisterUser)
	apiAuthGroup.POST("/logout", authController.LogoutUser)
	apiAuthGroup.OPTIONS("/session", authController.SessionUser)
	apiAuthGroup.GET("/session", authController.SessionUser)

  // MAIL CONFIRMATION
	apiAuthGroup.POST("/confirm", authController.ConfirmUser)
	apiAuthGroup.POST("/resend-code", authController.ResendConfirmationCode)

}
