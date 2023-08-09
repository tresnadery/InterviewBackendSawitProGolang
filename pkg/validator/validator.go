package validator

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"strings"
)

func Validate(i interface{}) error {
	errors := map[string][]string{}
	v, trans := newValidator()
	err := v.Struct(i)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := strings.ToLower(string(err.Field()[0])) + err.Field()[1:]
			errors[fieldName] = append(errors[fieldName], err.Translate(trans))
		}
		errorsMarshal, err := json.Marshal(errors)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(errorsMarshal))
	}
	return nil
}

func newValidator() (*validator.Validate, ut.Translator) {
	v := validator.New()
	trans := newTranslation(v)
	customValidation(v)
	customTranslation(v, trans)
	return v, trans
}

func newTranslation(v *validator.Validate) ut.Translator {
	var uni *ut.UniversalTranslator
	var trans ut.Translator
	en := en.New()
	uni = ut.New(en, en)

	trans, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v, trans)
	return trans
}
