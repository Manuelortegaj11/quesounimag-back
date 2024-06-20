// controllers/userController.go
package controllers

import (
	"net/http"
	"proyectoqueso/models"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

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

// GetUserById retrieves a user by their ID
func (uc *UserController) GetUserById(c echo.Context) error {
	// Obtain the value of the "token" cookie
	cookie, err := c.Cookie("token")
	if err != nil {
		// Handle the case where the cookie is not found or other error occurs
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Token not found",
		})
	}
	tokenValue := cookie.Value

	// Parse the token without validating the signature
	token, err := jwt.ParseWithClaims(tokenValue, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil // Replace "your-secret-key" with your actual secret key
	})
	if err != nil {
		// Handle the case where the token cannot be parsed
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid token",
		})
	}

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		// Handle the case where the claims cannot be extracted
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid token claims",
		})
	}
	userID := claims.Issuer // This should be the user ID

	// Query the user from the database using GORM
	var user models.User
	if err := uc.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		// Handle the case where the user with the given ID is not found
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	user.Password = ""

	// Return the user data
	return c.JSON(http.StatusOK, user)
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

func (uc *UserController) GetAllAddress(c echo.Context) error {
	// Obtener el token de la cookie
	cookie, err := c.Cookie("token")
	if err != nil {
		// Manejar el caso en que no se encuentra la cookie o ocurre otro error
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Token not found",
		})
	}
	tokenValue := cookie.Value

	// Analizar el token sin validar la firma
	token, err := jwt.ParseWithClaims(tokenValue, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil // Reemplaza "your-secret-key" con tu clave secreta real
	})
	if err != nil {
		// Manejar el caso en que no se puede analizar el token
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid token",
		})
	}

	// Extraer el ID del usuario de las reclamaciones del token
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		// Manejar el caso en que no se pueden extraer las reclamaciones
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid token claims",
		})
	}
	userID := claims.Issuer // Esto debería ser el ID del usuario

	// Convertir userID a UUID
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid user ID format",
		})
	}

	// Consultar el usuario en la base de datos usando GORM
	var user models.User
	if err := uc.DB.Where("id = ?", parsedUserID).First(&user).Error; err != nil {
		// Manejar el caso en que no se encuentra el usuario con el ID dado
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	// Consultar las direcciones del usuario
	var addresses []models.UserAddress
	if err := uc.DB.Where("user_id = ?", parsedUserID).Find(&addresses).Error; err != nil {
		// Manejar el caso en que no se encuentran direcciones
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Addresses not found",
		})
	}

	// Devolver las direcciones del usuario
	return c.JSON(http.StatusOK, addresses)
}
