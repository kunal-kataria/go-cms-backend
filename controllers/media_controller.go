package controllers

import (
	"go-cms-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMedia(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
	var media []models.Media



}

func GetMediaByID(c *gin.Context) {
    // Implement logic to retrieve media by ID
}

func CreateMedia(c *gin.Context) {
    // Implement logic to create new media
}

func DeleteMedia(c *gin.Context) {
    // Implement logic to delete media
}