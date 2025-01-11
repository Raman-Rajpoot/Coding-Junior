package middleware

import (
	"backend/utils"
	"github.com/gofiber/fiber/v2"
)

// Protect is the middleware to secure routes
func Protect(ctx *fiber.Ctx) error {
	// Retrieve the access token from cookies
	accessToken := ctx.Cookies("access_token")

	if accessToken == "" {
		return utils.NewApiError(fiber.StatusUnauthorized, "Missing access token").Handle(ctx)
	}

	// Validate the access token
	claims, err := utils.ValidateToken(accessToken, "access")
	if err != nil {
		return utils.NewApiError(fiber.StatusUnauthorized, "Invalid or expired access token").Handle(ctx)
	}

	// Store user claims in context for subsequent handlers
	ctx.Locals("user", claims)

	// Proceed to the next handler
	return ctx.Next()
}
