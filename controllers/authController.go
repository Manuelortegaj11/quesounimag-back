package controllers

import (
	"fmt"
	"net/http"
	"os"
	"proyectoqueso/models"
	"proyectoqueso/security"
	"proyectoqueso/utils"
	util "proyectoqueso/utils"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	// "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

type ErrorResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

var (
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
)

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func GenerateCode() {

}

func (au *AuthController) RegisterUser(c echo.Context) error {

	requestBody := map[string]interface{}{}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	requiredFields := []string{"birthday", "country", "state", "city", "email", "password", "firstName", "lastName", "phoneNumber", "streetAddress", "postalCode"}
	for _, field := range requiredFields {
		if _, ok := requestBody[field]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Missing %s field", field))
		}
	}

	email := requestBody["email"].(string)
	password := requestBody["password"].(string)
	firstName := requestBody["firstName"].(string)
	lastName := requestBody["lastName"].(string)
	phoneNumber := requestBody["phoneNumber"].(string)
	birthday := requestBody["birthday"].(string)
	streetAddress := requestBody["streetAddress"].(string)
	postalCode := requestBody["postalCode"].(string)
	country := requestBody["country"].(string)
	state := requestBody["state"].(string)
	city := requestBody["city"].(string)

	if _, err := time.Parse("2006-01-02", birthday); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid date format. Use yyyy-mm-dd.")
	}

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
		verifyPassword := security.VerifyPassword(hashedPassword, password)

		if verifyPassword != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Password can't verify")
		}

		confirmationCode := utils.GenerateConfirmationCode()
		subject := "Confirm your account"
		body := fmt.Sprintf("Your confirmation code is: %s", confirmationCode)

		if err := util.SendEmail(email, subject, body); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send confirmation email")
		}

		// Create user with data of the Request
		newUserID := uuid.New()
		newUser := &models.User{
			ID:               newUserID,
			Email:            email,
			Password:         string(hashedPassword),
			FirstName:        firstName,
			LastName:         lastName,
			Birthday:         birthday,
			ConfirmationCode: confirmationCode,
			IsConfirmed:      false,
			Addresses: []models.UserAddress{
				{
					UserID:        newUserID,
					PhoneNumber:   phoneNumber,
					Country:       country,
					State:         state,
					City:          city,
					StreetAddress: streetAddress,
					PostalCode:    postalCode,
				},
			},
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

	match, err := util.IsValidEmail(email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email")
	}

	var user models.User

	if err := au.DB.First(&user, "email = ?", email).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	verifyPassword := security.VerifyPassword(user.Password, password)
	if match {
		if verifyPassword != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Password can't verify")
		}
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		Issuer:    user.ID.String(),
	})

	token, err := claims.SignedString(jwtKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error generating token",
		})
	}

	// // Set cookie
	// OneWeek := time.Now().Add(time.Hour * 24 * 7)
	// cookie := new(http.Cookie)
	// cookie.Name = "token"
	// cookie.Value = token
	// cookie.Expires = OneWeek
 //  cookie.Domain = "quesocosteno.com"
 //  cookie.Path = "/"
	// cookie.Secure = true
	// cookie.HttpOnly = true
	// c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
    "token": token,
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
func (ac *AuthController) LogoutUser(c echo.Context) error {
	// Simplemente devuelve una respuesta exitosa, indicando al cliente que elimine el token almacenado
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Usuario ha cerrado sesión correctamente",
	})
}

func (ac *AuthController) SessionUser(c echo.Context) error {

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
	if err := ac.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		// Handle the case where the user with the given ID is not found
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	user.Password = ""

	// Return the user data
	return c.JSON(http.StatusOK, user)
}

func (au *AuthController) ConfirmUser(c echo.Context) error {
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

	if _, ok := requestBody["code"]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing confirmation code field")
	}

	email := requestBody["email"].(string)
	confirmationCode := requestBody["code"].(string)

	var user models.User
	if err := au.DB.Where("email = ? AND confirmation_code = ?", email, confirmationCode).First(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email or confirmation code")
	}

	user.IsConfirmed = true
	user.ConfirmationCode = ""

	if err := au.DB.Save(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to confirm user")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Account confirmed successfully",
	})
}

func (ac *AuthController) ResendConfirmationCode(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Cuenta activada correctamente.",
	})
}
