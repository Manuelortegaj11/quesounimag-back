package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/swaggo/echo-swagger"
	"proyectoqueso/config"
	_ "proyectoqueso/docs"
	"proyectoqueso/routes"
	"gorm.io/gorm"
)

var (
    db *gorm.DB
    err error
)

// @title Xhlar API 
// @version 1.0
// @description Control Panel API for Xhlar
// @termsOfService http://xhlar.com/terminos-condiciones/
// @contact.name API Support
// @contact.email soporte@xhlar.com
// @contact.url http://www.xhlar.com/soporte
// @host xhlar.com
// @BasePath /api/v1/
// @schemes http https
func main() {

  e := echo.New()

  config.InitEnv()
  db, _ := config.NewDB()
  config.Migrate(db)

  // Middleware
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.Use(middleware.CORS())
  

  // Routes
  routes.InitRoute(e, db)

  // Swagger
  e.GET("/swagger/*", echoSwagger.WrapHandler)

  // Start server
	e.Logger.SetLevel(log.INFO)
  e.Logger.Fatal(e.Start(":1323"))
}
