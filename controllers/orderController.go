package controllers

import (
	"net/http"
	"proyectoqueso/models"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type OrderController struct {
	DB *gorm.DB
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{DB: db}
}

func (au *OrderController) CreateOrder(c echo.Context) error {
	// Declara una estructura para almacenar el cuerpo de la solicitud
	type RequestBody struct {
		OrderDate    time.Time `json:"orderDate"`
		OrderTotal   int       `json:"orderTotal"`
		OrderDetails []struct {
			ODProdID   int `json:"odProdId"`
			ODQuantity int `json:"odQuantity"`
			ODPrice    int `json:"odPrice"`
		} `json:"orderDetails"`
		OrderAddress struct {
			ID            int    `json:"ID"`
			CreatedAt     string `json:"CreatedAt"`
			UpdatedAt     string `json:"UpdatedAt"`
			UserID        string `json:"UserID"`
			FullName      string `json:"FullName"`
			PhoneNumber   string `json:"PhoneNumber"`
			Country       string `json:"Country"`
			State         string `json:"State"`
			City          string `json:"City"`
			StreetAddress string `json:"StreetAddress"`
			PostalCode    string `json:"PostalCode"`
		} `json:"orderAddress"`
	}

	// Instancia para almacenar el cuerpo de la solicitud
	var requestBody RequestBody

	// Intenta enlazar el cuerpo de la solicitud a la estructura RequestBody
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if requestBody.OrderTotal == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing orderTotal field")
	}

	// Convertir UserID a uuid.UUID
	userID, err := uuid.Parse(requestBody.OrderAddress.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid user ID format",
		})
	}

	ID_ORDER := uuid.New()

	newOrder := &models.Order{
		ID:          ID_ORDER,
		UserID:      userID,
		TotalAmount: float64(requestBody.OrderTotal),
	}

	newOrderAddress := &models.OrderAddress{
		OrderID:       ID_ORDER,
		UserAddressID: int(requestBody.OrderAddress.ID),
	}

	if err := au.DB.Create(newOrder).Error; err != nil {
		return err
	}

	if err := au.DB.Create(newOrderAddress).Error; err != nil {
		return err
	}

	// Crear y guardar cada detalle del pedido
	for _, detail := range requestBody.OrderDetails {
		newOrderDetail := &models.OrderDetail{
			OrderID:   newOrder.ID,
			ProductID: uint(detail.ODProdID),
			Quantity:  detail.ODQuantity,
		}

		// Guardar el detalle del pedido en la base de datos
		if err := au.DB.Create(newOrderDetail).Error; err != nil {
			return err
		}

		// Actualizar el stock del producto
		var product models.Product
		if err := au.DB.First(&product, detail.ODProdID).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Product not found",
			})
		}

		// Descontar la cantidad solicitada del stock
		product.Stock -= detail.ODQuantity

		// Verificar que el stock no sea negativo
		if product.Stock < 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Insufficient stock for product ID " + strconv.Itoa(detail.ODProdID),
			})
		}

		// Guardar los cambios en el producto
		if err := au.DB.Save(&product).Error; err != nil {
			return err
		}

	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Order created successfully",
	})

}

func (au *OrderController) GetOrdersByUserID(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Authorization header not found",
		})
	}

	// Extract the token from the header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Token not found",
		})
	}

	// Parse the token without validating the signature
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid token",
		})
	}

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid token claims",
		})
	}
	userID := claims.Issuer // Esto debería ser el ID del usuario

	// Convertir userID a UUID
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid user ID format",
		})
	}

	// Consultar las órdenes del usuario en la base de datos
	var orders []models.Order
	if err := au.DB.Preload("OrderDetails.Product").Preload("OrderAddress").Where("user_id = ?", parsedUserID).Find(&orders).Error; err != nil {
		return err
	}

	// Construir la respuesta personalizada
	var orderResponses []map[string]interface{}
	for _, order := range orders {
		var orderDetails []map[string]interface{}
		for _, detail := range order.OrderDetails {
			orderDetails = append(orderDetails, map[string]interface{}{
				"odProdId":   detail.Product.ID,
				"odQuantity": detail.Quantity,
				"odPrice":    detail.Product.Price,
			})
		}
		orderResponse := map[string]interface{}{
			"ID":           order.ID,
			"CreatedAt":    order.CreatedAt,
			"orderTotal":   order.TotalAmount,
			"orderDetails": orderDetails,
			"orderAddress": order.OrderAddress.UserAddressID,
		}
		orderResponses = append(orderResponses, orderResponse)
	}

	// Devolver la respuesta personalizada como JSON
	return c.JSON(http.StatusOK, orderResponses)
}
