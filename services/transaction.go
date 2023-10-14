package services

import (
	"mime/multipart"
	"net/http"
)

func GetFileType(file *multipart.FileHeader) string {
	// Open the file and read the first 512 bytes to detect the MIME type
	src, err := file.Open()
	if err != nil {
		return ""
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return ""
	}

	// Detect the MIME type
	fileType := http.DetectContentType(buffer)
	return fileType
}
