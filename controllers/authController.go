package controllers

import (
	"net/http"
	"os"
	"time"
	"proyectoqueso/models"
	"proyectoqueso/security"
	util "proyectoqueso/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	// "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

type ErrorResponse struct {
  Message string `json:"message"`
  Errors []string `json:"errors"`
}

var (
  jwtKey = []byte(os.Getenv("JWT_SECRET"))
)

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (au *AuthController) RegisterUser(c echo.Context) error {
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

// Create a route for login an user with email and password with gorm
func (au *AuthController) LoginUser(c echo.Context) error {
  // JSON body
	var requestBody map[string]interface{}

	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if requestBody == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing JSON body")
	}

	if _, ok := requestBody["email"]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing email field")
	}

	if _, ok := requestBody["password"]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing password field")
	}

	email := requestBody["email"].(string)
	password := requestBody["password"].(string)
	// npassword := []byte(password)

	// Retrieve password from database
	var user models.User
	userPassword := user.Password

	if err := au.DB.First(&user, "email = ?", email).Error; err != nil {
      return c.JSON(http.StatusNotFound, "User not found")
	}

 //  if err := bcrypt.CompareHashAndPassword(userPassword, npassword); err != nil {
 //    return c.JSON(http.StatusBadRequest, map[string]string{
 //      "message": "Invalid email or password",
 //    })
 //  }

  security.VerifyPassword(userPassword, password)

  claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
    ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
    Issuer:    user.Email,
  })

  token, err := claims.SignedString(jwtKey)
  if err != nil {
    return c.JSON(http.StatusInternalServerError, map[string]string{
      "message": "Error generating token",
    })
  }

  // Set cookie
  OneWeek := time.Now().Add(time.Hour * 24 * 7)
  cookie := new(http.Cookie)
  cookie.Name = "jwt"
  cookie.Value = token
  cookie.Expires = OneWeek
  cookie.HttpOnly = true
  c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User Login successfully",
	})
}

// LogoutUser returns a empty cookie.
// @Summary Logout the current user.
// @Description Logout the current user.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Logout the current user."
// @Success 401 {string} string "The current user haven't logged-in yet.
// @Router /auth/logout [post]
func LogoutUser(c echo.Context) error {
  cookie := new(http.Cookie)
  cookie.Name = "jwt"
  cookie.Value = ""
  cookie.Expires = time.Now().Add(-time.Hour)
  cookie.HttpOnly = true
  c.SetCookie(cookie)
  return c.String(http.StatusOK, "Logout user!")

}
