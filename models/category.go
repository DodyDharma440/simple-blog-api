package models

import "time"

type Category struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:100;not null;unique" json:"name"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	// Article   []Article `json:"-"`
}