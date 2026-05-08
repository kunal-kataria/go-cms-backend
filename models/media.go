package models

import "time"

type Media struct {
	//ID field as uint with gorm tag for primary key and json tag
	ID uint `json:"id" gorm:"primaryKey"`

	//URL field as string with gorm tag for size limit (255) and not null constraint and json tag and binding tag to make it required
	URL string `json:"url" gorm:"size:255;not null" binding:"required"`

	//Type field as   string with gorm tag for size limit (50) and json tag and binding tag to make it required
	Type string `json:"type" gorm:"size:50" binding:"required"`

	//CreatedAt field as time.Time with gorm tag for automatic timestamp on creation and json tag
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	//UpdatedAt field as time.Time with gorm tag for automatic timestamp on updates and json tag
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
