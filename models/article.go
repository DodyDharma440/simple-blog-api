package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

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
	IsPublished bool              `gorm:"not null" json:"is_published"`
	UserID      uint              `json:"user_id"`
	CreatedAt   time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Tags        []ArticleTag      `json:"tags"`
	Categories  []ArticleCategory `json:"categories"`
	Comments    []ArticleComment  `json:"-"`
	User        User              `json:"author"`
}

func (a *Article) GetSlug(_ *gorm.DB) {
	slugSlice := strings.Split(a.Title, " ")
	slug := strings.ToLower(strings.Join(slugSlice, "-"))

	a.Slug = slug
}

func (a *Article) Delete(db *gorm.DB) error {
	var article Article
	var category ArticleCategory
	var tag ArticleTag

	if err := db.Where("id=?", a.ID).First(&article).Error; err != nil {
		return err
	}

	if err := db.Where("article_id=?", a.ID).Delete(&category).Error; err != nil {
		return err
	}

	if err := db.Where("article_id=?", a.ID).Delete(&tag).Error; err != nil {
		return err
	}

	if err := db.Delete(&article).Error; err != nil {
		return err
	}

	return nil
}

func (a *Article) InsertCategories(db *gorm.DB, ids []uint) []string {
	errs := []string{}

	for _, id := range ids {
		ca := Category{}

		if err := db.Where("id=?", id).First(&ca).Error; err != nil {
			errs = append(errs, err.Error())
		}

		caInput := ArticleCategory{
			ArticleID:  a.ID,
			CategoryID: uint(id),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := db.Create(&caInput).Error; err != nil {
			errs = append(errs, err.Error())
		}
	}

	return errs
}

func (a *Article) InsertTags(db *gorm.DB, ids []uint) []string {
	errs := []string{}

	for _, id := range ids {
		ca := Tag{}

		if err := db.Where("id=?", id).First(&ca).Error; err != nil {
			errs = append(errs, err.Error())
		}

		tInput := ArticleTag{
			ArticleID: a.ID,
			TagID:     uint(id),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&tInput).Error; err != nil {
			errs = append(errs, err.Error())
		}
	}

	return errs
}

func (a *Article) GetDetails(db *gorm.DB) {
	categories := []ArticleCategory{}
	tags := []ArticleTag{}

	db.Where("article_id=?", a.ID).Find(&a.Categories)
	db.Where("article_id=?", a.ID).Find(&a.Tags)
	db.Where("id=?", a.UserID).First(&a.User)

	for _, category := range a.Categories {
		db.Where("id = ?", category.CategoryID).First(&category.Category)
		categories = append(categories, category)
	}

	for _, tag := range a.Tags {
		db.Where("id = ?", tag.TagID).First(&tag.Tag)
		tags = append(tags, tag)
	}

	a.Categories = categories
	a.Tags = tags
}

type ArticleTag struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`
	ArticleID uint      `json:"article_id"`
	TagID     uint      `json:"tag_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Article   Article   `json:"-"`
	Tag       Tag       `json:"tag"`
}

type ArticleCategory struct {
	ID         uint      `gorm:"primary_key;auto_increment" json:"id"`
	ArticleID  uint      `json:"article_id"`
	CategoryID uint      `json:"category_id"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Article    Article   `json:"-"`
	Category   Category  `json:"category"`
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
