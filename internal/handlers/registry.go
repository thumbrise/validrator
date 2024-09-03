package handlers

import (
	"github.com/thumbrise/validrator/internal/validation"
)

// BuiltInHandlers is the default pack of validations.
var BuiltInHandlers = map[string]validation.RuleHandlerFunc{
	"len":          HasLengthOf,
	"boolean":      IsBoolean,
	"min":          HasMinOf,
	"max":          HasMaxOf,
	"eq":           IsEq,
	"ne":           IsNe,
	"lt":           IsLt,
	"lte":          IsLte,
	"gt":           IsGt,
	"gte":          IsGte,
	"jwt":          IsJWT,
	"alphaunicode": IsAlphaUnicode,
	"datetime":     IsDatetime,
	"number":       IsNumber,
	"email":        IsEmail,
	"url":          IsURL,
	"http_url":     IsHttpURL,
	"uri":          IsURI,
	"contains":     Contains,
	"oneof":        IsOneOf,
	// "eqfield":  isEqField,
	// "nefield":  isNeField,
	// "gtefield": isGteField,
	// "gtfield":  IsGtField,
	// "ltefield": isLteField,
	// "ltfield":  isLtField,
}
