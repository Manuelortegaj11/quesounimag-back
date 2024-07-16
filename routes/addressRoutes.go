package routes

import (
	"proyectoqueso/controllers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupAddressRoutes(e *echo.Echo, db *gorm.DB) {

	addressController := controllers.NewAddressController(db)
	apiAddressGroup := e.Group("/v1/address")
  apiAddressGroup.GET(":id", addressController.GetAddressByID)
	apiAddressGroup.GET("", addressController.GetAllAddress)
	apiAddressGroup.POST("", addressController.CreateAddress)
	apiAddressGroup.PUT(":id", addressController.UpdateAddress)
	apiAddressGroup.DELETE(":id", addressController.DeleteAddress)
}
