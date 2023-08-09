package validator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func customTranslation(v *validator.Validate, trans ut.Translator) {
	translateValidPassword(v, trans)
}

func translateValidPassword(v *validator.Validate, trans ut.Translator) {
	v.RegisterTranslation("valid_password", trans, func(ut ut.Translator) error {
		return ut.Add("valid_password", "Password must be contain at least 1 uppercase character, 1 special character, and 1 number", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_password", fe.Field())

		return t
	})
}

func translateIndonesianPhoneNumber(v *validator.Validate, trans ut.Translator) {
	v.RegisterTranslation("indonesian_phone_number", trans, func(ut ut.Translator) error {
		return ut.Add("indonesian_phone_number", `Phone numbers must start with the Indonesia country code  "+62"`, true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("indonesian_phone_number", fe.Field())

		return t
	})
}
