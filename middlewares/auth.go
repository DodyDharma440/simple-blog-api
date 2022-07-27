package middlewares

import (
	"final-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := utils.TokenValid(c)
		if err != nil {
			utils.CreateResponse(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}
