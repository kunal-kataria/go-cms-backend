package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kunal-kataria/go-cms-backend/models"
	"github.com/kunal-kataria/go-cms-backend/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetMedia(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"id", "url", "type", "created_at", "updated_at"}).
		AddRow(1, "abc.com", "photo", time.Now(), time.Now()).
		AddRow(2, "def.com", "video", time.Now(), time.Now())
	mock.ExpectQuery(`SELECT .* FROM "media"`).WillReturnRows(rows)

	router.GET("/media", GetMedia)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/media", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response []models.Media
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}
	if len(response) != 2 {
		t.Fatalf("Expected length 2 but got %d", len(response))
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMediaByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"id", "url", "type", "created_at", "updated_at"}).
		AddRow(1, "abc.com", "photo", time.Now(), time.Now())
	mock.ExpectQuery(`SELECT .+ FROM "media" WHERE "media"."id" = \$1 ORDER BY "media"."id" LIMIT \$2`).
		WithArgs(1, 1).WillReturnRows(rows)

	router.GET("/media/:id", GetMediaByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/media/1", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response models.Media
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshalling json: %v", err)
	}

	if response.ID != 1 {
		t.Fatalf("Expected ID=1, but got ID=%d", response.ID)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateMedia(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	media := models.Media{
		URL:  "abc.com",
		Type: "photo",
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "media"`).
		WithArgs(media.URL, media.Type, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)
	mock.ExpectCommit()

	body, err := json.Marshal(media)
	if err != nil {
		t.Fatalf("Failed to marshal media: %v", err)
	}

	router.POST("/media", CreateMedia)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/media", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, but got %d", w.Code)
	}

	var response utils.MessageResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshalling json: %v", err)
	}
	assert.Equal(t, "Media created successfully!", response.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteMedia(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "media" WHERE "media"."id" = \$1`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router.DELETE("/media/:id", DeleteMedia)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/media/1", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response utils.MessageResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshalling json: %v", err)
	}
	assert.Equal(t, "Media deleted successfully!", response.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}
