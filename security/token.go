package security

import (
	"errors"
	"fmt"
	"os"
	util "proyectoqueso/utils"
	"strings"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	JwtKey     = []byte(os.Getenv("JWT_SECRET_KEY"))
	JwtSigningMethod = jwt.SigningMethodHS256.Name
)

func NewToken(userId string) (string, error) {
	claims := jwt.StandardClaims{
		Id:        userId,
		Issuer:    userId,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func validateSignedMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return JwtKey, nil
}

func ParseToken(tokenString string) (*jwt.StandardClaims, error) {
	claims := new(jwt.StandardClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, validateSignedMethod)
	if err != nil {
		return nil, err
	}
	var ok bool
	claims, ok = token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return nil, util.ErrInvalidAuthToken
	}
	return claims, nil
}

func GetUserIDFromToken(c echo.Context) (uuid.UUID, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return uuid.Nil, errors.New("authorization header not found")
	}

	// Extract the token from the header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return uuid.Nil, errors.New("token not found")
	}

	// Parse the token without validating the signature
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		return uuid.Nil, errors.New("invalid token")
	}

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims")
	}

	// Parse userID as UUID
	userID, err := uuid.Parse(claims.Issuer)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID format")
	}

	return userID, nil
}
