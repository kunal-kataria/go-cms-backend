package controllers

import (
	"go-cms-backend/models"
	"go-cms-backend/utils"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPages(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var pages []models.Page

	if err := db.Find(&pages).Error; err != nil {
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to fetch data",
		})
		return
	}

	c.JSON(200, pages)
}

func GetPage(c *gin.Context) {
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

	var page models.Page
	if err := db.First(&page, id).Error; err != nil {
		c.JSON(404, utils.HTTPError{
			Code: 404,
			Message: "Record not found",
		})
		return
	}

	c.JSON(200, page)
}

func CreatePage(c *gin.Context) {

}

func UpdatePage(c *gin.Context) {

}

func DeletePage(c *gin.Context) {
	
}