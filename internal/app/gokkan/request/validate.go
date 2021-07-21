package request

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// Validator is the default request struct validator
// see github.com/go-playground/validator for validate tags.
type Validator struct {
	v *validator.Validate
}

// NewValidator creates a new validator.
func NewValidator() *Validator {
	return &Validator{v: validator.New()}
}

// Validate is our implementation of echo.Validator
// it is called for every request by echo.
func (v *Validator) Validate(i interface{}) error {
	if err := v.v.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
