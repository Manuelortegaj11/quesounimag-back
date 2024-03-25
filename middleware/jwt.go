package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func checkCookie(c echo.Context) bool {
  cookie, err := c.Cookie("jwt")
  if err != nil {
    return false
  }

  cookieValue := cookie.Value

  token, err := jwt.ParseWithClaims(cookieValue, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
    return jwtKey, nil
  })

  if err != nil {
    return false
  }

  if !token.Valid {
    return false
  }

  return true
}

func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        check := checkCookie(c)
        if !check {
          return c.JSON(http.StatusUnauthorized, map[string]string{
            "message": "Cookie Unauthorized",
          })  
        }
        return next(c)
    }
}

func CheckRole(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c echo.Context) error {
    check := checkCookie(c)
    if !check {
      return c.JSON(http.StatusUnauthorized, map[string]string{
        "message": "Cookie Unauthorized",
      })  
    }

    return next(c)
  }
}
