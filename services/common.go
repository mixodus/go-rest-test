package services

import (
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
