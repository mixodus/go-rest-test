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
	joinCondition := "players.id = players_banks.player_id"
	var player []models.Player
	models.DB.Joins("JOIN players_banks ON " + joinCondition).Find(&player)

	query := models.DB.Dialector.Explain("SELECT * FROM players JOIN players_banks ON " + joinCondition)
	fmt.Println(query)

	res := dto.Response{
		Status:  true,
		Message: "Success",
		Data:    player,
	}

	c.JSON(http.StatusOK, res)
}

func GetPlayerById(c *gin.Context) {
	id := c.Param("id")
	var player dto.Player
	if err := models.DB.Where("id = ?", id).First(&player).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			res := dto.Response{
				Status:  false,
				Message: "User not found",
				Data:    err.Error(),
			}
			c.JSON(http.StatusBadRequest, res)
			return
		default:
			res := dto.Response{
				Status:  false,
				Message: "Internal Server Error",
				Data:    err.Error(),
			}
			c.JSON(http.StatusBadRequest, res)
			return
		}
	}

	res := dto.Response{
		Status:  true,
		Message: "Get player by ID",
		Data:    player,
	}

	c.JSON(http.StatusOK, res)
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
			res := dto.Response{
				Status:  false,
				Message: out[0].Field + ", " + out[0].Message,
				Data:    nil,
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, res)
		}

		return
	}

	//HASH PASSWORD
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	//SET DATA TO STORE FROM INPUT
	dataPost := models.Player{
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Password:  userInput.Password,
		Email:     userInput.Email,
		Phone:     userInput.Phone,
	}

	//SAVE TO DB
	if err := models.DB.Create(&dataPost).Error; err != nil {
		log.Default().Println(err.Error())
		res := dto.Response{
			Status:  false,
			Message: "User already exist",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, res)
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
	res := dto.Response{
		Status:  true,
		Message: "Success",
		Data:    dataRes,
	}

	c.JSON(http.StatusOK, res)
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
			res := dto.Response{
				Status:  false,
				Message: out[0].Field + ", " + out[0].Message,
				Data:    nil,
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, res)
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
			res := dto.Response{
				Status:  false,
				Message: "Something went wrong!",
				Data:    err.Error(),
			}
			c.JSON(http.StatusBadRequest, res)
			return
		} else {
			dataRes := dto.Response{
				Status:  true,
				Message: "Successfully login.",
				Data:    res,
			}
			c.JSON(http.StatusOK, dataRes)
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
			res := dto.Response{
				Status:  false,
				Message: "User not found",
				Data:    err.Error(),
			}
			c.JSON(http.StatusBadRequest, res)
			return
		default:
			res := dto.Response{
				Status:  false,
				Message: "Internal Server Error",
				Data:    err.Error(),
			}
			c.JSON(http.StatusBadRequest, res)
			return
		}

	}

	//CHECK PASSWORD
	if err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(userInput.Password)); err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Wrong password",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	tokenz, err := services.GeneratePlayerToken(c, &player)
	if err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Something went wrong!",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, res)
	}

	//==== SET SESSION REDIS ====
	//set session to redis
	playerJson, _ := json.Marshal(player)
	//set user email as user's data key
	redis.Set(c, userInput.Email, playerJson, time.Hour*1)
	//set player id as token key
	redis.Set(c, player.Id, tokenz, time.Hour*24)
	// if errors := redis.Set(c, userInput.Email, playerJson, time.Hour*1).Err().Error(); errors != "" {
	// 	log.Fatal(errors) //somehow this just cause error
	// }
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

	dataRes := dto.Response{
		Status:  true,
		Message: "Successfully login.",
		Data:    data,
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, dataRes)
}

func Profile(c *gin.Context) {
	// id := c.Param("id")
	// fmt.Println("user id: ", id)
	player, err := services.GetUserByToken(c)
	if err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Something went wrong!",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, res)
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
	res := dto.Response{
		Status:  true,
		Message: "Token OK",
		Data:    dataRes,
	}
	c.JSON(http.StatusOK, res)
}

func Logout(c *gin.Context) {
	player, err := services.GetUserByToken(c)
	if err != nil {
		res := dto.Response{
			Status:  false,
			Message: "Something went wrong!",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}
	redis := services.GetRedisClient()
	redis.Del(c, player.Id)
	redis.Del(c, player.Email)
	res := dto.Response{
		Status:  true,
		Message: "Logout success",
		Data:    nil,
	}
	c.JSON(http.StatusOK, res)
}
