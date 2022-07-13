package api

import (
	"github.com/ariandi/ppob_go/util"
	"github.com/go-playground/validator/v10"
)

var validStatus validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if status, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedStatus(status)
	}
	return false
}
