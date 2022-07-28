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

type ChangePasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// Get All Users godoc
// @Summary     Get all users.
// @Tags        User
// @Produce     json
// @Success     200 {object} []models.User
// @Router      /users [get]
// @Security ApiKeyAuth
func GetUsers(c *gin.Context) {
	var users []models.User

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Find(&users).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, users)
}

// Get User By ID godoc
// @Summary     Get user.
// @Tags        User
// @Produce     json
// @Param id path string true "user id"
// @Success     200 {object} models.User
// @Router      /users/{id} [get]
// @Security ApiKeyAuth
func GetUser(c *gin.Context) {
	var user models.User

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id=?", c.Param("id")).First(&user).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	utils.CreateResponse(c, http.StatusOK, user)
}

// Create User godoc
// @Summary     Create user.
// @Tags        User
// @Produce     json
// @Param Body body UserInput true "body for create user"
// @Success     200 {object} models.User
// @Router      /users [post]
// @Security ApiKeyAuth
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

	db := c.MustGet("db").(*gorm.DB)

	validate := user.Validate(db)

	if len(validate) > 0 {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, validate)
		return
	}

	if err := user.BeforeSave(db); err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Create(&user).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusCreated, &user)
}

// Update User godoc
// @Summary     Update user.
// @Tags        User
// @Produce     json
// @Param id path string true "user id"
// @Param Body body UserInput true "body for update user"
// @Success     200 {object} models.User
// @Router      /users/{id} [put]
// @Security ApiKeyAuth
func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var user models.User

	if err := db.Where("id=?", c.Param("id")).First(&user).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	validate := user.Validate(db)

	if len(validate) > 0 {
		utils.CreateResponse(c, http.StatusUnprocessableEntity, validate)
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

// Delete User godoc
// @Summary     Delete user.
// @Tags        User
// @Produce     json
// @Param id path string true "user id"
// @Success     200 {object} bool
// @Router      /users/{id} [delete]
// @Security ApiKeyAuth
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

// Login User godoc
// @Summary     Login user.
// @Description Login User.
// @Tags        Auth
// @Param Body body LoginInput true "the body to login"
// @Produce     json
// @Router      /login [post]
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

// Register User godoc
// @Summary     Register user.
// @Tags        Auth
// @Produce     json
// @Param Body body UserInput true "body for register user"
// @Success     200 {object} models.User
// @Router      /register [post]
func RegisterUser(c *gin.Context) {
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

	db := c.MustGet("db").(*gorm.DB)

	if err := user.BeforeSave(db); err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Create(&user).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusCreated, &user)
}

// Change Password User godoc
// @Summary     Change Password user.
// @Tags        Auth
// @Produce     json
// @Param Body body ChangePasswordInput true "body for change user password"
// @Success     200 {object} models.User
// @Router      /change-password [patch]
// @Security ApiKeyAuth
func ChangePassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID, err := utils.ExtractTokenID(c)

	if err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input ChangePasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.CreateResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	u := models.User{}

	if err := db.Model(models.User{}).Where("id=?", userID).Take(&u).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, err.Error())
		return
	}

	if err := models.VerifyPassword(u.Password, input.OldPassword); err != nil {
		utils.CreateResponse(c, http.StatusBadRequest, "Password lama tidak cocok")
		return
	}

	var updated models.User
	updated.Password = input.NewPassword
	updated.BeforeSave(db)

	if err := db.Model(&u).Updates(updated).Error; err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreateResponse(c, http.StatusOK, "Berhasil ganti password")
}

// Get User Profile godoc
// @Summary     Get user profile.
// @Tags        Auth
// @Produce     json
// @Success     200 {object} models.User
// @Router      /my-profile [get]
// @Security ApiKeyAuth
func MyProfile(c *gin.Context) {
	userID, err := utils.ExtractTokenID(c)

	if err != nil {
		utils.CreateResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var user models.User

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id=?", userID).First(&user).Error; err != nil {
		utils.CreateResponse(c, http.StatusNotFound, "data not found")
		return
	}

	utils.CreateResponse(c, http.StatusOK, user)
}
