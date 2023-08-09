package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

func customValidation(v *validator.Validate) {
	v.RegisterValidation("valid_password", validPassword)
	v.RegisterValidation("indonesian_phone_number", indonesianPhoneNumber)
}

func validPassword(fl validator.FieldLevel) bool {
	var containString, containNumber, containSpecialCharacter bool
	password := fl.Field().Interface()

	var regex, err = regexp.Compile(`[A-Z]`)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	containString = regex.MatchString(password.(string))

	regex, err = regexp.Compile(`\d`)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	containNumber = regex.MatchString(password.(string))

	regex, err = regexp.Compile(`[^a-zA-Z\d]`)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	containSpecialCharacter = regex.MatchString(password.(string))
	return (containString && containNumber && containSpecialCharacter)
}

func indonesianPhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().Interface()
	return strings.HasPrefix(phoneNumber.(string), "+62")
}
