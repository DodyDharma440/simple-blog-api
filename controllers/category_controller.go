package controllers

import (
	"final-project/models"
	"final-project/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryInput struct {
	Name string `json:"name"`
}

// Get All Categories godoc
// @Summary     Get all categories.
// @Tags        Category
// @Produce     json
// @Success     200 {object} []models.Category
// @Router      /categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Find(&categories).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, categories)
}

// Get Category By ID godoc
// @Summary     Get Category.
// @Tags        Category
// @Produce     json
// @Param id path string true "category id"
// @Success     200 {object} models.Category
// @Router      /categories/{id} [get]
func GetCategory(c *gin.Context) {
	var category models.Category

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("id=?", c.Param("id")).First(&category).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	utils.CreateResponse(c, http.StatusOK, category)
}

// Create Category godoc
// @Summary     Create Category.
// @Tags        Category
// @Produce     json
// @Param Body body CategoryInput true "body for create category"
// @Success     200 {object} models.Category
// @Router      /categories [post]
// @Security ApiKeyAuth
func CreateCategory(c *gin.Context) {
	var input CategoryInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	category := models.Category{
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db := c.MustGet("db").(*gorm.DB)

	if err := category.Validate(); len(err) > 0 {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	if err := db.Create(&category).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusCreated, &category)
}

// Update Category godoc
// @Summary     Update Category.
// @Tags        Category
// @Produce     json
// @Param id path string true "category id"
// @Param Body body CategoryInput true "body for update category"
// @Success     200 {object} models.Category
// @Router      /categories/{id} [put]
// @Security ApiKeyAuth
func UpdateCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var category models.Category

	if err := db.Where("id", c.Param("id")).First(&category).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	var input CategoryInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	updated := models.Category{
		Name:      input.Name,
		UpdatedAt: time.Now(),
	}

	if err := updated.Validate(); len(err) > 0 {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	if err := db.Model(&category).Updates(updated).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, &category)
}

// Delete Category godoc
// @Summary     Delete Category.
// @Tags        Category
// @Produce     json
// @Param id path string true "category id"
// @Success     200 {object} bool
// @Router      /categories/{id} [delete]
// @Security ApiKeyAuth
func DeleteCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var category models.Category

	if err := db.Where("id=?", c.Param("id")).First(&category).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	if err := db.Delete(&category).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, true)
}
