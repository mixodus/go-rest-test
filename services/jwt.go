package services

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mixodus/go-rest-test/config"
	"github.com/mixodus/go-rest-test/models"
)

func GeneratePlayerToken(c *gin.Context, player *models.Player) (string, error) {
	//create token
	expTime := time.Now().Add(time.Hour * 24)
	calims := &config.JWTClaims{
		Id:    player.Id,
		Email: player.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Player",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	//declare token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, calims)
	//sign token
	tokenz, err := token.SignedString(config.JWT_KEY)
	if err != nil {
		return "", err
	}
	return tokenz, nil
}
