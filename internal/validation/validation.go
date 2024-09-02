// Package validation is core of repository
package validation

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	strings_thumbrise "github.com/thumbrise/validrator/internal/strings"
)

// RuleHandlerFunc is type for custom handler.
type RuleHandlerFunc func(v reflect.Value, ruleArgs []string) bool

// Validatable is type for input of validation.
type Validatable struct {
	JSON     map[string]interface{}
	Rules    map[string][]string
	Handlers map[string]RuleHandlerFunc
}

// Validate method processes validation of map by rules.
func Validate(validatable *Validatable) (*Error, error) {
	validationErrors := make(map[string]FieldValidationFail)

	for fieldKey, ruleSet := range validatable.Rules {
		fieldKey = strings_thumbrise.ToCamel(fieldKey)
		fieldValue, fieldExists := validatable.JSON[fieldKey]

		reflectedValue := reflect.ValueOf(fieldValue)

		if !fieldExists || (!reflectedValue.IsValid() && fieldValue == nil) {
			// Optional is special rule
			// If field empty BUT "optional" rule applied, then no "required" error
			if !slices.Contains(ruleSet, tagOptional) {
				validationErrors[fieldKey] = FieldValidationFail{
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

		// Handle nested rules

		fieldErrs, err := validateField(reflectedValue, ruleSet, validatable.Handlers)
		if err != nil {
			return nil, err
		}

		if len(fieldErrs) > 0 {
			validationErrors[fieldKey] = FieldValidationFail{
				Field: fieldKey,
				Rules: fieldErrs,
				Value: fieldValue,
			}
		}
	}

	if len(validationErrors) > 0 {
		return &Error{
			Failed: validationErrors,
		}, nil
	}

	return nil, nil //nolint:nilnil
}

func validateField(value reflect.Value, ruleSet []string, handlers map[string]RuleHandlerFunc) ([]string, error) {
	fieldErrs := make([]string, 0, 10)

	for _, rule := range ruleSet {
		handler, ok := handlers[rule]
		if !ok {
			return fieldErrs, fmt.Errorf("%w: %s", errInvalidRule, rule)
		}

		ruleArgs := parseRuleArgs(rule)
		if handler(value, ruleArgs) {
			return fieldErrs, nil
		}

		fieldErrs = append(fieldErrs, rule)
	}

	return fieldErrs, nil
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
