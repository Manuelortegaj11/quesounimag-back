package controllers

import (
	"net/http"
	"proyectoqueso/models"
	"proyectoqueso/security"
	util "proyectoqueso/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	// "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CategoryController struct {
	DB *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{DB: db}
}

func (au *CategoryController) CreateCategory(c echo.Context) error {
  // c.Logger().Info(config.GetEnv("JWT_SECRET"))
  // JSON body
	var requestBody map[string]interface{}

	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if requestBody == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing JSON body")
	}

	if _, ok := requestBody["name"]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing name field")
	}

	if _, ok := requestBody["email"]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing email field")
	}

	if _, ok := requestBody["password"]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing password field")
	}

	name := requestBody["name"].(string)
	email := requestBody["email"].(string)
	password := requestBody["password"].(string)

  match, err := util.IsValidEmail(email)
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "Invalid email")
  }

	if match {
		// Check if email already exists
		var user models.User
		queryUser := au.DB.Where("email = ?", email).First(&user)

		if queryUser.Error == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Email already exists")
		}

    hashedPassword, _ := security.EncryptPassword(password)
    verifyPassword := security.VerifyPassword(hashedPassword, password); 

    if verifyPassword != nil {
      return echo.NewHTTPError(http.StatusBadRequest, "Password can't verify");
    }

		// c.Logger().Info(hashedPassword)
		// c.Logger().Info(verifyPassword)

		// Create user with data of the Request
		newUser := &models.User{
      ID: uuid.New(),
			FirstName: name,
			Email:     email,
			Password:  string(hashedPassword),
		}

		if err := au.DB.Create(newUser).Error; err != nil {
			return err
		}

    return c.JSON(http.StatusOK, map[string]string{
      "message": "User created successfully",
    })

	} 
  return nil
}

