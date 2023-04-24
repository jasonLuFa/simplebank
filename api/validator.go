package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/jasonLuFa/simplebank/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// 方法一
	// return util.IsSupportedCurrency(fieldLevel.Field().String())
	// 方法二
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

var validAmount validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if amountString, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsAmountPositive(amountString)
	}
	return false
}
