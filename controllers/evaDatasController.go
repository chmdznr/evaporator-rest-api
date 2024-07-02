package controllers

import (
	"be-evaporator/database"
	"be-evaporator/models"
	"be-evaporator/models/reqresp"
	"github.com/gofiber/fiber/v2"
)

// CreateEvaData creates a new record for evaporator sensor data
// @Summary Create a new record for evaporator sensor data
// @Description Create a new record for evaporator sensor data
// @Tags Evaporator
// @Accept json
// @Produce json
// @Param evaData body models.EvaData true "Data Sensor pada evaporator"
// @Success 201 {object} reqresp.SuccessResponse
// @Failure 400 {object} reqresp.ErrorResponse
// @Failure 500 {object} reqresp.ErrorResponse
// @Router /eva-data [post]
func CreateEvaData(c *fiber.Ctx) error {
	evaData := new(models.EvaData)
	if err := c.BodyParser(evaData); err != nil {
		return c.Status(400).JSON(&reqresp.ErrorResponse{
			Status:  "error",
			Message: "Failed to parse request body: " + err.Error(),
		})
	}

	if err := database.DB.Create(&evaData).Error; err != nil {
		return c.Status(500).JSON(&reqresp.ErrorResponse{
			Status:  "error",
			Message: "Failed to create sensor data: " + err.Error(),
		})
	}

	return c.Status(201).JSON(&reqresp.SuccessResponse{
		Status:  "success",
		Message: "Berhasil menambahkan data evaporator",
		Data:    evaData,
	})
}
