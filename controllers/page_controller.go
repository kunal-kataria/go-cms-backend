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
	db := c.MustGet("db").(*gorm.DB)

	var page models.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(400, utils.HTTPError{
			Code: 400,
			Message: "Invalid Page data",
		})
		return
	}

	if err := page.Validate(); err != nil {
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

	if err := tx.Create(&page).Error; err != nil {
		tx.Rollback()
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to create Page",
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
		Message: "Page created!",
	})
}

func UpdatePage(c *gin.Context) {
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

	var existingPage models.Page
	if err := db.First(&existingPage, id).Error; err != nil {
		c.JSON(404, utils.HTTPError{
			Code: 404,
			Message: "Page not found!",
		})
		return
	}

	var page models.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(400, utils.HTTPError{
			Code: 400,
			Message: "Invalid Page data",
		})
		return
	}

	if err := page.Validate(); err != nil {
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

	if err := tx.Model(&models.Page{}).Where("id = ?", id).Updates(&page).Error; err != nil {
		tx.Rollback()
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to update Page",
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
		Message: "Page updated!",
	})

}

func DeletePage(c *gin.Context) {
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
			Message: "Page not found!",
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

	if err := tx.Delete(&page).Error; err != nil {
		tx.Rollback()
		c.JSON(500, utils.HTTPError{
			Code: 500,
			Message: "Failed to delete page",
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
		Message: "Page deleted!",
	})
}