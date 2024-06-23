// cmd/database/main.go
package main

import (
	"log"
	"os"
	"proyectoqueso/config"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "drop" {
		if err := config.DropAllTables(db); err != nil {
			log.Fatalf("failed to drop tables: %v", err)
		}
		log.Println("All tables dropped successfully")
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "createtestusers" {
		if err := config.CreateTestUsers(db); err != nil {
			log.Fatalf("failed to create user: %v", err)
		}
		log.Println("All test users created successfully")
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "droptestusers" {
		if err := config.DropTestUsers(db); err != nil {
      log.Fatalf("Error dropping test users: %v", err)
		}
		log.Println("All test users deleted successfully")
		return
	}

	if _, err := config.Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Migration completed successfully")
}
