package controllers

import (
	"be-evaporator/models/reqresp"
	"github.com/gofiber/fiber/v2"
)

// GetPing Get ping
// @Summary Get ping message
// @Description Get ping message to check if server is alive
// @Tags Ping
// @Accept json
// @Produce json
// @Success 200 {object} reqresp.SuccessResponse
// @Router /ping [get]
func GetPing(c *fiber.Ctx) error {
	return c.JSON(&reqresp.SuccessResponse{
		Status:  "success",
		Message: "pong",
		Data:    nil,
	})
}
