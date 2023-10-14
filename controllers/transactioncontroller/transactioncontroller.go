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
	"github.com/mixodus/go-rest-test/models"
	"github.com/mixodus/go-rest-test/services"
)

func TopUp(c *gin.Context) {
	//GET PLAYER TOKEN
	player, err := services.GetUserByToken(c)
	if err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Player might not exist...",
			Data:    err.Error(),
		}
		c.JSON(http.StatusNotFound, res)
		return
	}
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
		"image/jpeg": true,
		"image/png":  true,
		// Add more allowed MIME types as needed
	}

	//RETURN ERROR IF FILE TYPE NOT ALLOWED
	fileType := services.GetFileType(file)
	if !allowedMimeTypes[fileType] {
		res := dto.Response{
			Status:  false,
			Message: "Image JPEG & PNG only",
			Data:    fileType,
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	//CREATE UPLOAD DIRECTORY
	uploadDir := "./uploads/transaction"
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

	//GET PLAYER BANK
	var playersBank models.PlayersBank
	if err := models.DB.Where("player_id = ?", player.Id).First(&playersBank).Error; err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Player's might not have bank account",
			Data:    nil,
		}
		if err.Error() == "record not found" {
			c.AbortWithStatusJSON(http.StatusOK, res)
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	//SAVE TRANSACTION
	var transaction models.Transaction
	transaction.PlayerId = player.Id
	transaction.PlayersBankId = playersBank.Id
	transaction.Amount = topUpRequest.Amount
	transaction.TransactionType = models.DEPOSIT
	transaction.FileName = uniqueFileName
	if err := models.DB.Create(&transaction).Error; err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Error saving transaction",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := dto.Response{
		Status:  true,
		Message: "File",
		Data:    transaction,
	}

	c.JSON(http.StatusOK, res)
}
