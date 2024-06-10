package main

import (
	"proyectoqueso/config"
	_ "proyectoqueso/docs"
	"proyectoqueso/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

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
