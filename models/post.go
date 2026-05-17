package models

import (
	"errors"
	"strings"
	"time"
)

// Post struct that will represent blog posts in our CMS
type Post struct {
	// - ID (unsigned integer, primary key)
	ID uint `json:"id" gorm:"primaryKey"`

	// Title (string, required, with max length)
	Title string `json:"title" gorm:"size:255;not null" binding:"required"`

	// Content (text field, required)
	Content string `json:"content" gorm:"type:text;not null" binding:"required"`

	// Author (string, optional)
	Author string `json:"author" gorm:"size:100"`

	// CreatedAt (timestamp for creation date)
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// UpdatedAt (timestamp for last update)
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Media (slice of Media, representing a many-to-many relationship)
	Media []Media `json:"media" gorm:"many2many:post_media"`
}

func (p *Post) Validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("Post Title cannot be empty")
	}
	if strings.TrimSpace(p.Content) == "" {
		return errors.New("Post Content cannot be empty")
	}
	return nil
}
