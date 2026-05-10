package controllers

import (
	"go-cms-backend/models"
	"go-cms-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPosts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var posts []models.Post

	title := c.Query("title")
	author := c.Query("author")

	query := db
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if author != "" {
		query = query.Where("author = ?", author)
	}

	if err := query.Preload("Media").Find(&posts).Error; err != nil {
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to fetch posts",
		})
		return
	}

	c.JSON(200, posts)
}

func GetPost(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	strid := c.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		c.JSON(400, utils.HTTPError{
			Code: 400,
			Message: "Invalid Id",
		})
		return
	}

	title := c.Query("title")
	author := c.Query("author")

	query := db
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if author != "" {
		query = query.Where("author = ?", author)
	}

	var post models.Post
	if err := query.Preload("Media").First(&post, id).Error; err != nil {
		c.JSON(404, utils.HTTPError{
			Code: 404,
			Message: "Post not found",
		})
		return
	}

	c.JSON(200, post)
}

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, utils.HTTPError{
			Code: 400,
			Message: "Invalid post data",
		})
		return
	}

	if err := post.Validate(); err != nil {
		c.JSON(400, utils.HTTPError{
			Code: 400,
			Message: err.Error(),
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to start transaction",
		})
		return
	}

	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to create Post",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to commit transaction",
		})
		return
	}

	c.JSON(201, utils.MessageResponse{
		Message: "Post created!",
	})
}

func UpdatePost(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	strid := c.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		c.JSON(400, utils.HTTPError{
			Code: 400,
			Message: "Invalid ID",
		})
		return
	}
	
	var existingPost models.Post
	if err := db.First(&existingPost, id).Error; err != nil {
		c.JSON(404, utils.HTTPError{
			Code: 404,
			Message: "Post not found!",		
		})
		return
	}

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, utils.HTTPError{
			Code: 400,
			Message: "Invalid Post data",
		})
		return
	}

	if err := post.Validate(); err != nil {
		c.JSON(400, utils.HTTPError{
			Code: 400,
			Message: err.Error(),
		})
		return
	}

	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to start transaction",
		})
		return
	}

	if err := tx.Model(&models.Post{}).Where("id = ?", id).Updates(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to update post",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to commit transaction",
		})
		return
	}

	c.JSON(200, utils.MessageResponse{
		Message: "Post updated!",
	})
}

func DeletePost(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	strid := c.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		c.JSON(400, utils.HTTPError{
			Code: 400,
			Message: "Invalid ID",
		})
		return
	}
	
	var post models.Post
	if err := db.First(&post, id).Error; err != nil {
		c.JSON(404, utils.HTTPError{
			Code: 404,
			Message: "Post not found!",		
		})
		return
	}

	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to start transaction",
		})
		return
	}

	if  err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to delete Post",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to commit transaction",
		})
		return
	}

	c.JSON(200, utils.MessageResponse{
		Message: "Post deleted!",
	})
}