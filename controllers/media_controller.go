package controllers

import (
	"go-cms-backend/models"
	"go-cms-backend/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMedia(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
	var media []models.Media

	if err := db.Find(&media).Error; err != nil {
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Could not fetch the data"})
		return
	}

	c.JSON(200, media)
}

func GetMediaByID(c *gin.Context) {
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
	if id == 0 {
		c.JSON(404, utils.HTTPError{
			Code: 404,
			Message: "No media with that ID",
		})
		return
	}

	var media models.Media
	if err := db.First(&media, id).Error; err!= nil {
		c.JSON(404, utils.HTTPError{
			Code: 404,
			Message: "No media with that ID",
		})
		return
	}

	c.JSON(200, media)
}

func CreateMedia(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)

	var media models.Media
	if err := c.ShouldBindJSON(&media); err != nil { //Fetched JSON payload
		c.JSON(400, utils.HTTPError{
			Code: 400,
			Message: "Please provide valid data",
		})
		return
	}

	if strings.TrimSpace(media.URL) == "" || strings.TrimSpace(media.Type) == "" {
		c.JSON(400, utils.HTTPError{
			Code: 400,
			Message: "URL/Type cannot be empty",
		})
		return
	}

	if err := db.Create(&media).Error; err != nil {
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to create the data",
		})
		return
	}

	c.JSON(201, utils.MessageResponse{
		Message: "Media created successfully!",
	})
}

func DeleteMedia(c *gin.Context) {
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

	if err := db.Delete(&models.Media{}, id).Error; err != nil {
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to delete media",
		})
		return
	}

	c.JSON(200, utils.MessageResponse{
		Message: "Media deleted successfully!",
	})
}