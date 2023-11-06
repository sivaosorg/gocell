package utils

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ErrorFieldModel struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ParseError(errs ...error) []interface{} {
	var out []interface{}
	for _, err := range errs {
		switch typedError := any(err).(type) {
		case validator.ValidationErrors:
			for _, e := range typedError {
				out = append(out, ParseFieldsError(e.Field(), e))
			}
		case *json.UnmarshalTypeError:
			out = append(out, ParseMarshallingError(*typedError))
		default:
			out = append(out, err.Error())
		}
	}
	return out
}

func ParseFieldsError(field string, fe validator.FieldError) *ErrorFieldModel {

	switch fe.Tag() {
	case "required":
		return &ErrorFieldModel{field, "This field is required"}
	case "lte":
		return &ErrorFieldModel{field, fmt.Sprintf("Should be less than or equal to %v", fe.Param())}
	case "gte":
		return &ErrorFieldModel{field, fmt.Sprintf("Should be greater than or equal to %v", fe.Param())}
	default:
		english := en.New()
		translator := ut.New(english, english)
		if translatorInstance, found := translator.GetTranslator("en"); found {
			return &ErrorFieldModel{field, fe.Translate(translatorInstance)}
		} else {
			return &ErrorFieldModel{field, fmt.Errorf("%v", fe).Error()}
		}
	}
}

func ParseMarshallingError(e json.UnmarshalTypeError) *ErrorFieldModel {
	return &ErrorFieldModel{e.Field, fmt.Sprintf("The type must be a %s", e.Type.String())}
}
