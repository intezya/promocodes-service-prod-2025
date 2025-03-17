package validator

import "github.com/go-playground/validator/v10"

func New() *validator.Validate {
	val := validator.New()
	_ = val.RegisterValidation("country", validateCountryCode)
	_ = val.RegisterValidation("url", validateURL)
	_ = val.RegisterValidation("mode_logic", validateModeLogic)
	_ = val.RegisterValidation("business_name", validateBusinessName)
	_ = val.RegisterValidation("password", validatePassword)
	_ = val.RegisterValidation("name", validateName)
	_ = val.RegisterValidation("description", validateDescription)
	return val
}
