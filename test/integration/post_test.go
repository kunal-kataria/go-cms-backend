package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/kunal-kataria/go-cms-backend/models"
	"github.com/kunal-kataria/go-cms-backend/utils"

	"strings"
	"testing"
)

func TestPostIntegration(t *testing.T) {
	// Clear Database
	clearTables()

	// Create Test Media
	mediaID := createTestMedia(t)

	t.Run("Create Post with Media", func(t *testing.T) {
		// TODO: Create Post with Media
		body := fmt.Sprintf(`{
			"title": "Post Title",
			"content": "Post Content",
			"author": "Xyz",
			"media": [{"id": %d}]
		}`, mediaID)

		req := httptest.NewRequest("POST", "/api/v1/posts/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status 201, got %d: %s", w.Code, w.Body.String())
		}

		var response utils.MessageResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response.Message != "Post created!" {
			t.Fatalf("Expected message 'Post created!', got %s", response.Message)
		}
	})
}

// Helper function to create test media
func createTestMedia(t *testing.T) uint {
	body := `{
		"url": "www.abc.com/imge.jpg",
		"type": "image"
	}`

	req := httptest.NewRequest("POST", "/api/v1/media/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create test media, status: %d, body: %s", w.Code, w.Body.String())
	}

	var response models.Media
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to create test media: %v", err)
	}

	return response.ID
}
