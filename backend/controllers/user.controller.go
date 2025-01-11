package controllers

import (
	"backend/db"
	"backend/models"
	"backend/utils"
	"context"
	"time"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
)

// RegisterUser handles user registration
func RegisterUser(ctx *fiber.Ctx) error {
	println("RegisterUser endpoint hit!")
	collection := db.GetCollection("users")
	data := new(models.User)

	// Parse the request body into the data struct
	if err := ctx.BodyParser(&data); err != nil {
		return utils.NewApiError(fiber.StatusBadRequest, "Invalid input").Handle(ctx)
	}

	// Validate input data
	println(data.UserName, data.Email, data.Password, data.FullName)
	if data.UserName == "" || data.Email == "" || data.Password == "" || data.FullName == "" {
		return utils.NewApiError(fiber.StatusBadRequest, "All fields are required").Handle(ctx)
	}

	// Check if user already exists in database
	var existingUser models.User
	err := collection.FindOne(context.TODO(), bson.M{"$or": []bson.M{{"email": data.Email}, {"userName": data.UserName}}}).Decode(&existingUser)
	if err == nil {
		return utils.NewApiError(fiber.StatusBadRequest, "User with this email or username already exists").Handle(ctx)
	}

	// Hash the password
	println("Hashing password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		println(err)
		return utils.NewApiError(fiber.StatusInternalServerError, "Failed to hash password").Handle(ctx)
	}
	data.Password = string(hashedPassword)
	data.CreatedAt = time.Now()

	// Insert the new user into the database
	println("Storing data")
	_, err = collection.InsertOne(context.TODO(), data)
	if err != nil {
		println(err)
		return utils.NewApiError(fiber.StatusInternalServerError, "Failed to create user").Handle(ctx)
	}

	// Clear the password from the response data before sending it back
	data.Password = ""
	return utils.ApiResponse(ctx, fiber.StatusCreated, data, "User registered successfully")
}

// LoginUser handles user login
func LoginUser(ctx *fiber.Ctx) error {
	collection := db.GetCollection("users")
	data := new(models.LoginRequest)

	// Parse the request body into the data struct
	if err := ctx.BodyParser(&data); err != nil {
		return utils.NewApiError(fiber.StatusBadRequest, "Invalid input").Handle(ctx)
	}

	// Validate input
	if (data.UserName == "" && data.Email == "") || data.Password == "" {
		return utils.NewApiError(fiber.StatusBadRequest, "All fields are required").Handle(ctx)
	}

	// Find the user in the database
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"$or": []bson.M{{"email": data.Email}, {"userName": data.UserName}}}).Decode(&user)
	if err != nil {
		return utils.NewApiError(fiber.StatusNotFound, "User not found").Handle(ctx)
	}

	// Compare the provided password with the stored password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return utils.NewApiError(fiber.StatusUnauthorized, "Invalid password").Handle(ctx)
	}

	// Generate access token and refresh token
	accessToken, err := utils.GenerateAccessToken(user.UserName, user.Email)
	if err != nil {
		return utils.NewApiError(fiber.StatusInternalServerError, "Could not generate access token").Handle(ctx)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UserName, user.Email)
	if err != nil {
		return utils.NewApiError(fiber.StatusInternalServerError, "Could not generate refresh token").Handle(ctx)
	}
	// Set tokens in HTTP-only cookies
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   3600, // 1 hour
		SameSite: "Strict",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   24 * 3600, // 1 day
		SameSite: "Strict",
	})


	// Return the tokens
	return utils.ApiResponse(ctx, fiber.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"userName" : user.UserName,
		"email" : user.Email,
		"fullName" : user.FullName,
	}, "Login successful")
}

// RefreshToken handles refreshing the access token using the refresh token
func RefreshToken(ctx *fiber.Ctx) error {
	// Parse the request body to get the refresh token
	data := new(struct {
		RefreshToken string `json:"refresh_token"`
	})

	if err := ctx.BodyParser(data); err != nil {
		return utils.NewApiError(fiber.StatusBadRequest, "Invalid input").Handle(ctx)
	}

	// Validate the refresh token
	claims, err := utils.ValidateToken(data.RefreshToken, "refresh")
	if err != nil {
		return utils.NewApiError(fiber.StatusUnauthorized, "Invalid or expired refresh token").Handle(ctx)
	}

	// Generate a new access token using the claims from the refresh token
	accessToken, err := utils.GenerateAccessToken(claims.UserName, claims.Email)
	if err != nil {
		return utils.NewApiError(fiber.StatusInternalServerError, "Could not generate new access token").Handle(ctx)
	}
	
// Set the new access token in an HTTP-only cookie
ctx.Cookie(&fiber.Cookie{
	Name:     "access_token",
	Value:    accessToken,
	HTTPOnly: true,
	Secure:   true,
	MaxAge:   3600, // 1 hour
	SameSite: "Strict",
})

	return utils.ApiResponse(ctx, fiber.StatusOK, map[string]string{
		"access_token": accessToken,
	}, "Access token refreshed successfully")
}

// GetUserProfile retrieves the profile information of the authenticated user
func GetUserProfile(ctx *fiber.Ctx) error {
	// Get the user information from the context (set by middleware)
	userClaims := ctx.Locals("user").(*utils.Claims)
    println(userClaims)
	// Get the user collection from the database
	collection := db.GetCollection("users")
   println(collection)
	// Find the user by their username or email
	var user models.User
	err := collection.FindOne(ctx.Context(), bson.M{
		"$or": []bson.M{
			{"userName": userClaims.UserName},
			{"email": userClaims.Email}, // To handle if username is email
		},
	}).Decode(&user)

	if err != nil {
		return utils.NewApiError(fiber.StatusNotFound, "User not found").Handle(ctx)
	}

	// Remove the password from the profile response before sending it
	user.Password = ""

	// Return the user profile data
	return utils.ApiResponse(ctx, fiber.StatusOK, user, "User profile fetched successfully")
}
