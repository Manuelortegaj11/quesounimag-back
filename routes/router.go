package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoute(e *echo.Echo, db *gorm.DB) {
	SetupInitRoute(e, db)
	SetupAuthRoutes(e, db)
	SetupUserRoutes(e, db)
	SetupAddressRoutes(e, db)
	SetupAuthRoutes(e, db)
	SetupPaymentRoutes(e, db)
	SetupCategoryRoutes(e, db)
	SetupProductRoutes(e, db)
	SetupImagesRoutes(e, db)
  SetupLocationRoutes(e, db)
}

func SetupInitRoute(e *echo.Echo, db *gorm.DB) {
	e.GET("/", func(c echo.Context) error {

		jsonResponse := map[string]interface{}{
			"message": "Xhlar S.A.S",
		}

		return c.JSON(200, jsonResponse)
	})
}
