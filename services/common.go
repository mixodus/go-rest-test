package services

import (
	"github.com/gin-gonic/gin"
	"github.com/mixodus/go-rest-test/dto"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func LowerCaseTitleCase(str string) (lower string, title string) {
	casest := cases.Title(language.English)
	titleCase := casest.String(str)
	casesl := cases.Lower(language.English)
	lowerCase := casesl.String(str)
	return lowerCase, titleCase
}

func ToLowerCase(str string) string {
	casesl := cases.Lower(language.English)
	lowerCase := casesl.String(str)
	return lowerCase
}

func Response(c *gin.Context, code int, status bool, message string, data interface{}) {
	res := dto.Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	if code >= 200 && code < 400 {
		c.JSON(code, res)
	} else {
		c.AbortWithStatusJSON(code, res)
	}
}
