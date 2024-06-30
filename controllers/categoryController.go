package controllers

import (
	"net/http"
	"proyectoqueso/models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	// "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CategoryController struct {
	DB *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{DB: db}
}

func createSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}

func (au *CategoryController) CreateCategory(c echo.Context) error {
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
	slug := createSlug(name)

	// Check if category already exists
	var category models.Category
	queryUser := au.DB.Where("name = ?", name).First(&category)

	if queryUser.Error == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Category already exists")
	}

	newCategory := &models.Category{
		Name: name,
		Slug: slug,
	}

	if err := au.DB.Create(newCategory).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category created successfully",
	})
}

func (au *CategoryController) UpdateCategory(c echo.Context) error {
	id := strings.Trim(c.Param("id"), "/")

	var requestBody map[string]interface{} = make(map[string]interface{})

	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if requestBody == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing JSON body")
	}

	if _, ok := requestBody["name"]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing Name field")
	}

	newName := requestBody["name"].(string)
	slug := createSlug(newName)

	c.Logger().Info("Category updated successfully", newName)

	// Check if category exists
	var category models.Category
	queryCategory := au.DB.First(&category, id)

	if queryCategory.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}

	// Check if the new name is already taken by another category
	var existingCategory models.Category
	queryExistCategory := au.DB.Where("name = ?", newName).First(&existingCategory)

	if queryExistCategory.Error == nil && existingCategory.ID != category.ID {
		return echo.NewHTTPError(http.StatusBadRequest, "Category with the new name already exists")
	}

	category.Name = newName
	category.Slug = slug

	if err := au.DB.Save(&category).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category updated successfully",
	})
}

func (au *CategoryController) DeleteCategory(c echo.Context) error {
	// Obtener el ID de la categoría desde el parámetro de la URL
	idParam := strings.Trim(c.Param("id"), "/")

	// Verificar que el ID sea un número válido
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid category ID")
	}

	// Buscar la categoría por ID
	var category models.Category
	queryCategory := au.DB.Where("id = ?", id).First(&category)

	if queryCategory.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}

	// Eliminar la categoría
	if err := au.DB.Delete(&category).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category deleted successfully",
	})
}

func (au *CategoryController) GetCategories(c echo.Context) error {
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
