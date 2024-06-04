// database_config.go
package config

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "proyectoqueso/models"
    "os"
)

var (
  db *gorm.DB
  err error
  )

func InitDB(db *gorm.DB) (*gorm.DB) {
  NewDB()
  Migrate(db)
  return db 
}

func NewDB() (*gorm.DB, error) {

    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    // dsn := "root:123456@tcp(localhost:3306)/xhlartest?charset=utf8mb4&parseTime=True&loc=Local"
    dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    return db, nil


}

func Migrate(db *gorm.DB) (*gorm.DB, error){
    if err := db.AutoMigrate(
    &models.User{}, 
    &models.Product{}, 
    &models.Payment{}, 
    &models.Order{},
    &models.Category{},
    ); 

    err != nil {
        return nil, err 
    }

    db.Model(&models.Order{}).Association("User")
    db.Model(&models.Order{}).Association("Product")
    db.Model(&models.Payment{}).Association("User")
    db.Model(&models.Payment{}).Association("Product")

    return db, nil
}

