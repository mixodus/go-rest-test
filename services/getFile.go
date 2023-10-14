package services

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetImage(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(400, gin.H{
			"message": "File path is required",
		})
		return
	}
	var contentType string
	if isJPEG(filePath) {
		contentType = "image/jpeg"
	} else if isPNG(filePath) {
		contentType = "image/png"
	} else if isGIF(filePath) {
		contentType = "image/gif"
	} else {
		// Set a default Content-Type for other image formats
		contentType = "image/*"
	}

	// Set the Content-Type to "application/octet-stream" for generic file serving
	c.Header("Content-Disposition", "inline")
	// Set the Content-Type header
	c.Header("Content-Type", contentType)
	// Serve the file
	c.File(filePath) // Replace with the actual path to your file
}

func isJPEG(filename string) bool {
	return strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg")
}

func isPNG(filename string) bool {
	return strings.HasSuffix(filename, ".png")
}

func isGIF(filename string) bool {
	return strings.HasSuffix(filename, ".gif")
}
