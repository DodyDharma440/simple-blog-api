package utils

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateResponse(c *gin.Context, status int, data interface{}) {
	ok := false
	msg := "Failed"

	strStatus := strconv.Itoa(status)

	if strings.Split(strStatus, "")[0] == "2" {
		ok = true
		msg = "Success"
	}

	var res map[string]interface{}

	if ok {
		res = gin.H{
			"ok":      ok,
			"status":  status,
			"message": msg,
			"data":    data,
		}
	} else {
		res = gin.H{
			"ok":      ok,
			"status":  status,
			"message": msg,
			"error":   data,
		}
	}

	c.JSON(status, res)
}
