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
		services.Response(c, http.StatusNotFound, false, "Player might not exist...", err.Error())
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
			services.Response(c, http.StatusBadRequest, false, out[0].Field+", "+out[0].Message, nil)
		} else {
			services.Response(c, http.StatusBadRequest, false, "Something wrong!", err.Error())
		}
		return
	}

	//GET FILE
	file, err := c.FormFile("file")
	if err != nil {
		services.Response(c, http.StatusBadRequest, false, "No file uploaded", file)
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
		services.Response(c, http.StatusBadRequest, false, "Image JPEG & PNG only", fileType)
		return
	}

	//CREATE UPLOAD DIRECTORY
	uploadDir := "./uploads/transaction"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		services.Response(c, http.StatusInternalServerError, false, "Error creating directory", nil)
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
		services.Response(c, http.StatusInternalServerError, false, "Error saving the file", nil)
		return
	}

	//GET PLAYER BANK
	var playersBank models.PlayersBank
	if err := models.DB.Where("player_id = ?", player.Id).First(&playersBank).Error; err != nil {
		if err.Error() == "record not found" {
			services.Response(c, http.StatusOK, true, "Player's might not have bank account", nil)
			return
		}
		services.Response(c, http.StatusBadRequest, false, "Failed to get player's bank", nil)
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
		services.Response(c, http.StatusInternalServerError, false, "Error saving transaction", nil)
		return
	}

	services.Response(c, http.StatusOK, true, "Transaction success", transaction)
}

func Spent(c *gin.Context) {
	//GET PLAYER TOKEN
	player, err := services.GetUserByToken(c)
	if err != nil {
		services.Response(c, http.StatusNotFound, false, "Player might not exist...", err.Error())
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
			services.Response(c, http.StatusBadRequest, false, out[0].Field+", "+out[0].Message, nil)
		} else {
			services.Response(c, http.StatusBadRequest, false, "Something wrong!", err.Error())
		}
		return
	}

	//GET PLAYER BANK
	var playersBank models.PlayersBank
	if err := models.DB.Where("player_id = ?", player.Id).First(&playersBank).Error; err != nil {
		if err.Error() == "record not found" {
			services.Response(c, http.StatusOK, true, "Player's might not have bank account", nil)
			return
		}
		services.Response(c, http.StatusBadRequest, false, "Failed to get player's bank", err.Error())
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
		services.Response(c, http.StatusInternalServerError, false, "Error saving transaction", nil)
		return
	}

	services.Response(c, http.StatusOK, true, "Transaction success", transaction)
}

func SetAllDebitStatusSuccess(c *gin.Context) {
	//UPDATE ALL DEBIT TRANSACTION STATUS TO SUCCESS
	if err := models.DB.Model(&models.Transaction{}).Where("status = ? AND transaction_type = ?", models.PENDING, models.DEBIT).Update("status", models.SUCCESS).Error; err != nil {
		services.Response(c, http.StatusBadRequest, false, "Something Wrong!", err.Error())
		return
	}

	services.Response(c, http.StatusOK, true, "All debit transaction status set to success", nil)
}

func SetAllCreditStatusSuccess(c *gin.Context) {
	//UPDATE ALL CREDIT TRANSACTION STATUS TO SUCCESS
	if err := models.DB.Model(&models.Transaction{}).Where("status = ? AND transaction_type = ?", models.PENDING, models.CREDIT).Update("status", models.SUCCESS).Error; err != nil {
		services.Response(c, http.StatusBadRequest, false, "Something Wrong!", err.Error())
		return
	}

	services.Response(c, http.StatusOK, true, "All credit transaction status set to success", nil)
}

func GetAndUpdateWallet(c *gin.Context) {
	//==== GET AND UPDATE USER BALANCE ====

	//GET PLAYER TOKEN
	player, err := services.GetUserByToken(c)
	if err != nil {
		services.Response(c, http.StatusNotFound, false, "Player might not exist...", err.Error())
		return
	}

	//UPDATE PLAYER WALLET
	// GET PLAYER TRANSACTION WHERE STATUS SUCCESS AND TYPE DEBIT
	var debit *int64
	if err := models.DB.Table("transactions").Select("SUM(amount) as total_debit").Where("player_id = ? AND status = ? AND transaction_type = ?", player.Id, models.SUCCESS, models.DEBIT).Row().Scan(&debit); err != nil {
		if err.Error() != "record not found" {
			services.Response(c, http.StatusBadRequest, false, "Something Wrong!", err.Error())
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
		services.Response(c, http.StatusBadRequest, false, "Something wrong!", err.Error())
	}
	if credit == nil {
		temp := int64(0)
		credit = &temp
	}

	//UPDATE PLAYER WALLET
	player.Balance = *debit - *credit
	if err := models.DB.Save(&player).Error; err != nil {
		services.Response(c, http.StatusBadRequest, false, "Something Wrong!", err.Error())
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

	services.Response(c, http.StatusOK, true, "Get & update player's balance", dataRes)
}
