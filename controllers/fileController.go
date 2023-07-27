package controllers

import (
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ImageUpload(c *gin.Context) {

	r := gin.Default()

	r.Static("../assets", "../assets")
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// single file
	image, err := c.FormFile("image")

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "Faild to upload image" + err.Error(),
		})
		return
	}
	if !isImage(image) {
		c.JSON(http.StatusOK, gin.H{
			"error": "Only image files are allowed",
		})
		return
	}
	err = c.SaveUploadedFile(image, "assets/uploads/user/"+image.Filename)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "Faild to upload image",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "/assets/uploads/user/" + image.Filename,
	})
}
func isImage(file *multipart.FileHeader) bool {
	extension := filepath.Ext(file.Filename)
	switch extension {
	case ".jpg", ".JPG", ".jpeg", ".png", ".gif", ".bmp":
		return true
	}
	return false
}