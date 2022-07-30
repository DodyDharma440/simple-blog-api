package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SliceStringToUInt(data []string) []uint {
	ids := []uint{}

	for _, d := range data {
		id, err := strconv.Atoi(d)

		if err != nil {
			fmt.Println(err.Error())
		}

		ids = append(ids, uint(id))
	}
	return ids
}

func UploadFile(c *gin.Context, folder string, name string) (string, error) {
	file, header, err := c.Request.FormFile(name)

	if err != nil {
		return "", err
	}

	timeStamp := strconv.Itoa(int(time.Now().Unix()))
	fileName := timeStamp + "-" + header.Filename
	path := fmt.Sprintf("public/upload/%v/%v", folder, fileName)

	out, err := os.Create(path)
	defer out.Close()

	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	filepath := fmt.Sprintf("file/upload/%v/%v", folder, fileName)

	return filepath, nil
}
