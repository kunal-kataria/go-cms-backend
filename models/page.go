package models

import (
	"errors"
	"time"
	"strings"
)

// Page struct that will represent pages in our CMS
type Page struct {
	// - ID (unsigned integer, primary key)
	ID uint `json:"id" gorm:"primaryKey"`

	// Title (string, required, with max length)
	Title string `json:"title" gorm:"size:255;not null" binding:"required"`

	// Content (text field, required)
	Content string `json:"content" gorm:"type:text;not null" binding:"required"`
    
	// CreatedAt (timestamp for creation date)
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// UpdatedAt (timestamp for last update)
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (p *Page) Validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("Page Title cannot be empty")
	}
	if strings.TrimSpace(p.Content) == "" {
		return errors.New("Page Content cannot be empty")
	}
	return nil
}