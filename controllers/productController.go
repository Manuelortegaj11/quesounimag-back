package controllers

import (
	"fmt"
	"net/http"
	"proyectoqueso/models"
	"regexp"
	"strconv"
	"strings"

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

func generarSlug(nombre string) string {
	// Convertir a minúsculas
	slug := strings.ToLower(nombre)
	// Reemplazar espacios y caracteres especiales con guiones
	slug = regexp.MustCompile(`\s+`).ReplaceAllString(slug, "-")
	// Remover caracteres no deseados
	slug = regexp.MustCompile(`[^\w\-]`).ReplaceAllString(slug, "")
	return slug
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
	requiredFields := []string{"name", "description", "category_id"}
	for _, field := range requiredFields {
		if _, ok := requestBody[field]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing "+field+" field")
		}
	}

	// Genera el slug a partir del nombre
	name := requestBody["name"].(string)
	slug := generarSlug(name)

	// Crea un nuevo producto
	newProduct := &models.Product{
		Name:        name,
		Description: requestBody["description"].(string),
		// Price:       requestBody["price"].(float64),
		// Stock:       int(requestBody["stock"].(float64)),
		CategoryID:  int64(requestBody["category_id"].(float64)),
		Slug:        slug, // Agrega el slug al producto
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
	id := strings.TrimPrefix(c.Param("id"), "/")

	// Check if id is provided
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Product ID is required")
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	// Bind the incoming JSON data to a map
	requestBody := map[string]interface{}{}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Define required fields
	requiredFields := []string{"category_id", "description", "name"}
	for _, field := range requiredFields {
		if _, ok := requestBody[field]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Missing %s field", field))
		}
	}


	// Extract fields from the request body
	categoryID := int64(requestBody["category_id"].(float64))
	description := requestBody["description"].(string)
	name := requestBody["name"].(string)
	slug := generarSlug(name)
	// price := float64(requestBody["price"].(float64))
	// stock := int(requestBody["stock"].(float64))

	// Find the existing product
	var existingProduct models.Product
	if err := au.DB.First(&existingProduct, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Update the product with new values
	existingProduct.CategoryID = categoryID
	existingProduct.Description = description
	existingProduct.Name = name
	// existingProduct.Price = price
	// existingProduct.Stock = stock
  existingProduct.Slug = slug 

	// Save the updated product
	if err := au.DB.Save(&existingProduct).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category updated successfully",
	})
}

func (au *ProductController) DeleteProduct(c echo.Context) error {

	id := strings.TrimPrefix(c.Param("id"), "/")

	// Check if id is provided
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Product ID is required")
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	// Delete the product by ID
	var product models.Product
	if err := au.DB.First(&product, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := au.DB.Delete(&product).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Product deleted successfully"})
}

func (au *ProductController) GetProducts(c echo.Context) error {
    id := c.QueryParam("id")
    slug := c.QueryParam("slug")
    slugCategory := c.QueryParam("slugCategory")
    page := c.QueryParam("page")
    pageSize := c.QueryParam("pageSize")
    search := c.QueryParam("search")

    // Parse page and pageSize into integers
    var pageInt, pageSizeInt int
    if page != "" {
        pageInt, _ = strconv.Atoi(page)
    } else {
        pageInt = 1 // Default to page 1 if not provided
    }
    if pageSize != "" {
        pageSizeInt, _ = strconv.Atoi(pageSize)
    } else {
        pageSizeInt = 10 // Default page size to 10 if not provided
    }

    // Offset calculation for pagination
    offset := (pageInt - 1) * pageSizeInt

    // Check if id is provided
    if id != "" {
        var product models.Product
        if err := au.DB.First(&product, id).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                return echo.NewHTTPError(http.StatusNotFound, "Product not found")
            }
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
        return c.JSON(http.StatusOK, product)
    }

    // Check if slug is provided
    if slug != "" {
        var products []models.Product
        if err := au.DB.Where("slug LIKE ?", "%"+slug+"%").Offset(offset).Limit(pageSizeInt).Find(&products).Error; err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
        return c.JSON(http.StatusOK, products)
    }

    // Check if slugCategory is provided
    if slugCategory != "" {
        var products []models.Product
        if err := au.DB.Joins("Category").Where("Category.slug = ?", slugCategory).Offset(offset).Limit(pageSizeInt).Find(&products).Error; err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
        return c.JSON(http.StatusOK, products)
    }

    // Check if search is provided
    if search != "" {
        var products []models.Product
        if err := au.DB.Where("name LIKE ?", "%"+search+"%").Offset(offset).Limit(pageSizeInt).Find(&products).Error; err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
        return c.JSON(http.StatusOK, products)
    }

    // If neither id, slug, nor slugCategory is provided, return all products
    var products []models.Product
    if err := au.DB.Offset(offset).Limit(pageSizeInt).Find(&products).Error; err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusOK, products)
}
