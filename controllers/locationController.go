package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type LocationController struct{}

var countries = []map[string]string{
	{"key": "ar", "value": "ar", "flag": "ar", "text": "Argentina"},
	{"key": "au", "value": "au", "flag": "au", "text": "Australia"},
	{"key": "be", "value": "be", "flag": "be", "text": "Belgium"},
	// Añade más países según sea necesario
}

var states = map[string][]map[string]string{
	"ar": {
		{"key": "BA", "value": "BA", "text": "Buenos Aires"},
		{"key": "CBA", "value": "CBA", "text": "Córdoba"},
	},
	"br": {
		{"key": "SP", "value": "SP", "text": "São Paulo"},
		{"key": "RJ", "value": "RJ", "text": "Rio de Janeiro"},
	},
	// Añade más estados según sea necesario
}

var cities = map[string][]map[string]string{
	"BA": {
		{"key": "CAP", "value": "CAP", "text": "Capital Federal"},
		{"key": "LPA", "value": "LPA", "text": "La Plata"},
	},
	"SP": {
		{"key": "SAO", "value": "SAO", "text": "São Paulo"},
		{"key": "CAMP", "value": "CAMP", "text": "Campinas"},
	},
	// Añade más ciudades según sea necesario
}

func NewLocationController() *LocationController {
	return &LocationController{}
}

// GetCountries returns all countries
func (lc *LocationController) GetCountries(c echo.Context) error {
	return c.JSON(http.StatusOK, countries)
}

// GetStates returns states for a given country key
func (lc *LocationController) GetStates(c echo.Context) error {
	countryKey := c.Param("countryKey")
	statesForCountry, exists := states[countryKey]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "States not found for the given country")
	}
	return c.JSON(http.StatusOK, statesForCountry)
}

// GetCities returns cities for a given state key
func (lc *LocationController) GetCities(c echo.Context) error {
	stateKey := c.Param("stateKey")
	citiesForState, exists := cities[stateKey]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "Cities not found for the given state")
	}
	return c.JSON(http.StatusOK, citiesForState)
}
