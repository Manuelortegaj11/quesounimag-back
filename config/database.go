// database_config.go
package config

import (
	"log"
	"os"
	"proyectoqueso/models"
	"proyectoqueso/security"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func InitDB(db *gorm.DB) *gorm.DB {
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

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) (*gorm.DB, error) {
	if err := db.AutoMigrate(
		&models.User{},
		&models.UserAddress{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.UserPermission{},
		&models.Product{},
		&models.Payment{},
		&models.Order{},
		&models.OrdersDetail{},
		&models.Category{},
	); err != nil {
		return nil, err
	}

	db.Model(&models.UserAddress{}).Association("User")
	db.Model(&models.Order{}).Association("User")
	db.Model(&models.Order{}).Association("Product")
	db.Model(&models.Payment{}).Association("User")
	db.Model(&models.Payment{}).Association("Product")

	return db, nil
}

func DropAllTables(db *gorm.DB) error {
	var err error
	tables := []interface{}{
		&models.User{},
		&models.UserAddress{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.UserPermission{},
		&models.Product{},
		&models.Payment{},
		&models.Order{},
		&models.OrdersDetail{},
		&models.Category{},
	}

	for _, table := range tables {
		if err = db.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	return nil
}

func CreateTestUser(db *gorm.DB) error {

	hashedPassword, err := security.EncryptPassword("123456")
	if err != nil {
		log.Fatal("Error encrypting password: ", err)
		return err
	}

	newUserID := uuid.New()
	newUser2ID := uuid.New()
	users := []models.User{
		{
			ID:        newUserID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Password:  string(hashedPassword),
			Addresses: []models.UserAddress{
				{
					UserID:        newUserID,
					FullName:      "John Doe",
					PhoneNumber:   "123456",
					Country:       "United States",
					State:         "California",
					City:          "Los Angeles",
					StreetAddress: "123 Main Street",
					PostalCode:    "90210",
				},
			},
		},
		{
			ID:        newUser2ID,
			FirstName: "Jane Smith",
			LastName:  "Smith",
			Email:     "jane@example.com",
			Password:  string(hashedPassword),
			Addresses: []models.UserAddress{
				{
					UserID:        newUser2ID,
					FullName:      "Jane Smith",
					PhoneNumber:   "123456",
					Country:       "United States",
					State:         "California",
					City:          "Los Angeles",
					StreetAddress: "123 Main Street",
					PostalCode:    "90210",
				},
			},
		},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			return err
		}
	}
	return nil
}
