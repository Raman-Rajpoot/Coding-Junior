package utils

import "github.com/gofiber/fiber/v2"

func ApiResponse(ctx *fiber.Ctx, status int, data interface{}, message string) error {
	return ctx.Status(status).JSON(fiber.Map{
		"data":    data,
		"message": message,
	})
}
