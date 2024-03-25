// controllers/userController.go
package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	// "proyectoqueso/models"
)

// Configura las rutas utilizando los controladores importados
// e.GET("/api/users", controllers.GetUsers)
// e.POST("/api/users", controllers.CreateUser)
// e.GET("/api/users/:id", controllers.GetUser)
// e.PUT("/api/users/:id", controllers.UpdateUser)
// e.DELETE("/api/users/:id", controllers.DeleteUser)
//

type UserController struct {
	DB *gorm.DB // Agrega la instancia de la base de datos como un campo en el controlador
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) CreateUser(c echo.Context) error {
  return nil
}


func (uc *UserController) GetUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, "Get users")
}

// Create a route for get user with gorm
func (uc *UserController)GetUser(c echo.Context) error {
	return c.JSON(http.StatusOK, "Get user")
}

// Create a route for user profile with gorm
func Profile(c echo.Context) error {
	return nil
}

// Create a route for update user profile with gorm
func UpdateProfile(c echo.Context) error {
	return nil
}

// Create a route for delete user profile with gorm
func DeleteProfile(c echo.Context) error {
	return nil
}

// Create a route for user profile with gorm
func GetProfile(c echo.Context) error {
	return nil
}

// Create a route for user profile with gorm
func GetProfiles(c echo.Context) error {
	return nil
}

func ResetPassword(c echo.Context) error {
	return nil
}

// Create a route for change password an user with email and password with gorm
func ChangePassword(c echo.Context) error {
	return nil
}

// Create a route for change email an user with email and password with gorm
func ChangeEmail(c echo.Context) error {
	return nil
}

// Create a route for change email an user with email and password with gorm
func VerifyEmail(c echo.Context) error {
	return nil
}
