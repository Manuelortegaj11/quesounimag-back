// routes/domainRoutes.go
package routes

import (
	"proyectoqueso/controllers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupDomainRoutes configura las rutas relacionadas con los dominios
func SetupDomainRoutes(e *echo.Echo, db *gorm.DB) {

    domainController := controllers.NewDomainController(db)
    apiDomainsGroup := e.Group("/api/v1/domain")

    apiDomainsGroup.GET("", controllers.GetDomains)
    apiDomainsGroup.POST("", domainController.CreateDomain)
    apiDomainsGroup.GET("/:id", controllers.GetDomain)
    // apiDomainsGroup.PUT("/:id", controllers.UpdateDomain)
    // apiDomainsGroup.DELETE("/:id", controllers.DeleteDomain)
}

