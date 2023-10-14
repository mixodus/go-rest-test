package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mixodus/go-rest-test/config"
	"github.com/mixodus/go-rest-test/dto"
	"github.com/mixodus/go-rest-test/services"
)

func Authenticate(c *gin.Context) {
	// fmt.Println(c.Request.Header.Get("Authorization"))

	//get token from header
	token := c.Request.Header.Get("Authorization")
	//split token from Bearer and check if jwt exist
	if strings.Split(token, " ")[0] == "Bearer" {
		token = strings.Split(token, " ")[1]
	} else {
		//return unauthorized
		res := dto.Response{
			Status:  false,
			Message: "Unauthorized",
			Data:    nil,
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}

	//validate token
	claims := &config.JWTClaims{}
	tokenz, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	if err != nil {
		//return unauthorized
		v, _ := err.(*jwt.ValidationError)
		switch v.Errors {
		case jwt.ValidationErrorSignatureInvalid:
			fmt.Println("Error: Signature Invalid")
			res := dto.Response{
				Status:  false,
				Message: "Error: Signature Invalid",
				Data:    nil,
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		case jwt.ValidationErrorExpired:
			fmt.Println("Error: Token Expired")
			res := dto.Response{
				Status:  false,
				Message: "Token Expired",
				Data:    nil,
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		default:
			fmt.Println("Error: Can't handle this token")
			res := dto.Response{
				Status:  false,
				Message: "unauthorized",
				Data:    nil,
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
	}

	if !tokenz.Valid {
		res := dto.Response{
			Status:  false,
			Message: "unauthorized",
			Data:    nil,
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}

	// ==== SESSION FROM REDIS ====

	//check if token exist in redis
	redis := services.GetRedisClient()
	storedTokenz := redis.Get(c, claims.Id).Val()
	if storedTokenz == "" {
		res := dto.Response{
			Status:  false,
			Message: "unauthorized",
			Data:    nil,
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}
	//check if token is current active token
	if !(token == storedTokenz) {
		res := dto.Response{
			Status:  false,
			Message: "unauthorized",
			Data:    nil,
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}
	// ==== end ====

	//==== SESSION FROM DATABASE ====
	//!!UNUSED since we use redis for session management!!
	//==== check if token exist in db ====
	// var tokenSession models.TokenSession
	// if err := models.DB.Where("player_id = ?", claims.Id).First(&tokenSession).Error; err != nil {
	// 	res := dto.Response{
	// 		Status:  false,
	// 		Message: "unauthorized",
	// 		Data:    nil,
	// 	}
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, res)
	// 	return
	// }
	//==== end ===

	//==== if token is not current active token then throw error ====
	// if !(token == tokenSession.Token) {
	// 	res := dto.Response{
	// 		Status:  false,
	// 		Message: "unauthorized",
	// 		Data:    nil,
	// 	}
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, res)
	// 	return
	// }
	//==== end ===

	// fmt.Println("User ID: ", claims.Id)
	// fmt.Println("User Email: ", claims.Email)

	// c.Params = append(c.Params, gin.Param{Key: "id", Value: claims.Id})

	c.Next()
}
