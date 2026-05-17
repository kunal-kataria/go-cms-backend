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

func TestGetPosts(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock Data Creation
	postRows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(1, "Title", "post", "Xyz", time.Now(), time.Now()).
		AddRow(2, "Title2", "post2", "XyzAbc", time.Now(), time.Now())
	postMediaRows := sqlmock.NewRows([]string{"post_id", "media_id"}).
		AddRow(1, 1).
		AddRow(2, 2)
	mediaRows := sqlmock.NewRows([]string{"id", "url", "type", "created_at", "updated_at"}).
		AddRow(1, "abc.com", "photo", time.Now(), time.Now()).
		AddRow(2, "def.com", "video", time.Now(), time.Now())

	// Database Expectations
	mock.ExpectQuery(`SELECT .* FROM "posts"`).WillReturnRows(postRows)
	mock.ExpectQuery(`SELECT .* FROM "post_media" WHERE "post_media"\."post_id" IN \(\$1,\$2\)`).
		WithArgs(1, 2).
		WillReturnRows(postMediaRows)
	mock.ExpectQuery(`SELECT .* FROM "media" WHERE "media"\."id" IN \(\$1,\$2\)`).
		WithArgs(1, 2).
		WillReturnRows(mediaRows)

	// HTTP Test Setup
	router.GET("/posts", GetPosts)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	// Response Validation
	var response []models.Post
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to Unmarshal json: %v", err)
	}

	if len(response) != 2 {
		t.Fatalf("Expected 2 records but got %d", len(response))
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPost(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock Data Creation
	post := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(1, "Title", "post", "Xyz", time.Now(), time.Now())
	postMedia := sqlmock.NewRows([]string{"post_id", "media_id"}).
		AddRow(1, 1)
	media := sqlmock.NewRows([]string{"id", "url", "type", "created_at", "updated_at"}).
		AddRow(1, "abc.com", "photo", time.Now(), time.Now())

	// Database Expectations
	mock.ExpectQuery(`SELECT .+ FROM "posts" WHERE "posts"."id" = \$1 ORDER BY "posts"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(post)
	mock.ExpectQuery(`SELECT .+ FROM "post_media" WHERE "post_media"."post_id" = \$1`).
		WithArgs(1).
		WillReturnRows(postMedia)
	mock.ExpectQuery(`SELECT .+ FROM "media" WHERE "media"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(media)

	// HTTP Test Setup
	router.GET("/posts/:id", GetPost)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/posts/1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Response Validation
	var response models.Post
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshall json: %v", err)
	}

	if response.ID != 1 {
		t.Fatalf("Expected ID=1, but got ID= %d", response.ID)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePost(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock Data Creation
	postBody := models.Post{
		Title:   "ABC",
		Content: "Photo",
		Author:  "Xyz",
	}

	// Database Expectations
	post := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "posts"`).
		WithArgs(postBody.Title, postBody.Content, postBody.Author, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(post)
	mock.ExpectCommit()

	// HTTP Test Setup
	body, err := json.Marshal(postBody)
	if err != nil {
		t.Fatalf("Failed to Marshal post: %v", err)
	}
	router.POST("/posts", CreatePost)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Response Validation
	var response utils.MessageResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "Post created!", response.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePost(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock Data Creation
	postBody := models.Post{
		Title:   "ABC",
		Content: "Photo",
		Author:  "Xyz",
	}
	post := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(1, "Title", "post", "Xyza", time.Now(), time.Now())

	// Database Expectations
	mock.ExpectQuery(`SELECT .* FROM "posts" WHERE "posts"."id" = \$1 ORDER BY "posts"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(post)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "posts" SET .* WHERE id = \$5`).
		WithArgs("ABC", "Photo", "Xyz", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// HTTP Test Setup
	body, err := json.Marshal(postBody)
	if err != nil {
		t.Fatalf("Failed to Marshal post: %v", err)
	}
	router.PUT("/posts/:id", UpdatePost)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/posts/1", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Response Validation
	var response utils.MessageResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "Post updated!", response.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeletePost(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	post := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(1, "Title", "post", "Xyza", time.Now(), time.Now())

	// Database Expectations
	mock.ExpectQuery(`SELECT .* FROM "posts" WHERE "posts"."id" = \$1 ORDER BY "posts"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(post)
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "posts" WHERE "posts"."id" = \$1`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// HTTP Test Setup
	router.DELETE("/posts/:id", DeletePost)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/posts/1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Response Validation
	var response utils.MessageResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "Post deleted!", response.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}
