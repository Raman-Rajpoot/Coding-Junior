package main

import (
	"backend/db"    // Correct import path
	"backend/routes" // Correct import path
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Initialize MongoDB
	db.ConnectDB()
    println("db connected")
	// Create a new Fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Allow your frontend origin
		AllowMethods:     "GET,POST,PUT,DELETE",  // Allowed HTTP methods
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true, // Allow cookies if needed
	}))
	// Register routes
	routes.SetupRoutes(app)
	println("routes connected")
	// Start the server
	log.Fatal(app.Listen(":8000"))
}
