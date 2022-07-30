package controllers

import (
	"final-project/models"
	"final-project/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleInput struct {
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	Image       interface{} `json:"image"`
	Description string      `json:"description"`
	Tags        []uint      `json:"tag_ids"`
	Categories  []uint      `json:"category_ids"`
}

// Get All Articles godoc
// @Summary     Get all articles.
// @Tags        Article
// @Produce     json
// @Success     200 {object} []models.Article
// @Router      /articles [get]
func GetArticles(c *gin.Context) {
	var articles []models.Article

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Find(&articles).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(articles) > 0 {
		_articles := []models.Article{}

		for _, article := range articles {
			article.GetDetails(db)
			_articles = append(_articles, article)
		}

		articles = _articles
	}

	utils.CreateResponse(c, http.StatusOK, &articles)
}

// Get Article godoc
// @Summary     Get article by id.
// @Tags        Article
// @Produce     json
// @Param id path string true "article id"
// @Success     200 {object} models.Article
// @Router      /articles/{id} [get]
func GetArticle(c *gin.Context) {
	var article models.Article

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id=?", c.Param("id")).First(&article).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	article.GetDetails(db)
	utils.CreateResponse(c, http.StatusOK, &article)
}

// Get Article godoc
// @Summary     Get article by slug.
// @Tags        Article
// @Produce     json
// @Param slug path string true "slug"
// @Success     200 {object} models.Article
// @Router      /articles/slug/{slug} [get]
func GetArticleBySlug(c *gin.Context) {
	var article models.Article

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("slug=?", c.Param("slug")).First(&article).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	utils.CreateResponse(c, http.StatusOK, article)
}

func CreateArticle(c *gin.Context) {
	var errs = []string{}

	filepath, err := utils.UploadFile(c, "articles")

	if err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	article := models.Article{
		Title:       c.PostForm("title"),
		Content:     c.PostForm("content"),
		ImagePath:   filepath,
		Description: c.PostForm("description"),
		IsPublished: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	userID, err := utils.ExtractTokenID(c)

	if err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	article.GetSlug(db)
	article.UserID = userID

	if err := db.Create(&article).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	categories := strings.Split(c.PostForm("category_ids"), ",")
	tags := strings.Split(c.PostForm("tag_ids"), ",")

	category_ids := utils.SliceStringToUInt(categories)
	tag_ids := utils.SliceStringToUInt(tags)

	errs = append(errs, article.InsertCategories(db, category_ids)...)
	errs = append(errs, article.InsertTags(db, tag_ids)...)

	if len(errs) > 0 {
		if err := article.Delete(db); err != nil {
			fmt.Println(err.Error())
		}
		utils.CreateResponse(c, http.StatusBadRequest, errs)
		return
	}

	utils.CreateResponse(c, http.StatusCreated, article)
}

func UpdateArticle(c *gin.Context) {}

// Delete Article godoc
// @Summary     Delete article.
// @Tags        Article
// @Produce     json
// @Param id path string true "article id"
// @Success     200 {object} bool
// @Router      /articles/{id} [delete]
// @Security ApiKeyAuth
func DeleteArticle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var article models.Article

	if err := db.Where("id=?", c.Param("id")).First(&article).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	filePathSlice := strings.Split(article.ImagePath, "/")
	fileName := filePathSlice[len(filePathSlice)-1]

	err := os.Remove("public/upload/articles/" + fileName)

	if err != nil {
		log.Fatal(err)
	}

	if err := article.Delete(db); err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, true)
}

// Publish Article godoc
// @Summary     Publish Article.
// @Tags        Article
// @Produce     json
// @Param 		id path string true "article id"
// @Success     200 {object} models.User
// @Router      /articles/publish/{id} [patch]
// @Security ApiKeyAuth
func PublishArticle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var article models.Article

	if err := db.Where("id=?", c.Param("id")).First(&article).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	updated := models.Article{
		IsPublished: true,
		UpdatedAt:   time.Now(),
	}

	if err := db.Model(&article).Updates(updated).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, &article)
}

// Unpublish Article godoc
// @Summary     Unpublish Article.
// @Tags        Article
// @Produce     json
// @Param 		id path string true "article id"
// @Success     200 {object} models.User
// @Router      /articles/unpublish/{id} [patch]
// @Security ApiKeyAuth
func UnpublishArticle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var article models.Article

	if err := db.Where("id=?", c.Param("id")).First(&article).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	updated := models.Article{
		IsPublished: false,
		UpdatedAt:   time.Now(),
	}

	if err := db.Model(&article).Updates(updated).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, &article)
}
