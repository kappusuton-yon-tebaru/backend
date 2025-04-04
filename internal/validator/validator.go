package validator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	*validator.Validate
	translator ut.Translator
}

func New() (*Validator, error) {
	en := en.New()
	translator, _ := ut.New(en, en).GetTranslator("en")

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := en_translations.RegisterDefaultTranslations(validate, translator)
	if err != nil {
		return nil, err
	}

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := field.Tag.Get("json")

		if len(strings.TrimSpace(tag)) == 0 {
			tag = field.Tag.Get("form")
		}

		return tag
	})

	err = validate.RegisterValidation("kebabnum", func(fl validator.FieldLevel) bool {
		pattern := "^[a-z0-9]+(-[a-z0-9]+)*$"
		return regexp.MustCompile(pattern).MatchString(fl.Field().String())
	})
	if err != nil {
		return nil, err
	}

	return &Validator{
		validate,
		translator,
	}, nil
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
