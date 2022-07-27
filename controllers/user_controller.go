package controllers

import (
	"final-project/models"
	"net/http"
	"time"

	"final-project/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetUsers(c *gin.Context) {
	var users []models.User

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Find(&users).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	var user models.User

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id=?", c.Param("id")).First(&user).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	utils.CreateResponse(c, http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var input UserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := user.BeforeSave(); err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Create(&user).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusCreated, &user)
}

func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var user models.User

	if err := db.Where("id=?", c.Param("id")).First(&user).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	var input UserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var updated models.User

	updated.Name = input.Name
	updated.Email = input.Email
	updated.UpdatedAt = time.Now()

	if err := db.Model(&user).Updates(updated).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, &user)
}

func DeleteUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var user models.User

	if err := db.Where("id=?", c.Param("id")).First(&user).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, true)
}

func LoginUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.CreateResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	u := models.User{}

	u.Email = input.Email
	u.Password = input.Password

	token, err := u.LoginCheck(db)

	if err != nil {
		utils.CreateResponse(c, http.StatusBadRequest, "email atau password salah")
		return
	}

	utils.CreateResponse(c, http.StatusOK, token)
}
