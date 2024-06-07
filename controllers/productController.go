package controllers

import (
	"net/http"
	"proyectoqueso/models"
	"strconv"

	"github.com/labstack/echo/v4"

	// "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{DB: db}
}

func (au *ProductController) CreateProduct(c echo.Context) error {
	// Declara una variable para almacenar el cuerpo de la solicitud
	var requestBody map[string]interface{}

	// Intenta unir el cuerpo de la solicitud al mapa
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	// Verifica si el cuerpo de la solicitud es nulo
	if requestBody == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing JSON body")
	}

	// Verifica si se proporcionan los campos necesarios en el cuerpo de la solicitud
	requiredFields := []string{"name", "description", "price", "stock", "category_id"}
	for _, field := range requiredFields {
		if _, ok := requestBody[field]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing "+field+" field")
		}
	}

	// Crea un nuevo producto
	newProduct := &models.Product{
		Name:        requestBody["name"].(string),
		Description: requestBody["description"].(string),
		Price:       requestBody["price"].(float64),
		Stock:       int(requestBody["stock"].(float64)),
		CategoryID:  int64(requestBody["category_id"].(float64)),
	}

	// Guarda el nuevo producto en la base de datos
	if err := au.DB.Create(newProduct).Error; err != nil {
		return err
	}

	// Devuelve una respuesta JSON indicando que el producto se creó correctamente
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Product created successfully",
	})
}

func (au *ProductController) UpdateProduct(c echo.Context) error {

	var requestBody map[string]interface{}

	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if requestBody == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing JSON body")
	}

	if _, ok := requestBody["oldName"]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing oldName field")
	}

	if _, ok := requestBody["newName"]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing newName field")
	}

	newName := requestBody["newName"].(string)
	oldName := requestBody["oldName"].(string)

	// Check if category already exists
	var category models.Category

	queryCategory := au.DB.Where("name = ?", oldName).First(&category)

	if queryCategory.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}

	queryExistCategory := au.DB.Where("name = ?", newName).First(&category)

	if queryExistCategory.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category already exist")
	}

	category.Name = newName

	if err := au.DB.Save(&category).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category updated successfully",
	})
}

func (au *ProductController) DeleteProduct(c echo.Context) error {
	var requestBody map[string]interface{}

	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if requestBody == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing JSON body")
	}

	if _, ok := requestBody["name"]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing name field")
	}

	name := requestBody["name"].(string)

	// Check if category exists
	var category models.Category

	queryCategory := au.DB.Where("name = ?", name).First(&category)

	if queryCategory.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}

	// Delete the category
	if err := au.DB.Delete(&category).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category deleted successfully",
	})
}

func (au *ProductController) GetProducts(c echo.Context) error {
	id := c.QueryParam("id")
	slug := c.QueryParam("slug")

	// Check if id is provided
	if id != "" {
		categoryID, err := strconv.Atoi(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid category ID")
		}

		// Retrieve a single category by ID
		var category models.Category
		if err := au.DB.First(&category, categoryID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return echo.NewHTTPError(http.StatusNotFound, "Category not found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, category)
	}

	// Check if slug is provided
	if slug != "" {
		var categories []models.Category
		if err := au.DB.Where("name LIKE ?", "%"+slug+"%").Find(&categories).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, categories)
	}

	// Retrieve all categories if no ID or slug is provided
	var categories []models.Category
	if err := au.DB.Find(&categories).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, categories)
}
