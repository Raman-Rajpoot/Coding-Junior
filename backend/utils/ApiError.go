package utils

import "github.com/gofiber/fiber/v2"

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
// NewApiError creates a new ApiError instance
func NewApiError(status int, message string) *ApiError {
	return &ApiError{Status: status, Message: message}
}

func (e *ApiError) Handle(ctx *fiber.Ctx) error {
	return ctx.Status(e.Status).JSON(fiber.Map{"error": e.Message})
}
