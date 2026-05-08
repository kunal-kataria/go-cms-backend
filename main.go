package main

import (
	"go-cms-backend/models"
	"go-cms-backend/routes"
	"go-cms-backend/utils"
	"os"

	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Initialize database connection
	db, err := utils.DbConnect()
	if err != nil {
		log.Fatalf("Failed to connect to Database: %v", err)
	}
	
	// Get the underlying *sql.DB instance and defer its closure
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get Database instance: %v", err)
	}
	defer sqlDB.Close()

	// Get the environment variable
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	// Conditionally run AutoMigrate in development environment
	if env == "development" {
		log.Println("Running Automigrator")
		if err := db.AutoMigrate(&models.Page{}, &models.Post{}, &models.Media{}); err != nil {
			log.Fatalf("Failed to automigrate database: %v", err)
		}
	}

	// Set Gin mode based on environment
	if env == "production" {
		gin.SetMode(gin.ReleaseMode) // ReleaseMode reduces log noise and disables debug features. 
	}

	router := gin.Default()

	// Initialize routes
	routes.InitializeRoutes(router, db)

	// Run the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
