package middlewares

import (
	"final-project/models"
	"final-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(*gorm.DB)
		userID, err := utils.ExtractTokenID(c)
		user := models.User{}

		if err != nil {
			utils.CreateResponse(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		if err := db.Where("id=?", userID).First(&user).Error; err != nil {
			utils.CreateResponse(c, http.StatusNotFound, "user : "+err.Error())
			c.Abort()
			return
		}

		if user.Role != "admin" {
			utils.CreateResponse(c, http.StatusBadRequest, "hanya admin yang dapat melakukan aksi ini")
			c.Abort()
			return
		}

		c.Next()
	}
}
