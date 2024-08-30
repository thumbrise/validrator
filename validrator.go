// Package validrator provides validation of golang structure, golang map, raw json.
// At the same time, it makes it possible to convert json to a golang structure.
// Provides a flexible interaction interface for the possibilities of internationalization,
// redefinition and addition of validation error messages, redefinition and addition of validation handler rules.
// Unlike most other golang validators, it makes it possible to use normal json data,
// as it would be in other programming languages in other frameworks.
package validrator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"

	"github.com/thumbrise/validrator/internal/convert"
)

var errInvalidRule = errors.New("invalid rule")

const (
	// tagOptional define type able to skip if value empty or field even does not exist.
	tagOptional = "optional"

	// tagRequired define rule which returns in validation error if optional not applied but field empty or even does not exist.
	tagRequired = "required"
)

// Validrator is main struct of package. Create via constructor.
type Validrator struct {
	handlers map[string]RuleHandlerFunc
}

// ValidationError is set of fail entries.
type ValidationError struct {
	Failed map[string]FieldValidationFail
}

// FieldValidationFail is fail entry of field.
type FieldValidationFail struct {
	Field string
	Rules []string
	Value interface{}
}

// RuleHandlerFunc is type for custom handler.
type RuleHandlerFunc func(v reflect.Value, ruleArgs []string) bool

func (v *ValidationError) Error() string {
	builder := strings.Builder{}

	for fieldName, fail := range v.Failed {
		rulesStr := strings.Join(fail.Rules, ", ")
		builder.WriteString(fmt.Sprintf("field=%s rules=%s value=%+v\n", fieldName, rulesStr, fail.Value))
	}

	return builder.String()
}

// NewValidrator constructor.
func NewValidrator() *Validrator {
	return &Validrator{
		handlers: make(map[string]RuleHandlerFunc),
	}
}

func parseRuleArgs(tag string) []string {
	colonIndex := strings.Index(tag, ":")
	if colonIndex == -1 {
		return []string{}
	}

	params := tag[colonIndex+1:]

	// Разделяем параметры по запятой
	return strings.Split(params, ",")
}

// ValidateMap method processes validation of map by rules.
func (v *Validrator) ValidateMap(data map[string]interface{}, rules map[string][]string) error {
	errs := make(map[string]FieldValidationFail)

	for fieldKey, ruleSet := range rules {
		fieldValue, fieldExists := data[fieldKey]

		reflectedValue := reflect.ValueOf(fieldValue)

		if !fieldExists || (!reflectedValue.IsValid() && fieldValue == nil) {
			// Optional is special rule
			// If field empty BUT "optional" rule applied, then no "required" error
			if !slices.Contains(ruleSet, tagOptional) {
				errs[fieldKey] = FieldValidationFail{
					Field: fieldKey,
					Rules: []string{tagRequired},
					Value: nil,
				}
			}

			continue
		}

		ruleSet = slices.DeleteFunc(ruleSet, func(s string) bool {
			return s == tagOptional
		})

		fieldErrs, err := v.validateField(reflectedValue, ruleSet)
		if err != nil {
			return err
		}

		if len(fieldErrs) > 0 {
			errs[fieldKey] = FieldValidationFail{
				Field: fieldKey,
				Rules: fieldErrs,
				Value: fieldValue,
			}
		}
	}

	if len(errs) > 0 {
		return &ValidationError{
			Failed: errs,
		}
	}

	return nil
}

func (v *Validrator) validateField(value reflect.Value, ruleSet []string) ([]string, error) {
	fieldErrs := make([]string, 0, 10)

	for _, rule := range ruleSet {
		handler, ok := v.handlers[rule]
		if !ok {
			return fieldErrs, fmt.Errorf("%w: %s", errInvalidRule, rule)
		}

		if handler(value, parseRuleArgs(rule)) {
			return fieldErrs, nil
		}

		fieldErrs = append(fieldErrs, rule)
	}

	return fieldErrs, nil
}

// ValidateJSON method processes validation of map by rules.
func (v *Validrator) ValidateJSON(input []byte, rules map[string][]string) error {
	data := make(map[string]interface{})

	err := convert.JSONToMap(input, data)
	if err != nil {
		return errors.Unwrap(err)
	}

	return v.ValidateMap(data, rules)
}

// ValidateJSONReaderToStruct method processes validation of map by rules and marshall to struct.
func (v *Validrator) ValidateJSONReaderToStruct(input io.Reader, rules map[string][]string, output any) error {
	buf := bytes.Buffer{}

	_, err := io.Copy(&buf, input)
	if err != nil {
		return errors.Unwrap(err)
	}

	return v.ValidateJSONToStruct(buf.Bytes(), rules, output)
}

// ValidateJSONToStruct method processes validation of map by rules and marshall to struct.
func (v *Validrator) ValidateJSONToStruct(input []byte, rules map[string][]string, output any) error {
	validationErrors := v.ValidateJSON(input, rules)
	if validationErrors != nil {
		return validationErrors
	}

	err := convert.JSONToStruct(input, output)
	if err != nil {
		return errors.Unwrap(err)
	}

	return validationErrors
}

// AddRuleHandler register new custom rule with handler function.
func (v *Validrator) AddRuleHandler(rule string, handlerFunc RuleHandlerFunc) {
	v.handlers[rule] = handlerFunc
}
