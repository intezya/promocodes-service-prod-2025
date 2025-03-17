package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

func validateCountryCode(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) == 2
}

func validateURL(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^(http|https)://`)
	if len(fl.Field().String()) > 350 {
		return false
	}
	return re.MatchString(fl.Field().String())
}

func validateModeLogic(fl validator.FieldLevel) bool {
	parent := fl.Parent()
	modeField := fl.Field()
	tag := fl.Param()

	fields := strings.Split(tag, ",")
	if len(fields) != 2 {
		return false
	}

	promoCommon := parent.FieldByName(fields[0])
	promoUnique := parent.FieldByName(fields[1])
	maxCountField := parent.FieldByName("MaxCount")

	var promoCommonValue string
	if !promoCommon.IsNil() {
		promoCommonValue = promoCommon.Elem().String()
	}

	var promoUniqueValue []string
	if !promoUnique.IsNil() {
		promoUniqueValue = promoUnique.Elem().Interface().([]string)
	}

	var maxCountValue int
	if !maxCountField.IsNil() {
		maxCountValue = int(maxCountField.Int())
	}

	switch modeField.String() {
	case "COMMON":
		return promoCommonValue != "" && len(promoUniqueValue) == 0 && maxCountValue >= 0
	case "UNIQUE":
		return promoCommonValue == "" && len(promoUniqueValue) > 0
	default:
		return false
	}
}

func validateBusinessName(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if len(v) < 5 || len(v) > 50 {
		return false
	}
	return true
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasLowercase := false
	hasUppercase := false
	hasDigit := false
	hasSpecialChar := false

	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLowercase = true
		case 'A' <= char && char <= 'Z':
			hasUppercase = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case containsSpecialChar(char):
			hasSpecialChar = true
		}
	}

	return regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,60}$`).MatchString(password) &&
		hasLowercase &&
		hasUppercase &&
		hasDigit &&
		hasSpecialChar
}

func containsSpecialChar(char rune) bool {
	specialChars := "@$!%*?&"
	for _, special := range specialChars {
		if char == special {
			return true
		}
	}
	return false
}

func validateName(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if len(v) < 3 || len(v) > 200 {
		return false
	}
	return true
}

func validateDescription(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if len(v) < 10 || len(v) > 300 {
		return false
	}
	return true
}
