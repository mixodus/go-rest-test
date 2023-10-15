package playercontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mixodus/go-rest-test/dto"
	"github.com/mixodus/go-rest-test/models"
	"github.com/mixodus/go-rest-test/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	//CALL GET PLAYER AND FILTER SERVICE QUERY
	player, err := services.FilterPlayer(c)
	if err != nil {
		services.Response(c, http.StatusBadRequest, false, "Something went wrong!", err.Error())
		return
	}
	services.Response(c, http.StatusOK, true, "Success", player)
}

func GetPlayerById(c *gin.Context) {
	id := c.Param("id")
	var player dto.Player
	if err := models.DB.Where("id = ?", id).Preload("PlayersBank.Bank").First(&player).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			services.Response(c, http.StatusBadRequest, false, "User not found", err.Error())
			return
		default:
			services.Response(c, http.StatusBadRequest, false, "User not found", err.Error())
			return
		}
	}

	services.Response(c, http.StatusOK, true, "Success", player)
}

func Register(c *gin.Context) {

	//GET INPUT REQUEST & VALIDATE
	var userInput dto.RegisterRequest
	if err := c.ShouldBindJSON(&userInput); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]services.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = services.ErrorMsg{Field: fe.Field(), Message: services.GetErrorMsg(fe)}
			}
			services.Response(c, http.StatusBadRequest, false, out[0].Field+", "+out[0].Message, nil)
			return
		} else {
			services.Response(c, http.StatusBadRequest, false, "Something wrong!", err.Error())
		}
	}

	//HASH PASSWORD
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	//SET DATA TO STORE FROM INPUT
	lowerCaseEmail := services.ToLowerCase(userInput.Email)
	dataPost := models.Player{
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Password:  userInput.Password,
		Email:     lowerCaseEmail,
		Phone:     userInput.Phone,
	}

	//SAVE TO DB
	if err := models.DB.Create(&dataPost).Error; err != nil {
		log.Default().Println(err.Error())
		services.Response(c, http.StatusBadRequest, false, "User already exist", err.Error())
		return
	}

	//SET DTO FOR RESPONSE
	dataRes := dto.RegisterResponse{
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Email:     userInput.Email,
		Phone:     userInput.Phone,
	}

	//RES DATA
	services.Response(c, http.StatusOK, true, "Success", dataRes)
}

func Login(c *gin.Context) {
	// GET REDIS CLIENT
	var redis = services.GetRedisClient()
	// GET INPUT REQUEST & VALIDATE
	var userInput dto.LoginRequest
	if err := c.ShouldBindJSON(&userInput); err != nil {
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

	// ==== LOGIN VIA REDIS ====
	// GET PLAYER DATA STORED IN REDIS
	session := redis.Get(c, userInput.Email).Val()
	fmt.Println("GET SESSION REDIS")
	// IF PLAYER DATA EXIST THEN LOGIN VIA REDIS
	if session != "" {
		res, err := services.LoginViaRedis(c, &userInput, session)
		if err != nil {
			services.Response(c, http.StatusBadRequest, false, "Something went wrong!", err.Error())
			return
		} else {
			services.Response(c, http.StatusOK, true, "Successfully login!", res)
			return
		}
	}
	// ==== END LOGIN VIA REDIS ====

	// ==== LOGIN VIA DATABASE ====
	//GET USERDATA BY USERNAME
	var player models.Player
	if err := models.DB.Where("email = ?", userInput.Email).First(&player).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			services.Response(c, http.StatusBadRequest, false, "User not found", err.Error())
			return
		default:
			services.Response(c, http.StatusBadRequest, false, "Something went wrong!", err.Error())
			return
		}

	}

	//CHECK PASSWORD
	if err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(userInput.Password)); err != nil {
		services.Response(c, http.StatusBadRequest, false, "Wrong password", err.Error())
		return
	}

	tokenz, err := services.GeneratePlayerToken(c, &player)
	if err != nil {
		services.Response(c, http.StatusBadRequest, false, "Something went wrong!", err.Error())
	}

	//==== SET SESSION REDIS ====
	//set session to redis
	playerJson, _ := json.Marshal(player)
	//set user email as user's data key
	redis.Set(c, userInput.Email, playerJson, time.Hour*1)
	//set player id as token key
	redis.Set(c, player.Id, tokenz, time.Hour*24)
	fmt.Println("SET SESSION REDIS")

	//!!UNUSED since we use redis for session management!!
	// ==== save token to db ====
	// var tokenSession models.TokenSession
	// tokenSession.PlayerId = player.Id
	// tokenSession.Token = tokenz
	// if err := models.DB.Where("player_id = ?", player.Id).Assign(tokenSession).FirstOrCreate(&tokenSession).Error; err != nil {
	// 	log.Default().Println(err.Error())
	// 	res := dto.Response{
	// 		Status:  false,
	// 		Message: "Save token failed",
	// 		Data:    err.Error(),
	// 	}
	// 	c.JSON(http.StatusBadRequest, res)
	// 	return
	// }
	// ==== end ====

	// ==== save token to redis ====
	//set token to redis

	//set dto
	data := map[string]interface{}{
		"token": tokenz,
	}

	services.Response(c, http.StatusOK, true, "Successfully login!", data)
}

func Profile(c *gin.Context) {
	// id := c.Param("id")
	// fmt.Println("user id: ", id)
	player, err := services.GetUserByToken(c)
	if err != nil {
		services.Response(c, http.StatusBadRequest, false, "Something went wrong!", err.Error())
		return
	}
	dataRes := dto.ProfileResponse{
		FirstName: player.FirstName,
		LastName:  player.LastName,
		Email:     player.Email,
		Phone:     player.Phone,
		Id:        player.Id,
		Balance:   player.Balance,
	}
	services.Response(c, http.StatusOK, true, "Token OK", dataRes)
}

func Logout(c *gin.Context) {
	player, err := services.GetUserByToken(c)
	if err != nil {
		services.Response(c, http.StatusBadRequest, false, "Something went wrong!", err.Error())
		return
	}
	redis := services.GetRedisClient()
	redis.Del(c, player.Id)
	redis.Del(c, player.Email)

	services.Response(c, http.StatusOK, true, "Logout success", nil)
}
