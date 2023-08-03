package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type RequestValidator struct {
	v *validator.Validate
}

func (rv *RequestValidator) Validate(v interface{}) error {
	if e := rv.v.Struct(v); e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, e.Error())
	}
	return nil
}

func New() *RequestValidator {
	return &RequestValidator{v: validator.New()}
}
