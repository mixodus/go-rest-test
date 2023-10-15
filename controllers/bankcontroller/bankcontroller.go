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
	res := dto.Response{
		Status:  true,
		Message: "Success",
		Data:    banks,
	}
	c.JSON(http.StatusOK, res)
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
			res := dto.Response{
				Status:  false,
				Message: out[0].Field + ", " + out[0].Message,
				Data:    nil,
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, res)
		}

		return
	}

	//CHECK BANK EXIST
	var bank models.Bank
	if err := models.DB.Where("bank_code = ?", playerBank.BankCode).First(&bank).Error; err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Bank not found",
			Data:    nil,
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	//GET USER BY TOKEN
	player, err := services.GetUserByToken(c)
	if err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Failed",
			Data:    nil,
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
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
			res := dto.Response{
				Status:  false,
				Message: "Failed to update players bank id",
				Data:    err,
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}

		res := dto.Response{
			Status:  true,
			Message: "Success",
			Data:    playersBank,
		}
		c.JSON(http.StatusOK, res)
		return
	}

	//IF PLAYER BANK EXIST THEN RETURN ERROR
	res := dto.Response{
		Status:  false,
		Message: "Player bank already exist",
		Data:    nil,
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, res)
}

// DELETE PLAYER BANK BY PLAYER TOKEN
func RemovePlayerBank(c *gin.Context) {
	//GET PLAYER TOKEN
	player, err := services.GetUserByToken(c)
	if err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Failed",
			Data:    nil,
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	//CHECK PLAYER BANK EXIST, IF NOT EXIST THEN RETURN ERROR
	var playersBank models.PlayersBank
	if err := models.DB.Where("player_id = ?", player.Id).First(&playersBank).Error; err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Player's bank not found",
			Data:    nil,
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	//UPDATE PLAYER'S BANK ID
	player.PlayersBankId = nil
	if err := models.DB.Where("id = ?", player.Id).Save(&player).Error; err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Failed to update players bank id",
			Data:    err,
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	//DELETE PLAYER BANK
	models.DB.Delete(&playersBank)
	res := dto.Response{
		Status:  true,
		Message: "Success",
		Data:    nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetPlayerBank(c *gin.Context) {
	//GET PLAYER TOKEN
	player, err := services.GetUserByToken(c)
	if err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Failed",
			Data:    nil,
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	//CHECK PLAYER BANK EXIST, IF NOT EXIST THEN RETURN ERROR
	var playersBank models.PlayersBank
	if err := models.DB.Joins("Bank").Where("player_id = ?", player.Id).First(&playersBank).Error; err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Player's bank not found",
			Data:    nil,
		}
		if err.Error() == "record not found" {
			c.AbortWithStatusJSON(http.StatusOK, res)
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := dto.Response{
		Status:  true,
		Message: "Success",
		Data:    playersBank,
	}
	c.JSON(http.StatusOK, res)
}
