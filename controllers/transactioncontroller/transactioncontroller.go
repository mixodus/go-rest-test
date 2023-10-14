package transactioncontroller

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mixodus/go-rest-test/dto"
	"github.com/mixodus/go-rest-test/services"
)

func TopUp(c *gin.Context) {
	//GET INPUT REQUEST & VALIDATE
	var topUpRequest dto.TopUpRequest
	if err := c.ShouldBind(&topUpRequest); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]services.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = services.ErrorMsg{Field: fe.Field(), Message: services.GetErrorMsg(fe)}
			}
			res := dto.Response{
				Status:  false,
				Message: out[0].Field + ", " + out[0].Message,
				Data:    nil,
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, res)
		} else {
			res := dto.Response{
				Status:  false,
				Message: "Something wrong!",
				Data:    err,
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, res)
		}
		return
	}

	//GET FILE
	file, err := c.FormFile("file")
	if err != nil {
		res := dto.Response{
			Status:  false,
			Message: "No file uploaded",
			Data:    file,
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	//CHECK FILE TYPE
	allowedMimeTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"application/pdf": true,
		// Add more allowed MIME types as needed
	}

	//RETURN ERROR IF FILE TYPE NOT ALLOWED
	fileType := services.GetFileType(file)
	if !allowedMimeTypes[fileType] {
		res := dto.Response{
			Status:  false,
			Message: "JPEG, PNG, PDF only",
			Data:    fileType,
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	//CREATE UPLOAD DIRECTORY
	uploadDir := "./uploads/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Error creating directory",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	currentTime := time.Now()
	// Convert the time.Time to a Unix timestamp (in seconds)
	timestamp := currentTime.Unix()
	timestamp = int64(timestamp)
	// Convert the timestamp to a string
	timestampStr := strconv.FormatInt(timestamp, 10)
	uniqueFileName := filepath.Join(uploadDir, timestampStr+"-"+file.Filename)
	// Save the file to the specified directory
	if err := c.SaveUploadedFile(file, uniqueFileName); err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Error saving the file",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := dto.Response{
		Status:  true,
		Message: "File",
		Data:    file,
	}

	c.JSON(http.StatusOK, res)
}
