package models

import "time"

// id: number;
// title: string;
// image_path: string;
// slug: string;
// content: string;
// description: string;
// readonly createdAt?: string;
// readonly updatedAt?: string;
// readonly deletedAt?: string;

type Article struct {
	ID          uint              `gorm:"primary_key;auto_increment" json:"id"`
	Title       string            `gorm:"size:100;unique;not null" json:"title"`
	ImagePath   string            `gorm:"size:255;not null" json:"image_path"`
	Slug        string            `gorm:"size:100;unique;not null" json:"slug"`
	Content     string            `gorm:"not null" json:"content"`
	Description string            `gorm:"size:255;not null" json:"description"`
	CreatedAt   time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Tags        []ArticleTag      `json:"-"`
	Categories  []ArticleCategory `json:"-"`
	Comments    []ArticleComment  `json:"-"`
}

type ArticleTag struct {
	ArticleID uint      `json:"article_id"`
	TagID     uint      `json:"tag_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Article   Article   `json:"-"`
	Tag       Tag       `json:"-"`
}

type ArticleCategory struct {
	ArticleID  uint      `json:"article_id"`
	CategoryID uint      `json:"category_id"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Article    Article   `json:"-"`
	Category   Category  `json:"-"`
}

type ArticleComment struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint      `json:"user_id"`
	ArticleID uint      `json:"article_id"`
	CommentID uint      `gorm:"default:null;" json:"comment_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	User      User      `json:"-"`
	Article   Article   `json:"-"`
}
