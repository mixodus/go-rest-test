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
	transaction.TransactionType = models.DEBIT
	transaction.FileName = uniqueFileName
	transaction.Notes = "TOP UP"
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

func Spent(c *gin.Context) {
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
	var spentRequest dto.SpentRequest
	if err := c.ShouldBindJSON(&spentRequest); err != nil {
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
	transaction.Amount = spentRequest.Amount
	transaction.TransactionType = models.CREDIT
	transaction.Notes = spentRequest.Notes
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
		Message: "Transaction success",
		Data:    transaction,
	}

	c.JSON(http.StatusOK, res)
}

func SetAllDebitStatusSuccess(c *gin.Context) {
	//UPDATE ALL DEBIT TRANSACTION STATUS TO SUCCESS
	if err := models.DB.Model(&models.Transaction{}).Where("status = ? AND transaction_type = ?", models.PENDING, models.DEBIT).Update("status", models.SUCCESS).Error; err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Something Wrong!",
			Data:    err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := dto.Response{
		Status:  true,
		Message: "All debit transaction status set to success",
		Data:    nil,
	}

	c.JSON(http.StatusOK, res)
}

func SetAllCreditStatusSuccess(c *gin.Context) {
	//UPDATE ALL CREDIT TRANSACTION STATUS TO SUCCESS
	if err := models.DB.Model(&models.Transaction{}).Where("status = ? AND transaction_type = ?", models.PENDING, models.CREDIT).Update("status", models.SUCCESS).Error; err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Something Wrong!",
			Data:    err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := dto.Response{
		Status:  true,
		Message: "All credit transaction status set to success",
		Data:    nil,
	}

	c.JSON(http.StatusOK, res)
}

func GetAndUpdateWallet(c *gin.Context) {
	//==== GET AND UPDATE USER BALANCE ====

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

	//UPDATE PLAYER WALLET
	// GET PLAYER TRANSACTION WHERE STATUS SUCCESS AND TYPE DEBIT
	var debit *int64
	if err := models.DB.Table("transactions").Select("SUM(amount) as total_debit").Where("player_id = ? AND status = ? AND transaction_type = ?", player.Id, models.SUCCESS, models.DEBIT).Row().Scan(&debit); err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Something Wrong!",
			Data:    err,
		}
		if err.Error() != "record not found" {
			c.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
	}
	if debit == nil {
		temp := int64(0)
		debit = &temp
	}

	// GET PLAYER TRANSACTION WHERE STATUS SUCCESS AND TYPE CREDIT
	var credit *int64
	if err := models.DB.Table("transactions").Select("SUM(amount) as total_credit").Where("player_id = ? AND status = ? AND transaction_type = ?", player.Id, models.SUCCESS, models.CREDIT).Row().Scan(&credit); err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Something Wrong!",
			Data:    err.Error(),
		}
		if err.Error() != "record not found" {
			c.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
	}
	if credit == nil {
		temp := int64(0)
		credit = &temp
	}

	//UPDATE PLAYER WALLET
	player.Balance = *debit - *credit
	if err := models.DB.Save(&player).Error; err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Something Wrong!",
			Data:    err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	//SUCCESS RESPONSE
	dataRes := dto.ProfileResponse{
		FirstName: player.FirstName,
		LastName:  player.LastName,
		Email:     player.Email,
		Phone:     player.Phone,
		Balance:   player.Balance,
		Id:        player.Id,
	}
	res := dto.Response{
		Status:  true,
		Message: "Get & update player's balance",
		Data:    dataRes,
	}

	c.JSON(http.StatusOK, res)
}
