package models

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type Article struct {
	ID          uint              `gorm:"primary_key;auto_increment" json:"id"`
	Title       string            `gorm:"size:100;unique;not null" json:"title"`
	ImageUrl    string            `gorm:"size:255;not null" json:"image_url"`
	Slug        string            `gorm:"size:100;unique;not null" json:"slug"`
	Content     string            `gorm:"not null" json:"content"`
	Description string            `gorm:"size:255;not null" json:"description"`
	IsPublished bool              `gorm:"not null" json:"is_published"`
	UserID      uint              `json:"user_id"`
	CreatedAt   time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Tags        []ArticleTag      `gorm:"many2many" json:"tags"`
	Categories  []ArticleCategory `gorm:"many2many" json:"categories"`
	Comments    []ArticleComment  `gorm:"many2many" json:"-"`
	User        User              `json:"author"`
}

func (a *Article) Validate(_ *gorm.DB) []string {
	errs := []string{}

	if a.Title == "" {
		errs = append(errs, "Title artikel harus diisi")
	}

	if a.ImageUrl != "" && !strings.Contains(a.ImageUrl, "http://") && !strings.Contains(a.ImageUrl, "https://") {
		errs = append(errs, "Image url tidak valid")
	}

	return errs
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

func (a *Article) RestoreUpdate(db *gorm.DB, details *Article) {
	fmt.Println("details => ", details)

	category_ids := []uint{}
	for _, c := range details.Categories {
		category_ids = append(category_ids, c.CategoryID)
	}

	tag_ids := []uint{}
	for _, t := range details.Tags {
		tag_ids = append(tag_ids, t.TagID)
	}

	db.Model(&a).Updates(details)
	a.InsertCategories(db, category_ids)
	a.InsertTags(db, tag_ids, []string{})
}

func (a *Article) BeforeUpdate(db *gorm.DB) error {
	var category ArticleCategory
	var tag ArticleTag

	if err := db.Where("article_id=?", a.ID).Delete(&category).Error; err != nil {
		return err
	}

	if err := db.Where("article_id=?", a.ID).Delete(&tag).Error; err != nil {
		return err
	}

	return nil
}

func (a *Article) InsertCategories(db *gorm.DB, ids []uint) []string {
	errs := []string{}

	if len(ids) > 0 && ids[0] != 0 {
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
	}

	return errs
}

func (a *Article) InsertTags(db *gorm.DB, ids []uint, tags []string) []string {
	errs := []string{}

	if len(ids) > 0 && ids[0] != 0 {
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
	}

	fmt.Println("t => ", tags)

	if len(tags) > 0 && tags[0] != "" {
		for _, name := range tags {
			existTag := Tag{}

			tInput := ArticleTag{
				ArticleID: a.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := db.Where("name=?", name).First(&existTag).Error; err != nil {
				tag := &Tag{
					Name:      name,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				if err := db.Create(&tag).Error; err != nil {
					fmt.Println(err.Error())
				}

				tInput.TagID = tag.ID

				if err := db.Create(&tInput).Error; err != nil {
					errs = append(errs, err.Error())
				}
			}

			if !slices.Contains(ids, existTag.ID) {
				tInput.TagID = existTag.ID

				if err := db.Create(&tInput).Error; err != nil {
					errs = append(errs, err.Error())
				}
			}
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
	Content   string    `json:"content"`
	IsReply   bool      `json:"is_reply"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Article   Article   `json:"-"`
	User      User      `json:"user"`
	// Replies   []ReplyArticleComment `json:"replies"`
}

func (co *ArticleComment) GetDetails(db *gorm.DB) error {
	user := User{}

	if err := db.Where("id=?", co.UserID).First(&user).Error; err != nil {
		return err
	}

	co.User = user
	return nil
}

type ReplyArticleComment struct {
	ID        uint           `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint           `json:"user_id"`
	ArticleID uint           `json:"article_id"`
	ParentID  uint           `json:"parent_id"`
	CommentID uint           `gorm:"default:null" json:"comment_id"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Article   Article        `json:"-"`
	Parent    ArticleComment `json:"parent"`
	User      User           `json:"user"`
}

func (r *ReplyArticleComment) GetParent(db *gorm.DB) error {
	var parent ArticleComment
	var user User

	if err := db.Where("id=?", r.ParentID).First(&parent).Error; err != nil {
		return err
	}

	if err := db.Where("id", parent.UserID).First(&user).Error; err != nil {
		return err
	}

	r.Parent = parent
	r.Parent.User = user

	return nil
}

func (co *ReplyArticleComment) GetDetails(db *gorm.DB) error {
	user := User{}

	if err := db.Where("id=?", co.UserID).First(&user).Error; err != nil {
		return err
	}

	co.User = user
	return nil
}
