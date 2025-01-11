package utils

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func AsyncHandler(handler func(ctx *fiber.Ctx) error) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		go func() {
			if err := handler(ctx); err != nil {
				log.Printf("Error in async handler: %v", err)
			}
		}()
		return nil
	}
}
