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

func FilterPlayer(c *gin.Context) (interface{}, error) {
	firstName := c.DefaultQuery("first_name", "")
	lastName := c.DefaultQuery("last_name", "")
	email := c.DefaultQuery("email", "")
	phone := c.DefaultQuery("phone", "")
	balance_bigger_than := c.DefaultQuery("balance_bigger_than", "")
	balance_less_than := c.DefaultQuery("balance_less_than", "")
	account_name := c.DefaultQuery("account_name", "")
	account_number := c.DefaultQuery("account_number", "")
	player_created_at_after := c.DefaultQuery("player_created_at_after", "")
	player_created_at_before := c.DefaultQuery("player_created_at_before", "")

	var player []dto.Player
	// models.DB.Joins("PlayersBank").Find(&player)
	query := models.DB.Table("players").
		Joins("LEFT JOIN players_banks ON players.id = players_banks.player_id").
		Joins("LEFT JOIN banks ON players_banks.bank_id = banks.id")
	if firstName != "" {
		lowerCase, titleCase := LowerCaseTitleCase(firstName)
		query = query.Where("players.first_name LIKE ?", "%"+lowerCase+"%").
			Or("players.first_name LIKE ?", "%"+titleCase+"%")
	}
	if lastName != "" {
		lowerCase, titleCase := LowerCaseTitleCase(lastName)
		query = query.Where("players.last_name LIKE ?", "%"+lowerCase+"%").
			Or("players.last_name LIKE ?", "%"+titleCase+"%")
	}
	if email != "" {
		lowerCaseEmail := ToLowerCase(email)
		query = query.Where("players.email LIKE ?", "%"+lowerCaseEmail+"%")
	}
	if phone != "" {
		query = query.Where("players.phone LIKE ?", "%"+phone+"%")
	}
	if balance_bigger_than != "" {
		query = query.Where("players.balance >= ?", balance_bigger_than)
	}
	if balance_less_than != "" {
		query = query.Where("players.balance <= ?", balance_less_than)
	}
	if account_name != "" {
		lowerCase, titleCase := LowerCaseTitleCase(firstName)
		query = query.Where("players_banks.bank_account_name LIKE ?", "%"+lowerCase+"%").
			Or("players_banks.bank_account_name LIKE ?", "%"+titleCase+"%")
	}
	if account_number != "" {
		query = query.Where("players_banks.bank_account_number LIKE ?", "%"+account_number+"%")
	}
	fmt.Println("player_created_at_after: " + player_created_at_after)
	if player_created_at_after != "" {
		created_at_after, err := time.Parse("2006-01-02", player_created_at_after)
		if err != nil {
			return nil, err
		}
		fmt.Println(created_at_after)
		query = query.Where("players.created_at >= ?", created_at_after)
	}
	fmt.Println("player_created_at_before: " + player_created_at_before)
	if player_created_at_before != "" {
		created_at_before, err := time.Parse("2006-01-02", player_created_at_before)
		if err != nil {
			return nil, err
		}
		fmt.Println(created_at_before)
		query = query.Where("players.created_at <= ?", created_at_before)
	}

	query.Preload("PlayersBank.Bank").Find(&player)

	return player, nil
}
