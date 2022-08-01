package controllers

import (
	"final-project/models"
	"final-project/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleInput struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	Tags        string `json:"tag_ids"`
	TagsNew     string `json:"tags"`
	Categories  string `json:"category_ids"`
}

type CommentInput struct {
	Content string `json:"content"`
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

	if err := db.Order("created_at DESC").Find(&articles).Error; err != nil {
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

// Get Article godoc
// @Summary     Get article by tag name.
// @Tags        Article
// @Produce     json
// @Param tag path string true "tag name"
// @Success     200 {object} []models.Article
// @Router      /articles/tag/{tag} [get]
func GetArticleByTag(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var tag models.Tag

	if err := db.Where("name=?", c.Param("tag")).First(&tag).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	var articles []models.Article
	var tags []models.ArticleTag

	if err := db.Where("tag_id=?", tag.ID).Find(&tags).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	for _, tag := range tags {
		article := models.Article{}
		err := db.Where("id=?", tag.ArticleID).First(&article).Error

		article.GetDetails(db)
		if err != nil {
			return
		}
		articles = append(articles, article)
	}

	utils.CreateResponse(c, http.StatusOK, &articles)
}

// Get Article godoc
// @Summary     Get article by category id.
// @Tags        Article
// @Produce     json
// @Param id path string true "category id"
// @Success     200 {object} []models.Article
// @Router      /articles/category/{id} [get]
func GetArticleByCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var category models.Category

	if err := db.Where("id=?", c.Param("id")).First(&category).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	var articles []models.Article
	var categories []models.ArticleCategory

	if err := db.Where("category_id=?", category.ID).Find(&categories).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	for _, category := range categories {
		article := models.Article{}
		err := db.Where("id=?", category.ArticleID).First(&article).Error

		article.GetDetails(db)
		if err != nil {
			return
		}
		articles = append(articles, article)
	}

	utils.CreateResponse(c, http.StatusOK, &articles)
}

// Create Article godoc
// @Summary     Create Article.
// @Tags        Article
// @Produce     json
// @Param Body body ArticleInput true "body for create article (example ids input: '1,2,3')"
// @Success     201 {object} models.Article
// @Router      /articles [post]
// @Security ApiKeyAuth
func CreateArticle(c *gin.Context) {
	var errs = []string{}
	var input ArticleInput

	db := c.MustGet("db").(*gorm.DB)

	userID, err := utils.ExtractTokenID(c)

	if err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	article := models.Article{
		Title:       input.Title,
		ImageUrl:    input.ImageUrl,
		Content:     input.Content,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	article.GetSlug(db)
	article.UserID = userID

	if err := article.Validate(db); len(err) > 0 {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	if err := db.Create(&article).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	categories := strings.Split(input.Categories, ",")
	tags := strings.Split(input.Tags, ",")

	category_ids := utils.SliceStringToUInt(categories)
	tag_ids := utils.SliceStringToUInt(tags)

	errs = append(errs, article.InsertCategories(db, category_ids)...)
	errs = append(errs, article.InsertTags(db, tag_ids, strings.Split(input.TagsNew, ","))...)

	if len(errs) > 0 {
		if err := article.Delete(db); err != nil {
			fmt.Println(err.Error())
		}
		utils.CreateResponse(c, http.StatusBadRequest, errs)
		return
	}

	article.GetDetails(db)
	utils.CreateResponse(c, http.StatusCreated, &article)
}

// Update Article godoc
// @Summary     Update Article.
// @Tags        Article
// @Produce     json
// @Param id path string true "article id"
// @Param Body body ArticleInput true "body for update article (example ids input: '1,2,3')"
// @Success     200 {object} models.Article
// @Router      /articles/{id} [put]
// @Security ApiKeyAuth
func UpdateArticle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input ArticleInput
	article := models.Article{}
	details := models.Article{}

	if err := db.Where("id=?", c.Param("id")).First(&article).First(&details).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	details.GetDetails(db)

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var errs = []string{}

	updated := models.Article{
		Title:       input.Title,
		Content:     input.Content,
		Description: input.Description,
		ImageUrl:    input.ImageUrl,
		UpdatedAt:   time.Now(),
	}

	updated.GetSlug(db)

	if err := updated.Validate(db); len(err) > 0 {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	if err := article.BeforeUpdate(db); err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Model(&article).Updates(updated).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	categories := strings.Split(input.Categories, ",")
	tags := strings.Split(input.Tags, ",")

	category_ids := utils.SliceStringToUInt(categories)
	tag_ids := utils.SliceStringToUInt(tags)

	errs = append(errs, article.InsertCategories(db, category_ids)...)
	errs = append(errs, article.InsertTags(db, tag_ids, strings.Split(input.TagsNew, ","))...)

	if len(errs) > 0 {
		article.RestoreUpdate(db, &details)
		utils.CreateResponse(c, http.StatusBadRequest, errs)
		return
	}

	article.GetDetails(db)
	utils.CreateResponse(c, http.StatusOK, article)
}

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
	var details models.Article

	if err := db.Where("id=?", c.Param("id")).First(&article).First(&details).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	details.GetDetails(db)

	category_ids := []uint{}
	for _, c := range details.Categories {
		category_ids = append(category_ids, c.CategoryID)
	}

	tag_ids := []uint{}
	for _, t := range details.Tags {
		tag_ids = append(tag_ids, t.TagID)
	}

	if err := db.Model(&article).Update("is_published", true).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var errs = []string{}

	errs = append(errs, article.InsertCategories(db, category_ids)...)
	errs = append(errs, article.InsertTags(db, tag_ids, []string{})...)

	if len(errs) > 0 {
		if err := article.Delete(db); err != nil {
			fmt.Println(err.Error())
		}
		utils.CreateResponse(c, http.StatusBadRequest, errs)
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
	var details models.Article

	if err := db.Where("id=?", c.Param("id")).First(&article).First(&details).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	details.GetDetails(db)

	category_ids := []uint{}
	for _, c := range details.Categories {
		category_ids = append(category_ids, c.CategoryID)
	}

	tag_ids := []uint{}
	for _, t := range details.Tags {
		tag_ids = append(tag_ids, t.TagID)
	}

	if err := db.Model(&article).Update("is_published", false).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var errs = []string{}

	errs = append(errs, article.InsertCategories(db, category_ids)...)
	errs = append(errs, article.InsertTags(db, tag_ids, []string{})...)

	if len(errs) > 0 {
		if err := article.Delete(db); err != nil {
			fmt.Println(err.Error())
		}
		utils.CreateResponse(c, http.StatusBadRequest, errs)
		return
	}

	utils.CreateResponse(c, http.StatusOK, &article)
}

// Get Comments by Article ID godoc
// @Summary     Get Comments by Article ID.
// @Tags        Article Comment
// @Produce     json
// @Param id path string true "article id"
// @Success     200 {object} []models.ArticleComment
// @Router      /articles/{id}/comments [get]
func GetComments(c *gin.Context) {
	var comments []models.ArticleComment

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("article_id=? AND is_reply=false", c.Param("id")).Find(&comments).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	_comments := []models.ArticleComment{}

	for _, comment := range comments {
		if err := comment.GetDetails(db); err != nil {
			utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		_comments = append(_comments, comment)
	}

	utils.CreateResponse(c, http.StatusOK, &_comments)
}

// Create Comment godoc
// @Summary     Create Comment.
// @Tags        Article Comment
// @Produce     json
// @Param id path string true "article id"
// @Param Body body CommentInput true "body for create user"
// @Success     200 {object} models.ArticleComment
// @Router      /articles/{id}/comments [post]
// @Security ApiKeyAuth
func CreateComment(c *gin.Context) {
	var input CommentInput
	var article models.Article

	userID, err := utils.ExtractTokenID(c)

	if err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("id=?", c.Param("id")).First(&article).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	comment := models.ArticleComment{
		Content:   input.Content,
		ArticleID: article.ID,
		UserID:    userID,
		IsReply:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&comment).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusCreated, &comment)
}

// Delete Comment godoc
// @Summary     Delete Comment.
// @Tags        Article Comment
// @Produce     json
// @Param id path string true "comment id"
// @Success     200 {object} bool
// @Router      /articles/comments/{id} [delete]
// @Security ApiKeyAuth
func DeleteComment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var comment models.ArticleComment

	userID, err := utils.ExtractTokenID(c)

	if err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Where("id=?", c.Param("id")).First(&comment).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	if userID != comment.UserID {
		utils.CreateResponse(c, http.StatusBadRequest, "hanya pembuat komentar yang dapat menghapus komentar")
		return
	}

	if err := db.Delete(&comment).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, true)
}

// Get Comments by Comment ID godoc
// @Summary     Get Reply Comments by Comment ID.
// @Tags        Article Comment
// @Produce     json
// @Param id path string true "comment id"
// @Success     200 {object} []models.ReplyArticleComment
// @Router      /articles/comments/{id}/replies [get]
func GetReplyComments(c *gin.Context) {
	var comments []models.ReplyArticleComment

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("parent_id=?", c.Param("id")).Find(&comments).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	_comments := []models.ReplyArticleComment{}

	for _, comment := range comments {
		if err := comment.GetDetails(db); err != nil {
			utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		if err := comment.GetParent(db); err != nil {
			utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		_comments = append(_comments, comment)
	}

	utils.CreateResponse(c, http.StatusOK, &_comments)
}

// Create Reply Comment godoc
// @Summary     Create Reply Comment.
// @Tags        Article Comment
// @Produce     json
// @Param id path string true "comment id"
// @Param Body body CommentInput true "body for create reply comment"
// @Success     200 {object} models.ReplyArticleComment
// @Router      /articles/comments/{id}/replies [post]
// @Security ApiKeyAuth
func CreateReplyComment(c *gin.Context) {
	var input CommentInput
	var parent models.ArticleComment

	userID, err := utils.ExtractTokenID(c)

	if err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("id=?", c.Param("id")).First(&parent).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	comment := models.ArticleComment{
		Content:   input.Content,
		ArticleID: parent.ArticleID,
		UserID:    userID,
		IsReply:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&comment).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	replyComment := models.ReplyArticleComment{
		Content:   input.Content,
		ArticleID: parent.ArticleID,
		ParentID:  parent.ID,
		UserID:    userID,
		CommentID: comment.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&replyComment).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusCreated, &comment)
}

// Delete Reply Comment godoc
// @Summary     Delete Reply Comment.
// @Tags        Article Comment
// @Produce     json
// @Param id path string true "reply comment id"
// @Success     200 {object} bool
// @Router      /articles/comments/replies/{id} [delete]
// @Security ApiKeyAuth
func DeleteReplyComment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var replyComment models.ReplyArticleComment
	var comment models.ArticleComment

	userID, err := utils.ExtractTokenID(c)

	if err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Where("id=?", c.Param("id")).First(&replyComment).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	if err := db.Where("id=?", replyComment.CommentID).First(&comment).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	if userID != comment.UserID && userID != replyComment.UserID {
		utils.CreateResponse(c, http.StatusBadRequest, "hanya pembuat komentar yang dapat menghapus komentar")
		return
	}

	if err := db.Delete(&comment).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Delete(&replyComment).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, true)
}
