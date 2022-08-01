package controllers

import (
	"final-project/models"
	"final-project/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TagInput struct {
	Name string `json:"name"`
}

// Get All Tags godoc
// @Summary     Get all tags.
// @Tags        Tag
// @Produce     json
// @Success     200 {object} []models.Tag
// @Router      /tags [get]
func GetTags(c *gin.Context) {
	var tags []models.Tag

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Find(&tags).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, tags)
}

// Get Tag By ID godoc
// @Summary     Get Tag.
// @Tags        Tag
// @Produce     json
// @Param id path string true "tag id"
// @Success     200 {object} models.Tag
// @Router      /tags/{id} [get]
func GetTag(c *gin.Context) {
	var tag models.Tag

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("id=?", c.Param("id")).First(&tag).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	utils.CreateResponse(c, http.StatusOK, tag)
}

// Create Tag godoc
// @Summary     Create Tag.
// @Tags        Tag
// @Produce     json
// @Param Body body TagInput true "body for create tag"
// @Success     200 {object} models.Tag
// @Router      /tags [post]
// @Security ApiKeyAuth
func CreateTag(c *gin.Context) {
	var input TagInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	tag := models.Tag{
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db := c.MustGet("db").(*gorm.DB)

	if err := tag.Validate(); len(err) > 0 {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	if err := db.Create(&tag).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusCreated, &tag)
}

// Update Tag godoc
// @Summary     Update Tag.
// @Tags        Tag
// @Produce     json
// @Param id path string true "tag id"
// @Param Body body TagInput true "body for update tag"
// @Success     200 {object} models.Tag
// @Router      /tags/{id} [put]
// @Security ApiKeyAuth
func UpdateTag(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var tag models.Tag

	if err := db.Where("id", c.Param("id")).First(&tag).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	var input TagInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	updated := models.Tag{
		Name:      input.Name,
		UpdatedAt: time.Now(),
	}

	if err := updated.Validate(); len(err) > 0 {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	if err := db.Model(&tag).Updates(updated).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, &tag)
}

// Delete Tag godoc
// @Summary     Delete Tag.
// @Tags        Tag
// @Produce     json
// @Param id path string true "tag id"
// @Success     200 {object} bool
// @Router      /tags/{id} [delete]
// @Security ApiKeyAuth
func DeleteTag(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var tag models.Tag

	if err := db.Where("id=?", c.Param("id")).First(&tag).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	if err := db.Delete(&tag).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, true)
}
