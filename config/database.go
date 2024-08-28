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

var modelsToMigrate = []interface{}{
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
	&models.OrderDetail{},
	&models.OrderAddress{},
	&models.Category{},
	&models.CollectionCenterInventory{},
	&models.CollectionCenter{},
}

func Migrate(db *gorm.DB) (*gorm.DB, error) {
	if err := db.AutoMigrate(modelsToMigrate...); err != nil {
		return nil, err
	}

	db.Model(&models.UserAddress{}).Association("User")
	db.Model(&models.Order{}).Association("User")
	db.Model(&models.Payment{}).Association("User")
	db.Model(&models.Payment{}).Association("Product")
	db.Model(&models.OrderAddress{}).Association("Order")
	db.Model(&models.OrderAddress{}).Association("UserAddress")

	return db, nil
}

func DropAllTables(db *gorm.DB) error {
	var err error
	tables := []interface{}{
    modelsToMigrate,
	}

	for _, table := range tables {
		if err = db.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	return nil
}

func CreateTestUsers(db *gorm.DB) error {

	hashedPassword, err := security.EncryptPassword("123456")
	if err != nil {
		log.Fatal("Error encrypting password: ", err)
		return err
	}

	newUserID := uuid.New()
	newUser2ID := uuid.New()
	newUser3ID := uuid.New()
	users := []models.User{
		{
			ID:        newUserID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Password:  string(hashedPassword),
      Roles: []models.Role{
        {
          Name: "Proveedor",
        },
      },
      Permissions: []models.Permission{
        {
          Name: "create_products",
        },
        {
          Name: "read_products",
        },
        {
          Name: "update_products",
        },
        {
          Name: "delete_products",
        },
      },
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
				{
					UserID:        newUserID,
					FullName:      "Jane Smith",
					PhoneNumber:   "654321",
					Country:       "Argentina ",
					State:         "Buenos Aires",
					City:          "Buenos aires ",
					StreetAddress: "Calle 10",
					PostalCode:    "00000",
				},
			},
		},
    // SEGUNDO USUARIO
		{
			ID:        newUser3ID,
			FirstName: "Michael ",
			LastName:  "Smith",
			Email:     "michael@example.com",
			Password:  string(hashedPassword),
			Addresses: []models.UserAddress{
				{
					UserID:        newUser3ID,
					FullName:      "Michael Smith",
					PhoneNumber:   "123456",
					Country:       "United States",
					State:         "California",
					City:          "Los Angeles",
					StreetAddress: "123 Main Street",
					PostalCode:    "90210",
				},
			},
		},
    // TERCER USUARIO
		{
			ID:        newUser2ID,
			FirstName: "Jane Smith",
			LastName:  "Smith",
			Email:     "jane@example.com",
			Password:  string(hashedPassword),
      Roles: []models.Role{
        {
          Name: "Administrador",
        },
      },
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

func DropTestUsers(db *gorm.DB) error {
    // Correos electrónicos de usuarios de prueba
    testUserEmails := []string{"john@example.com", "jane@example.com", "michael@example.com"}

    // // Eliminar usuarios de prueba
    if err := db.Where("email IN (?)", testUserEmails).Delete(&models.User{}).Error; err != nil {
        return err
    }

    // Obtener los UUIDs de todos los usuarios de prueba
    var testUserUUIDs []string
    if err := db.
        Unscoped().
        Model(&models.User{}).
        Where("email IN (?)", testUserEmails).
        Pluck("id", &testUserUUIDs).
        Error; err != nil {
        return err
    }

    // Imprimir los UUIDs obtenidos
    for _, uuid := range testUserUUIDs {
        log.Println("UUID de usuario de prueba:", uuid)
    }

    // Borrar permanentemente las direcciones (Unscoped)
    if err := db.
        Unscoped().
        Where("user_id IN (?)", testUserUUIDs).
        Delete(&models.UserAddress{}).
        Error; err != nil {
        return err
    }

    // Borrar permanentemente los usuarios eliminados (Unscoped)
    if err := db.
        Unscoped().
        Where("email IN (?)", testUserEmails).
        Delete(&models.User{}).
        Error; err != nil {
        return err
    }

    log.Println("Usuarios de prueba borrados permanentemente")

    return nil
}
