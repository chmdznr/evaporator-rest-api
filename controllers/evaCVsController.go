package controllers

import (
	"be-evaporator/database"
	"be-evaporator/models"
	"be-evaporator/models/reqresp"
	"be-evaporator/utils"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// detectContentType inspects the file content to determine the MIME type
func detectContentType(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Read the first 512 bytes to detect content type
	buffer := make([]byte, 512)
	if _, err := f.Read(buffer); err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

// CreateEvaCv creates a new record for evaporator CV data
// @Summary Create a new record for evaporator CV data
// @Description Create a new record for evaporator CV data
// @Tags Evaporator
// @Accept multipart/form-data
// @Produce json
// @Param trial_code formData string true "Trial Code"
// @Param data_type formData string true "Data Type (image or video)" Enums(image, video)
// @Param notes formData string false "Notes"
// @Param file formData file true "File"
// @Success 201 {object} reqresp.SuccessResponse
// @Failure 400 {object} reqresp.ErrorResponse
// @Failure 500 {object} reqresp.ErrorResponse
// @Router /eva-cv [post]
func CreateEvaCv(c *fiber.Ctx) error {
	// Parse form data
	trialCode := c.FormValue("trial_code")
	dataType := c.FormValue("data_type")
	notes := c.FormValue("notes")
	file, err := c.FormFile("file")

	if err != nil {
		return c.Status(400).JSON(&reqresp.ErrorResponse{
			Status:  "error",
			Message: "Failed to get file: " + err.Error(),
		})
	}

	// Validate data_type
	if dataType != "image" && dataType != "video" {
		return c.Status(400).JSON(&reqresp.ErrorResponse{
			Status:  "error",
			Message: "Invalid data_type, must be either 'image' or 'video'",
		})
	}

	// Validate file MIME type
	fileType, err := detectContentType(file)
	if err != nil {
		return c.Status(400).JSON(&reqresp.ErrorResponse{
			Status:  "error",
			Message: "Failed to detect file type: " + err.Error(),
		})
	}

	if (dataType == "image" && !strings.HasPrefix(fileType, "image/")) ||
		(dataType == "video" && !strings.HasPrefix(fileType, "video/")) {
		return c.Status(400).JSON(&reqresp.ErrorResponse{
			Status:  "error",
			Message: "Invalid file type for the given data_type",
		})
	}

	// Create EvaCv record
	evaCv := models.EvaCv{
		TrialCode: trialCode,
		DataType:  dataType,
		Notes:     notes,
	}

	if err := database.DB.Create(&evaCv).Error; err != nil {
		return c.Status(500).JSON(&reqresp.ErrorResponse{
			Status:  "error",
			Message: "Failed to create evaporator CV: " + err.Error(),
		})
	}

	// Clean the file name to be used in the media table
	cleanName := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename)), "_")
	name := utils.RandString(13) + "_" + cleanName
	fileName := name + filepath.Ext(file.Filename)

	// Create Media record
	media := models.Media{
		ModelType:            "App\\Models\\NewbornCv",
		ModelID:              evaCv.ID,
		UUID:                 uuid.New().String(),
		CollectionName:       "file",
		Name:                 name,
		FileName:             fileName,
		MimeType:             file.Header.Get("Content-Type"),
		Disk:                 "public",
		Size:                 uint(file.Size),
		Manipulations:        `[]`,
		CustomProperties:     `[]`,
		GeneratedConversions: `{"thumb":false,"preview":false}`,
		ResponsiveImages:     `[]`,
		OrderColumn:          1,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := database.DB.Create(&media).Error; err != nil {
		return c.Status(500).JSON(&reqresp.ErrorResponse{
			Status:  "error",
			Message: "Failed to create media record: " + err.Error(),
		})
	}

	// Save file to a specified path named after the media ID
	dir := filepath.Join(os.Getenv("BASE_UPLOAD_FOLDER"), strconv.Itoa(int(media.ID)))
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return c.Status(500).JSON(&reqresp.ErrorResponse{
			Status:  "error",
			Message: "Failed to create directory: " + err.Error(),
		})
	}

	filePath := filepath.Join(dir, fileName)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(&reqresp.ErrorResponse{
			Status:  "error",
			Message: "Failed to save file: " + err.Error(),
		})
	}

	// Set file permissions to -rw-r--r--
	if err := os.Chmod(filePath, 0644); err != nil {
		return c.Status(500).JSON(&reqresp.ErrorResponse{
			Status:  "error",
			Message: "Failed to set file permissions: " + err.Error(),
		})
	}

	return c.Status(201).JSON(&reqresp.SuccessResponse{
		Status:  "success",
		Message: "Berhasil menambahkan data CV evaporator baru",
		Data:    evaCv,
	})
}
