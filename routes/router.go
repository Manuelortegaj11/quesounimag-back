package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoute(e *echo.Echo, db *gorm.DB){
  SetupAuthRoutes(e, db)
  SetupUserRoutes(e, db)
  SetupAuthRoutes(e, db)
  SetupPaymentRoutes(e, db)
}
