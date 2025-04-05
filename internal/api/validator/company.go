package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/vcsfrl/xm/internal/model"
	"slices"
)

func CompanyValidator(logger zerolog.Logger) *validator.Validate {
	var validate = validator.New()

	// register a custom validation for company type
	err := validate.RegisterValidation("company_type", func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface().(model.CompanyType)

		return slices.Contains(model.CompanyTypes, value)
	})

	if err != nil {
		logger.Error().Err(err).Msg("failed to register custom validation for company type")
	}

	return validate
}
