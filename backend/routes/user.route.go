package routes

import (
	"backend/controllers"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	userRoutes := api.Group("/user")

	userRoutes.Post("/register", controllers.RegisterUser)
	userRoutes.Post("/login", controllers.LoginUser)
	userRoutes.Post("/refresh-token", controllers.RefreshToken) // Refresh token route

	// Protected routes example
	protected := userRoutes.Group("/", middleware.Protect)
	protected.Get("/profile", controllers.GetUserProfile)
}
