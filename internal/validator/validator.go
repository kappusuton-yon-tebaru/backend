package validator

import (
	"reflect"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	*validator.Validate
	translator ut.Translator
}

func New() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())

	en := en.New()
	translator, _ := ut.New(en, en).GetTranslator("en")

	validate = validator.New()
	en_translations.RegisterDefaultTranslations(validate, translator)

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})

	return &Validator{
		validate,
		translator,
	}
}

func (v *Validator) Translate(err error) []string {
	errs := []string{}
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errs = append(errs, e.Translate(v.translator))
		}
	}

	return errs
}
