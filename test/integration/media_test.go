package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kunal-kataria/go-cms-backend/models"
	"github.com/kunal-kataria/go-cms-backend/utils"
)

func TestMediaIntegration(t *testing.T) {
	// Clear Database
	clearTables()

	t.Run("Create Media", func(t *testing.T) {
		// Test Data
		body := `{
			"url": "http://abc.com/test.jpg",
			"type": "image"
		}`

		// HTTP Request
		req := httptest.NewRequest("POST", "/api/v1/media/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Execute Request
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status 201, got %d: %s", w.Code, w.Body.String())
		}

		// Verify Response
		var response utils.MessageResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response.Message != "Media created successfully!" {
			t.Fatalf("Expected message 'Media created successfully!', got %s", response.Message)
		}
	})

	t.Run("Get All Media", func(t *testing.T) {
		// STEP 1: Create & Execute HTTP Request
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/media/", nil)
		testRouter.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		// Verify Response
		var response []models.Media
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response) != 1 {
			t.Errorf("Expected 1 record, got %d", len(response))
		}

	})

}
