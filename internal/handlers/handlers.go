// Package handlers contains builtin validation handlers
package handlers

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// HasLengthOf is the validation function for validating if the current field's value is equal to the param's value.
func HasLengthOf(field reflect.Value, params []string) bool {
	if len(params) < 1 {
		return false
	}

	param := params[0]

	switch field.Kind() { //nolint:exhaustive
	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) == p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() == p

	case reflect.Float32:
		p := asFloat32(param)

		return field.Float() == p

	case reflect.Float64:
		p := asFloat64(param)

		return field.Float() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsLt is the validation function for validating if the current field's value is less than the param's value.
func IsLt(field reflect.Value, params []string) bool { //nolint:dupl
	if len(params) < 1 {
		return false
	}

	param := params[0]

	switch field.Kind() { //nolint:exhaustive
	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) < p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) < p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)

		return field.Int() < p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() < p

	case reflect.Float32:
		p := asFloat32(param)

		return field.Float() < p

	case reflect.Float64:
		p := asFloat64(param)

		return field.Float() < p

	case reflect.Struct:
		if field.Type().ConvertibleTo(timeType) {
			return field.Convert(timeType).Interface().(time.Time).Before(time.Now().UTC()) //nolint:forcetypeassert
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsGt is the validation function for validating if the current field's value is greater than the param's value.
func IsGt(field reflect.Value, params []string) bool { //nolint:dupl
	if len(params) < 1 {
		return false
	}

	param := params[0]

	switch field.Kind() { //nolint:exhaustive
	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) > p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) > p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)

		return field.Int() > p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() > p

	case reflect.Float32:
		p := asFloat32(param)

		return field.Float() > p

	case reflect.Float64:
		p := asFloat64(param)

		return field.Float() > p

	case reflect.Struct:
		if field.Type().ConvertibleTo(timeType) {
			return field.Convert(timeType).Interface().(time.Time).After(time.Now().UTC()) //nolint:forcetypeassert
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsHttpURL is the validation function for validating if the current field's value is a valid HTTP(s) URL.
func IsHttpURL(field reflect.Value, params []string) bool { //nolint:revive,stylecheck
	if !IsURL(field, params) {
		return false
	}

	switch field.Kind() { //nolint:gocritic,exhaustive
	case reflect.String:
		s := strings.ToLower(field.String())

		url, err := url.Parse(s)
		if err != nil || url.Host == "" {
			return false
		}

		return url.Scheme == "http" || url.Scheme == "https"
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsURI is the validation function for validating if the current field's value is a valid URI.
func IsURI(field reflect.Value, _ []string) bool {
	switch field.Kind() { //nolint:gocritic,exhaustive
	case reflect.String:
		str := field.String()

		// checks needed as of Go 1.6 because of change https://github.com/golang/go/commit/617c93ce740c3c3cc28cdd1a0d712be183d0b328#diff-6c2d018290e298803c0c9419d8739885L195
		// emulate browser and strip the '#' suffix prior to validation. see issue-#237
		if i := strings.Index(str, "#"); i > -1 {
			str = str[:i]
		}

		if len(str) == 0 {
			return false
		}

		_, err := url.ParseRequestURI(str)

		return err == nil
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// Contains is the validation function for validating that the field's value Contains the text specified within the param.
func Contains(field reflect.Value, params []string) bool {
	if len(params) < 1 {
		return false
	}

	param := params[0]

	return strings.Contains(field.String(), param)
}

// IsOneOf godoc.
func IsOneOf(field reflect.Value, params []string) bool {
	var val string

	switch field.Kind() { //nolint:exhaustive
	case reflect.String:
		val = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = strconv.FormatUint(field.Uint(), 10)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}

	for _, param := range params {
		if param == val {
			return true
		}
	}

	return false
}

// isFileURL is the helper function for validating if the `path` valid file URL as per RFC8089.
func isFileURL(path string) bool {
	if !strings.HasPrefix(path, "file:/") {
		return false
	}

	_, err := url.ParseRequestURI(path)

	return err == nil
}

// IsNumber is the validation function for validating if the current field's value is a valid number.
func IsNumber(field reflect.Value, _ []string) bool {
	switch field.Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		return true
	default:
		return numberRegex().MatchString(field.String())
	}
}

// IsBoolean is the validation function for validating if the current field's value is a valid boolean value or can be safely converted to a boolean value.
func IsBoolean(field reflect.Value, _ []string) bool {
	switch field.Kind() { //nolint:exhaustive
	case reflect.Bool:
		return true
	default:
		_, err := strconv.ParseBool(field.String())

		return err == nil
	}
}

// IsURL is the validation function for validating if the current field's value is a valid URL.
func IsURL(field reflect.Value, _ []string) bool {
	switch field.Kind() { //nolint:gocritic,exhaustive
	case reflect.String:
		str := strings.ToLower(field.String())

		if len(str) == 0 {
			return false
		}

		if isFileURL(str) {
			return true
		}

		url, err := url.Parse(str)
		if err != nil || url.Scheme == "" {
			return false
		}

		if url.Host == "" && url.Fragment == "" && url.Opaque == "" {
			return false
		}

		return true
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsEmail is the validation function for validating if the current field's value is a valid email address.
func IsEmail(field reflect.Value, _ []string) bool {
	return emailRegex().MatchString(field.String())
}

// IsAlphaUnicode is the validation function for validating if the current field's value is a valid alpha unicode value.
func IsAlphaUnicode(field reflect.Value, _ []string) bool {
	return alphaUnicodeRegex().MatchString(field.String())
}

// HasMinOf is the validation function for validating if the current field's value is greater than or equal to the param's value.
func HasMinOf(fl reflect.Value, params []string) bool {
	return IsGte(fl, params)
}

// HasMaxOf is the validation function for validating if the current field's value is less than or equal to the param's value.
func HasMaxOf(fl reflect.Value, params []string) bool {
	return IsLte(fl, params)
}

// IsDatetime is the validation function for validating if the current field's value is a valid datetime string.
func IsDatetime(field reflect.Value, params []string) bool {
	if len(params) < 1 {
		return false
	}

	param := params[0]

	if field.Kind() == reflect.String {
		_, err := time.Parse(param, field.String())

		return err == nil
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsJWT is the validation function for validating if the current field's value is a valid JWT string.
func IsJWT(field reflect.Value, _ []string) bool {
	return jWTRegex().MatchString(field.String())
}

// Bool validate false, true, 1, 0, "true", "false", "0", "1".
func Bool(v reflect.Value, _ []string) bool {
	switch v.Interface() {
	case 1, 0, false, true, "true", "false", "0", "1":
		return true
	}

	return false
}

// Len validate length for next types: String, Slice, Map, Array, Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Uintptr, Float32, Float64.
func Len(val reflect.Value, args []string) bool {
	if len(args) < 1 {
		panic("Len expects 1 argument")
	}

	param := args[0]

	switch val.Kind() { //nolint:exhaustive
	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(val.String())) == p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(val.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(val.Type(), param)

		return val.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return val.Uint() == p

	case reflect.Float32:
		p := asFloat32(param)

		return val.Float() == p

	case reflect.Float64:
		p := asFloat64(param)

		return val.Float() == p
	default:
		panic(fmt.Sprintf("Bad field type %T", val.Interface()))
	}
}

// IsGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func IsGte(field reflect.Value, params []string) bool { //nolint:cyclop
	if len(params) < 1 {
		return false
	}

	param := params[0]

	switch field.Kind() { //nolint:exhaustive
	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) >= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) >= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)

		return field.Int() >= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() >= p

	case reflect.Float32:
		p := asFloat32(param)

		return field.Float() >= p

	case reflect.Float64:
		p := asFloat64(param)

		return field.Float() >= p

	case reflect.Struct:
		if field.Type().ConvertibleTo(timeType) {
			now := time.Now().UTC()
			t := field.Convert(timeType).Interface().(time.Time) //nolint:forcetypeassert

			return t.After(now) || t.Equal(now)
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func IsLte(field reflect.Value, params []string) bool { //nolint:cyclop
	if len(params) < 1 {
		return false
	}

	param := params[0]

	switch field.Kind() { //nolint:exhaustive
	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) <= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) <= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)

		return field.Int() <= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() <= p

	case reflect.Float32:
		p := asFloat32(param)

		return field.Float() <= p

	case reflect.Float64:
		p := asFloat64(param)

		return field.Float() <= p

	case reflect.Struct:
		if field.Type().ConvertibleTo(timeType) {
			now := time.Now().UTC()

			t, ok := field.Convert(timeType).Interface().(time.Time)
			if ok {
				return t.Before(now) || t.Equal(now)
			}
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsNe is the validation function for validating that the field's value does not equal the provided param value.
func IsNe(field reflect.Value, params []string) bool {
	return !IsEq(field, params)
}

// IsEq is the validation function for validating if the current field's value is equal to the param's value.
func IsEq(field reflect.Value, params []string) bool {
	if len(params) < 1 {
		return false
	}

	param := params[0]

	switch field.Kind() { //nolint:exhaustive
	case reflect.String:
		return field.String() == param

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asIntFromType(field.Type(), param)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() == p

	case reflect.Float32:
		p := asFloat32(param)

		return field.Float() == p

	case reflect.Float64:
		p := asFloat64(param)

		return field.Float() == p

	case reflect.Bool:
		p := asBool(param)

		return field.Bool() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}
