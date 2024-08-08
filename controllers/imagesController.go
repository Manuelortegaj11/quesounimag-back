package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	maxFileSize = 10 << 20 // 10 MB
	productDir  = "./uploads/products/"
	categoryDir = "./uploads/categories/"
	avatarDir   = "./uploads/avatars/"
)

type ImageController struct {
	DB *gorm.DB
}

func NewImageController(db *gorm.DB) *ImageController {
	return &ImageController{DB: db}
}

func (ic *ImageController) GetImage(c echo.Context) error {
	imageType := c.Param("type")
	imageId := c.Param("id")
	var imageDir string

	switch imageType {
	case "product":
		imageDir = productDir
	case "category":
		imageDir = categoryDir
	case "avatar":
		imageDir = avatarDir
	default:
		return c.String(http.StatusBadRequest, "Tipo de imagen no válido")
	}

	imagePath := filepath.Join(imageDir, imageId)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return c.String(http.StatusNotFound, "Imagen no encontrada")
	}

	return c.File(imagePath)
}

func (ic *ImageController) UpdateImageProduct(c echo.Context) error {
	imageId := c.Param("id")
	return uploadImage(c, productDir, imageId)
}

func (ic *ImageController) UpdateImageCategory(c echo.Context) error {
	imageId := c.Param("id")
	return uploadImage(c, categoryDir, imageId)
}

func (ic *ImageController) UpdateImageAvatar(c echo.Context) error {
	imageId := c.Param("id")
	return uploadImage(c, avatarDir, imageId)
}

func uploadImage(c echo.Context, uploadDir string, filename string) error {
	// Verificar que el directorio de subida exista, de lo contrario crearlo
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// Obtener el archivo de la solicitud
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println("Error al obtener el archivo:", err)
		return c.String(http.StatusBadRequest, "No se pudo obtener el archivo")
	}

	fmt.Println("Archivo recibido:", file.Filename)

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

	// Si la extensión ya está presente, no la agregamos de nuevo
	if !strings.HasSuffix(filename, filepath.Ext(sanitizedFilename)) {
		filename += filepath.Ext(sanitizedFilename)
	}

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
