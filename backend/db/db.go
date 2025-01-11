package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson" // Import bson for MongoDB operations
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var DB *mongo.Database

// ConnectDB connects to MongoDB using the URI from environment variables
func ConnectDB() {
	log.Println("Starting DB connection...")

	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatalf("MONGODB_URI is not set in .env file")
	}
	log.Println("Mongo URI Loaded")

	// Set the database name
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatalf("DB_NAME is not set in .env file")
	}
	log.Println("Database name loaded:", dbName)

	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Set a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check if the connection is successful
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	} else {
		log.Println("Successfully connected to MongoDB")
	}

	// Assign the database and client instances
	mongoClient = client
	DB = client.Database(dbName)

	// Log the connected database
	log.Printf("Using database: %s", dbName)

	// Verify the collections in the database
	collections, err := DB.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	} else {
		log.Println("Collections in database:", collections)
	}

	log.Println("DB CONNECTED")
}

// DisconnectDB gracefully disconnects the MongoDB client
func DisconnectDB() {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := mongoClient.Disconnect(ctx)
		if err != nil {
			log.Fatalf("Error disconnecting MongoDB client: %v", err)
		}
		log.Println("Disconnected from MongoDB")
	}
}

// GetCollection returns a handle to the specified MongoDB collection
func GetCollection(name string) *mongo.Collection {
	// Log when trying to retrieve the collection
	log.Printf("Fetching collection: %s", name)

	// Ensure DB is initialized before returning the collection
	if DB == nil {
		log.Fatal("Database not connected!")
	}

	collection := DB.Collection(name)
	if collection == nil {
		log.Fatalf("Failed to get collection: %s", name)
	} else {
		log.Printf("Successfully fetched collection: %s", name)
	}

	return collection
}
