// controllers/userController.go
package controllers

import (
	"fmt"
	"net/http"
	"proyectoqueso/models"
	"proyectoqueso/security"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AddressController struct {
	DB *gorm.DB
}

func NewAddressController(db *gorm.DB) *AddressController {
	return &AddressController{DB: db}
}

func (uc *AddressController) GetAllAddress(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Authorization header not found",
		})
	}

	// Extract the token from the header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Token not found",
		})
	}

	// Parse the token without validating the signature
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid token",
		})
	}

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
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

func (uc *AddressController) GetAddressByID(c echo.Context) error {
    authHeader := c.Request().Header.Get("Authorization")
    if authHeader == "" {
        return c.JSON(http.StatusUnauthorized, map[string]string{
            "message": "Authorization header not found",
        })
    }

    // Extraer el token del encabezado
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == "" {
        return c.JSON(http.StatusUnauthorized, map[string]string{
            "message": "Token not found",
        })
    }

    // Validar y parsear el token con la firma
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Verificar el método de firma
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        // Devolver la clave secreta utilizada para firmar
        return jwtKey, nil
    })
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "message": "Invalid token",
        })
    }

    // Verificar si el token es válido
    if !token.Valid {
        return c.JSON(http.StatusUnauthorized, map[string]string{
            "message": "Token is not valid",
        })
    }

    // Obtener el ID de la dirección desde el parámetro de la URL
    addressID := c.QueryParam("id")

    // Consultar la dirección por ID directamente
    var address models.UserAddress
    if err := uc.DB.Where("id = ?", addressID).First(&address).Error; err != nil {
        // Manejar el caso en que no se encuentra la dirección
        return c.JSON(http.StatusNotFound, map[string]string{
            "message": "Address not found",
        })
    }

    // Devolver la dirección encontrada
    return c.JSON(http.StatusOK, address)
}

func (uc *AddressController) CreateAddress(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Authorization header not found",
		})
	}

	// Extract the token from the header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Token not found",
		})
	}

	// Parse the token without validating the signature
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid token",
		})
	}

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
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

	requestBody := map[string]interface{}{}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	requiredFields := []string{"addCountry", "addAddress", "addCity", "addName", "addPhone", "addPostalCode", "addState"}
	for _, field := range requiredFields {
		if _, ok := requestBody[field]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Missing %s field", field))
		}
	}

	addCountry := requestBody["addCountry"].(string)
	addAddress := requestBody["addAddress"].(string)
	addCity := requestBody["addCity"].(string)
	addName := requestBody["addName"].(string)
	addPhone := requestBody["addPhone"].(string)
	addPostalCode := requestBody["addPostalCode"].(string)
	addState := requestBody["addState"].(string)

	address := models.UserAddress{
		UserID:        parsedUserID,
		FullName:      addName,
		PhoneNumber:   addPhone,
		Country:       addCountry,
		State:         addState,
		City:          addCity,
		StreetAddress: addAddress,
		PostalCode:    addPostalCode,
	}

	if err := uc.DB.Create(&address).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to add address",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Address successfully added",
	})
}

func (uc *AddressController) DeleteAddress(c echo.Context) error {
	// Extract addressID from URL parameters
	addressID := strings.TrimPrefix(c.Param("id"), "/")
	fmt.Println("Extracting address ID from URL parameters: ", addressID)

	// Parse addressID into uint
	id, err := strconv.ParseUint(addressID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid address ID",
		})
	}

	// Get user ID from token claims
	userID, err := security.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized",
		})
	}

	// Check if the address belongs to the user
	var existingAddress models.UserAddress
	if err := uc.DB.Where("id = ? AND user_id = ?", id, userID).First(&existingAddress).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Address not found or does not belong to the user",
		})
	}

	// Delete the address from the database
	if err := uc.DB.Delete(&existingAddress).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to delete address",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":   "Address deleted",
		"addressId": addressID,
	})
}

func (uc *AddressController) UpdateAddress(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Authorization header not found",
		})
	}

	// Extract the token from the header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Token not found",
		})
	}

	// Parse the token without validating the signature
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid token",
		})
	}

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid token claims",
		})
	}

	userID := claims.Issuer
	// Convert userID to UUID
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid user ID format",
		})
	}

	// Get addressID from URL parameters
	addressID := strings.Trim(c.Param("id"), "/")
	parsedAddressID, err := strconv.ParseUint(addressID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid address ID format",
		})
	}

	requestBody := map[string]interface{}{}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	requiredFields := []string{"addCountry", "addAddress", "addCity", "addName", "addPhone", "addPostalCode", "addState"}
	for _, field := range requiredFields {
		if _, ok := requestBody[field]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Missing %s field", field))
		}
	}

	addCountry := requestBody["addCountry"].(string)
	addAddress := requestBody["addAddress"].(string)
	addCity := requestBody["addCity"].(string)
	addName := requestBody["addName"].(string)
	addPhone := requestBody["addPhone"].(string)
	addPostalCode := requestBody["addPostalCode"].(string)
	addState := requestBody["addState"].(string)

	// Check if the address belongs to the user
	var existingAddress models.UserAddress
	if err := uc.DB.Where("id = ? AND user_id = ?", parsedAddressID, parsedUserID).First(&existingAddress).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Address not found or does not belong to the user",
		})
	}

	existingAddress.FullName = addName
	existingAddress.PhoneNumber = addPhone
	existingAddress.Country = addCountry
	existingAddress.State = addState
	existingAddress.City = addCity
	existingAddress.StreetAddress = addAddress
	existingAddress.PostalCode = addPostalCode

	// Update the address in the database
	if err := uc.DB.Save(&existingAddress).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to update address",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Address updated",
	})
}
