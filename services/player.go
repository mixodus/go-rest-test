package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mixodus/go-rest-test/config"
	"github.com/mixodus/go-rest-test/dto"
	"github.com/mixodus/go-rest-test/models"
	"golang.org/x/crypto/bcrypt"
)

func GetUserByToken(c *gin.Context) (models.Player, error) {
	//get user by token
	token := c.Request.Header.Get("Authorization")
	token = strings.Split(token, " ")[1]
	//split token from Bearer and check if jwt exist
	//validate token
	claims := &config.JWTClaims{}
	tokenz, _ := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})
	fmt.Println(tokenz)
	var player models.Player
	if err := models.DB.Where("id = ?", claims.Id).First(&player).Error; err != nil {
		return player, err
	}
	return player, nil
}

func LoginViaRedis(c *gin.Context, userInput *dto.LoginRequest, session string) (interface{}, error) {
	var userData models.Player
	json.Unmarshal([]byte(session), &userData)
	if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(userInput.Password)); err != nil {
		return nil, err
	} else {
		tokenz, err := GeneratePlayerToken(c, &userData)
		if err != nil {
			return nil, err
		}

		//set dto
		data := map[string]interface{}{
			"token": tokenz,
		}

		//store token to redis
		redis := GetRedisClient()
		redis.Set(c, userData.Id, tokenz, time.Hour*24)

		return data, nil
	}
}
