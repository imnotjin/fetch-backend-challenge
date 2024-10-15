package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/imnotjin/fetch-backend-challenge/docs"
	"github.com/imnotjin/fetch-backend-challenge/handlers"
	"github.com/imnotjin/fetch-backend-challenge/models"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// main is the entry point for the Fetch Rewards API server. It sets up the Swagger documentation,
// initializes the production database, runs necessary migrations, and starts the Gin HTTP server.
func main() {
	// Setup Swagger configuration
	docs.SwaggerInfo.Title = "Fetch Rewards API"
	docs.SwaggerInfo.Description = "This is a points management API server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8000"
	docs.SwaggerInfo.BasePath = "/"

	// Setup production database
	prodDB, err := setupProdDB()
	if err != nil {
		log.Fatalf("Failed to connect to production database: %v", err)
	}

	// Migrate the schema
	err = prodDB.AutoMigrate(&models.Transaction{})
	if err != nil {
		fmt.Printf("Failed to run migrations on production database: %v\n", err)
		os.Exit(1)
	}

	// Setup and start the HTTP router
	r := setupProdRouter(prodDB)
	r.Run(":8000")
}

// setupProdDB initializes a connection to the production PostgreSQL database using the
// connection details from the .env file. It returns a pointer to the gorm.DB instance.
func setupProdDB() (*gorm.DB, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Build DSN (Data Source Name) string for the production database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Open the database connection using gorm
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// setupProdRouter initializes the Gin router and registers the API routes for managing points.
// It configures routes for adding points, spending points, and fetching the account balance.
// The Swagger API documentation is also served at /swagger/*any.
func setupProdRouter(prodDB *gorm.DB) *gin.Engine {
	r := gin.Default()
	h := &handlers.Handler{DB: prodDB}

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Points management endpoints
	r.POST("/add", h.AddPoints)
	r.POST("/spend", h.SpendPoints)
	r.GET("/balance", h.GetBalance)

	return r
}
