package controllers

import (
	"fmt"
	"net/http"
	"proyectoqueso/models"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{DB: db}
}

func generarSlug(nombre string) string {
	slug := strings.ToLower(nombre)
	slug = regexp.MustCompile(`\s+`).ReplaceAllString(slug, "-")
	slug = regexp.MustCompile(`[^\w\-]`).ReplaceAllString(slug, "")
	return slug
}

func (au *ProductController) CreateProduct(c echo.Context) error {
	var requestBody map[string]interface{}

	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if requestBody == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing JSON body")
	}

	requiredFields := []string{"name", "description", "category_id", "price", "price_min", "price_max"}
	for _, field := range requiredFields {
		if _, ok := requestBody[field]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing "+field+" field")
		}
	}

	name := requestBody["name"].(string)
	slug := generarSlug(name)

	newProduct := &models.Product{
		Name:        name,
		Description: requestBody["description"].(string),
		Price:       requestBody["price"].(float64),
		PriceMin:    requestBody["price_min"].(float64),
		PriceMax:    requestBody["price_max"].(float64),
		CategoryID:  int64(requestBody["category_id"].(float64)),
		Slug:        slug,
	}

	if err := au.DB.Create(newProduct).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Product created successfully",
	})
}

func (au *ProductController) UpdateProduct(c echo.Context) error {
	id := strings.TrimPrefix(c.Param("id"), "/")

	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Product ID is required")
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	requestBody := map[string]interface{}{}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	requiredFields := []string{"category_id", "description", "name", "price", "price_min", "price_max"}
	for _, field := range requiredFields {
		if _, ok := requestBody[field]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Missing %s field", field))
		}
	}

	categoryID := int64(requestBody["category_id"].(float64))
	description := requestBody["description"].(string)
	name := requestBody["name"].(string)
	price := requestBody["price"].(float64)
	priceMin := requestBody["price_min"].(float64)
	priceMax := requestBody["price_max"].(float64)
	slug := generarSlug(name)

	var existingProduct models.Product
	if err := au.DB.First(&existingProduct, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	existingProduct.CategoryID = categoryID
	existingProduct.Description = description
	existingProduct.Name = name
	existingProduct.Price = price
	existingProduct.PriceMin = priceMin
	existingProduct.PriceMax = priceMax
	existingProduct.Slug = slug

	if err := au.DB.Save(&existingProduct).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Product updated successfully",
	})
}

func (au *ProductController) DeleteProduct(c echo.Context) error {
	id := strings.TrimPrefix(c.Param("id"), "/")

	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Product ID is required")
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

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

	var pageInt, pageSizeInt int
	if page != "" {
		pageInt, _ = strconv.Atoi(page)
	} else {
		pageInt = 1
	}
	if pageSize != "" {
		pageSizeInt, _ = strconv.Atoi(pageSize)
	} else {
		pageSizeInt = 10
	}

	offset := (pageInt - 1) * pageSizeInt
	var totalItems int64

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

	var products []models.Product

	if slug != "" {
		if err := au.DB.Model(&models.Product{}).Where("slug LIKE ?", "%"+slug+"%").Count(&totalItems).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if err := au.DB.Where("slug LIKE ?", "%"+slug+"%").Offset(offset).Limit(pageSizeInt).Find(&products).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"totalItems": totalItems,
			"data":       products,
		})
	}

	if slugCategory != "" {
		if err := au.DB.Joins("Category").Where("Category.slug = ?", slugCategory).Model(&models.Product{}).Count(&totalItems).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if err := au.DB.Joins("Category").Where("Category.slug = ?", slugCategory).Offset(offset).Limit(pageSizeInt).Find(&products).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"totalItems": totalItems,
			"data":       products,
		})
	}

	if search != "" {
		if err := au.DB.Model(&models.Product{}).Where("name LIKE ?", "%"+search+"%").Count(&totalItems).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if err := au.DB.Where("name LIKE ?", "%"+search+"%").Offset(offset).Limit(pageSizeInt).Find(&products).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"totalItems": totalItems,
			"data":       products,
		})
	}

	if err := au.DB.Model(&models.Product{}).Count(&totalItems).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := au.DB.Offset(offset).Limit(pageSizeInt).Find(&products).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"totalItems": totalItems,
		"data":       products,
	})
}
