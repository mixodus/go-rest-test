package bankcontroller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mixodus/go-rest-test/dto"
	"github.com/mixodus/go-rest-test/models"
	"github.com/mixodus/go-rest-test/services"
)

// GET ALL BANK LIST
func BankList(c *gin.Context) {
	var banks []models.Bank
	models.DB.Find(&banks)
	services.Response(c, http.StatusOK, true, "Success", banks)
}

// ADD PLAYER BANK
func AddPlayerBank(c *gin.Context) {
	//GET INPUT REQUEST & VALIDATE
	var playerBank dto.AddBankRequest
	if err := c.ShouldBindJSON(&playerBank); err != nil {
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

	//CHECK BANK EXIST
	var bank models.Bank
	if err := models.DB.Where("bank_code = ?", playerBank.BankCode).First(&bank).Error; err != nil {
		services.Response(c, http.StatusBadRequest, false, "Bank not found", nil)
		return
	}

	//GET USER BY TOKEN
	player, err := services.GetUserByToken(c)
	if err != nil {
		services.Response(c, http.StatusBadRequest, false, "Failed to get user by token", nil)
		return
	}

	//CHECK PLAYER BANK EXIST, IF NOT EXIST THEN CREATE
	var playersBank models.PlayersBank
	if err := models.DB.Where("player_id = ? AND deleted_at IS NULL", player.Id).First(&playersBank).Error; err != nil {
		playersBank := models.PlayersBank{
			PlayerId:          player.Id,
			BankId:            bank.Id,
			BankAccountNumber: playerBank.AccountNumber,
			BankAccountName:   playerBank.AccountName,
			Bank:              bank,
		}
		models.DB.Create(&playersBank)

		//UPDATE PLAYER'S BANK ID
		playersBankId := &playersBank.Id
		player.PlayersBankId = playersBankId
		if err := models.DB.Where("id = ?", player.Id).Save(&player).Error; err != nil {
			services.Response(c, http.StatusBadRequest, false, "Failed to update players bank id", err.Error())
			return
		}

		services.Response(c, http.StatusOK, true, "Success", playersBank)
		return
	}

	//IF PLAYER BANK EXIST THEN RETURN ERROR
	services.Response(c, http.StatusBadRequest, false, "Player bank already exist", nil)
}

// DELETE PLAYER BANK BY PLAYER TOKEN
func RemovePlayerBank(c *gin.Context) {
	//GET PLAYER TOKEN
	player, err := services.GetUserByToken(c)
	if err != nil {
		services.Response(c, http.StatusBadRequest, false, "Failed to get user by token", nil)
		return
	}

	//CHECK PLAYER BANK EXIST, IF NOT EXIST THEN RETURN ERROR
	var playersBank models.PlayersBank
	if err := models.DB.Where("player_id = ?", player.Id).First(&playersBank).Error; err != nil {
		services.Response(c, http.StatusBadRequest, false, "Player's bank not found", nil)
		return
	}

	//UPDATE PLAYER'S BANK ID
	player.PlayersBankId = nil
	if err := models.DB.Where("id = ?", player.Id).Save(&player).Error; err != nil {
		services.Response(c, http.StatusBadRequest, false, "Failed to update players bank id", err.Error())
		return
	}

	//DELETE PLAYER BANK
	models.DB.Delete(&playersBank)
	services.Response(c, http.StatusOK, true, "Success", nil)
}

func GetPlayerBank(c *gin.Context) {
	//GET PLAYER TOKEN
	player, err := services.GetUserByToken(c)
	if err != nil {
		services.Response(c, http.StatusBadRequest, false, "Failed to get user by token", nil)
		return
	}

	//CHECK PLAYER BANK EXIST, IF NOT EXIST THEN RETURN ERROR
	var playersBank models.PlayersBank
	if err := models.DB.Joins("Bank").Where("player_id = ?", player.Id).First(&playersBank).Error; err != nil {
		if err.Error() == "record not found" {
			services.Response(c, http.StatusBadRequest, false, "Player's bank not found", nil)
			return
		}
		services.Response(c, http.StatusBadRequest, false, "Failed to get player's bank", err.Error())
		return
	}

	services.Response(c, http.StatusOK, true, "Success", playersBank)
}
