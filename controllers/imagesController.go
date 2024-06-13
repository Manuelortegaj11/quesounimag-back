package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	//"proyectoqueso/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	maxFileSize = 10 << 20 // 10 MB
	avatarDir   = "./uploads/avatars/"
	productDir  = "./uploads/products/"
)

type Imagecontroller struct {
	DB *gorm.DB
}

func NewImageController(db *gorm.DB) *Imagecontroller {
	return &Imagecontroller{DB: db}
}

func (au *Imagecontroller) GetImage(c echo.Context) error {
	// Devuelve una respuesta JSON indicando que el producto se creó correctamente
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Product created successfully",
	})
}

func (au *Imagecontroller) UploadImage(c echo.Context) error {
	// Devuelve una respuesta JSON indicando que el producto se creó correctamente
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Product created successfully",
	})
}

// Handler para subir avatares
func uploadAvatar(c echo.Context) error {
	return uploadImage(c, avatarDir)
}

// Handler para subir imágenes de productos
func uploadProduct(c echo.Context) error {
	return uploadImage(c, productDir)
}

// Función genérica para subir imágenes
func uploadImage(c echo.Context, uploadDir string) error {
	// Verificar que el directorio de subida exista, de lo contrario crearlo
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// Obtener el archivo de la solicitud
	file, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, "No se pudo obtener el archivo")
	}

	// Verificar el tipo MIME del archivo
	if file.Header.Get("Content-Type") != "image/jpeg" && file.Header.Get("Content-Type") != "image/png" {
		return c.String(http.StatusBadRequest, "Solo se permiten archivos JPEG o PNG")
	}

	// Limitar el tamaño del archivo
	if file.Size > maxFileSize {
		return c.String(http.StatusBadRequest, "El archivo es demasiado grande. El tamaño máximo permitido es de 10 MB")
	}

	// Abrir el archivo
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "No se pudo abrir el archivo")
	}
	defer src.Close()

	// Sanitizar el nombre del archivo
	sanitizedFilename := sanitizeFilename(file.Filename)

	// Crear un nombre de archivo único
	filename := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(sanitizedFilename))

	// Crear un destino en el servidor para guardar el archivo
	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		return c.String(http.StatusInternalServerError, "No se pudo crear el archivo en el servidor")
	}
	defer dst.Close()

	// Copiar el contenido del archivo al destino
	if _, err = io.Copy(dst, src); err != nil {
		return c.String(http.StatusInternalServerError, "No se pudo guardar el archivo")
	}

	return c.String(http.StatusOK, fmt.Sprintf("Archivo subido con éxito: %s", filename))
}

// Sanitizar el nombre del archivo para evitar inyecciones de ruta
func sanitizeFilename(filename string) string {
	// Eliminar directorios y solo mantener el nombre del archivo
	filename = filepath.Base(filename)

	// Reemplazar espacios con guiones bajos
	filename = strings.ReplaceAll(filename, " ", "_")

	// Eliminar caracteres no permitidos usando una expresión regular
	re := regexp.MustCompile(`[^\w\-.]`)
	filename = re.ReplaceAllString(filename, "")

	return filename
}
