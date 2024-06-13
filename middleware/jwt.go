package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func CheckCookie(c echo.Context) bool {
    cookie, err := c.Cookie("token")
    if err != nil {
        c.Logger().Errorf("Error fetching cookie: %v", err)
        return false
    }

    cookieValue := cookie.Value

    token, err := jwt.ParseWithClaims(cookieValue, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        c.Logger().Errorf("Error parsing token: %v", err)
        return false
    }

    if !token.Valid {
        c.Logger().Errorf("Invalid token")
        return false
    }

    return true
}

func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        if !CheckCookie(c) {
            return c.JSON(http.StatusUnauthorized, map[string]string{
                "message": "Cookie Unauthorized",
            })
        }
        return next(c)
    }
}

func CheckRole(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c echo.Context) error {
    check := CheckCookie(c)
    if !check {
      return c.JSON(http.StatusUnauthorized, map[string]string{
        "message": "Cookie Unauthorized",
      })  
    }

    return next(c)
  }
}
