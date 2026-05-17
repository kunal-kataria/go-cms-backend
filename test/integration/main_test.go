package integration

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/kunal-kataria/go-cms-backend/models"
	"github.com/kunal-kataria/go-cms-backend/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB        // Test database connection
var testRouter *gin.Engine // Gin router instance

func TestMain(m *testing.M) {
	// Environment Setup
	if err := setup(); err != nil {
		log.Fatalf("Failed to set up test environment: %v", err)
	}

	// Execute all integration tests
	exitCode := m.Run()

	// Cleanup & Exit with test result code
	cleanup()
	os.Exit(exitCode)
}

func setup() error {
	// Gin to TestMode
	gin.SetMode(gin.TestMode)

	// Database Connection
	dbUser := os.Getenv("TEST_DB_USER")
	dbPassword := os.Getenv("TEST_DB_PASSWORD")
	dbName := os.Getenv("TEST_DB_NAME")
	dbHost := os.Getenv("TEST_DB_HOST")
	dbPort := os.Getenv("TEST_DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort)
	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Schema Migration
	if err := testDB.AutoMigrate(&models.Page{}, &models.Post{}, &models.Media{}); err != nil {
		return fmt.Errorf("failed to automigrate database: %w", err)
	}

	// Router Setup
	testRouter = gin.Default()
	routes.InitializeRoutes(testRouter, testDB)

	return nil

}

func cleanup() {
	sqlDB, err := testDB.DB()
	if err != nil {
		log.Fatalf("Failed to get Database instance: %v", err)
	}

	// Database Cleanup
	if err := testDB.Exec("DROP TABLE IF EXISTS post_media").Error; err != nil {
		log.Printf("Failed to drop table post_media: %v", err)
	}
	if err := testDB.Exec("DROP TABLE IF EXISTS posts").Error; err != nil {
		log.Printf("Failed to drop table posts: %v", err)
	}
	if err := testDB.Exec("DROP TABLE IF EXISTS media").Error; err != nil {
		log.Printf("Failed to drop table media: %v", err)
	}
	if err := testDB.Exec("DROP TABLE IF EXISTS pages").Error; err != nil {
		log.Printf("Failed to drop table pages: %v", err)
	}

	// Connection Cleanup
	if err := sqlDB.Close(); err != nil {
		log.Printf("Failed to close database: %v", err)
	}
}

func clearTables() {
	// Data Cleanup
	if err := testDB.Exec("DELETE FROM post_media").Error; err != nil {
		log.Printf("Failed to clear post_media: %v", err)
	}
	if err := testDB.Exec("DELETE FROM posts").Error; err != nil {
		log.Printf("Failed to clear posts: %v", err)
	}
	if err := testDB.Exec("DELETE FROM media").Error; err != nil {
		log.Printf("Failed to clear media: %v", err)
	}
	if err := testDB.Exec("DELETE FROM pages").Error; err != nil {
		log.Printf("Failed to clear pages: %v", err)
	}
}
