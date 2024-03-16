package controller

import (
	"api/app/lib"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/spf13/viper"
)

func TestPostLogsProduction(t *testing.T) {
	app := fiber.New()
	app.Post("/logs", PostLogs)

	// Set up a mock environment
	viper.Set("ENVIRONTMENT", "production")

	// Create a test logger file with some logs
	loggerFileContent := []byte("Log line 1\nLog line 2\n")
	loggerFilePath := "logger.log"
	err := ioutil.WriteFile(loggerFilePath, loggerFileContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create logger file: %v", err)
	}
	defer os.Remove(loggerFilePath)

	// Make a request to the PostLogs handler
	response, _, err := lib.PostTest(app, "/logs", nil)
	if err != nil {
		t.Fatalf("PostLogs request failed: %v", err)
	}

	// Verify the response
	utils.AssertEqual(t, 200, response.StatusCode, "HTTP Status")

	// Clean up the mock environment
	viper.Set("ENVIRONTMENT", "")
}

func TestPostLogsDevelopment(t *testing.T) {
	app := fiber.New()
	app.Post("/logs", PostLogs)

	// Set up a mock environment
	viper.Set("ENVIRONTMENT", "development")

	// Create a test logger file with some logs
	loggerFileContent := []byte("Custom log line 1\nCustom log line 2\n")
	loggerFilePath := "logs/logger.log"
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}
	err = ioutil.WriteFile(loggerFilePath, loggerFileContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create custom logger file: %v", err)
	}
	defer os.Remove(loggerFilePath)
	defer os.RemoveAll("logs")

	// Make a request to the PostLogs handler
	response, _, err := lib.PostTest(app, "/logs", nil)
	if err != nil {
		t.Fatalf("PostLogs request failed: %v", err)
	}

	// Verify the response
	utils.AssertEqual(t, 200, response.StatusCode, "HTTP Status")

	// Clean up the mock environment
	viper.Set("ENVIRONTMENT", "")
}

func TestPostLogsErrorProduction(t *testing.T) {
	app := fiber.New()
	app.Post("/logs", PostLogs)

	// Set up a mock environment
	viper.Set("ENVIRONTMENT", "production")

	// Make a request to the PostLogs handler
	response, _, err := lib.PostTest(app, "/logs", nil)
	if err != nil {
		t.Fatalf("PostLogs request failed: %v", err)
	}

	// Verify the response
	utils.AssertEqual(t, 500, response.StatusCode, "HTTP Status")

	// Clean up the mock environment
	viper.Set("ENVIRONTMENT", "")
}

func TestPostLogsErrorDevelopment(t *testing.T) {
	app := fiber.New()
	app.Post("/logs", PostLogs)

	// Set up a mock environment
	viper.Set("ENVIRONTMENT", "development")

	// Make a request to the PostLogs handler
	response, _, err := lib.PostTest(app, "/logs", nil)
	if err != nil {
		t.Fatalf("PostLogs request failed: %v", err)
	}

	// Verify the response
	utils.AssertEqual(t, 500, response.StatusCode, "HTTP Status")

	// Clean up the mock environment
	viper.Set("ENVIRONTMENT", "")
}
