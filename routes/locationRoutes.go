package routes

import (
	"proyectoqueso/controllers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupLocationRoutes(e *echo.Echo, db *gorm.DB) {

	locationController := controllers.NewLocationController()

	// Configurar rutas
	e.GET("/countries", locationController.GetCountries)
	e.GET("/states/:countryKey", locationController.GetStates)
	e.GET("/cities/:stateKey", locationController.GetCities)

}
