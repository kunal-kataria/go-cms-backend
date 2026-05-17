package controllers

import (
	"github.com/kunal-kataria/go-cms-backend/models"
	"github.com/kunal-kataria/go-cms-backend/utils"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetPages(t *testing.T) {
	gin.SetMode(gin.TestMode) // set gin to Test mode

	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "First Page", "Content 1", time.Now(), time.Now()).
		AddRow(2, "Second Page", "Content 2", time.Now(), time.Now())
	mock.ExpectQuery(`SELECT \* FROM "pages"`).WillReturnRows(rows)

	router.GET("/pages", GetPages)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/pages", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response []models.Page
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error Unmarshaling response: %v", err)
	}
	if len(response) != 2 {
		t.Fatalf("Expected 2 pages, but got %d", len(response))
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPage(t *testing.T) {
	// Add test for GetPage
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock Data Creation
	row := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "First Page", "Content 1", time.Now(), time.Now())

	// Database Expectations
	mock.ExpectQuery(`SELECT .+ FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(row)

	// HTTP Test Setup
	router.GET("/pages/:id", GetPage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/pages/1", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected code 200, but got %d", w.Code)
	}

	// Response Validation
	var response models.Page
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error Unmarshaling response: %v", err)
	}

	if response.ID != 1 {
		t.Fatalf("Expected ID = 1, got ID = %d", response.ID)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	page := models.Page{
		Title:   "ABC",
		Content: "Reels",
	}

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(`INSERT INTO "pages"`).
		WithArgs(page.Title, page.Content, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)
	mock.ExpectCommit()

	body, err := json.Marshal(page)
	if err != nil {
		t.Fatalf("Failed to marshal page: %v", err)
	}

	router.POST("/pages", CreatePage)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/pages", bytes.NewBuffer(body))
	req.Header.Set("Content-type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response utils.MessageResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, "Page created!", response.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePage(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Database Expectations
	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "Old Title", "Old Content", time.Now(), time.Now())
	mock.ExpectQuery(`SELECT .* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "pages" SET .* WHERE id = \$4`).
		WithArgs("Updated Title", "Updated Content", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Request Preparation
	update := models.Page{
		Title:   "Updated Title",
		Content: "Updated Content",
	}

	body, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Failed to marshal update: %v", err)
	}

	// HTTP Test Setup
	router.PUT("/pages/:id", UpdatePage)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/pages/1", bytes.NewBuffer(body))
	req.Header.Set("Content-type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Response Validation
	var response utils.MessageResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "Page updated!", response.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeletePage(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Database Expectations
	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "Title", "Content", time.Now(), time.Now())
	mock.ExpectQuery(`SELECT .* FROM "pages" WHERE "pages"."id" = \$1 ORDER BY "pages"."id" LIMIT \$2`).
		WithArgs(1, 1).WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "pages" WHERE "pages"."id" = \$1`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	//  HTTP Test Setup
	router.DELETE("/pages/:id", DeletePage)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/pages/1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Response Validation
	var response utils.MessageResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "Page deleted!", response.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/*
TESTING HINTS:
1. Use sqlmock.AnyArg() for timestamp fields
2. Remember to escape special characters in SQL patterns
3. Each database operation needs proper error handling
4. Content-Type header is required for POST/PUT requests
5. Transaction tests need Begin/Commit expectations
6. Use proper argument matching in mock expectations
7. Consider testing error cases:
   - Invalid IDs
   - Missing required fields
   - Database errors
*/
