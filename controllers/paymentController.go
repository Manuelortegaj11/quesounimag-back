package controllers

import (
  // "net/http"
	"gorm.io/gorm"
  "github.com/labstack/echo/v4"
)

type PaymentController struct {
	DB *gorm.DB 
}

func NewPaymentController(db *gorm.DB) *PaymentController {
	return &PaymentController{DB: db}
}

func (uc *PaymentController) CreatePayment(c echo.Context) error {
  return nil
}

// Create a route for update a payment with gorm
func UpdatePayment(c echo.Context) error {
  return nil
}

// Create a route for delete a payment with gorm
func DeletePayment(c echo.Context) error {
  return nil
}

// Create a route for get a payment with gorm
func GetPayment(c echo.Context) error {
  return nil
}

// Create a route for get all payments with gorm
func GetPayments(c echo.Context) error {
  return nil
}

// Create a route for get all payments with gorm
func GetPaymentsByUser(c echo.Context) error {
  return nil
}

