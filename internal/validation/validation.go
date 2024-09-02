// Package validation is core of repository
package validation

import (
	"fmt"
	"reflect"
	"regexp"
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

const iterativeRuleRegexStr = ".*\\.\\*$"

var iterativeRuleRegex = regexp.MustCompile(iterativeRuleRegexStr)

const iterativeRuleCompatibleRegexStr = "(.*)(\\..*)$"

var iterativeRuleCompatibleRegex = regexp.MustCompile(iterativeRuleCompatibleRegexStr)

func unwrapIterativeRules(validatable *Validatable) {
	newRules := make(map[string][]string)

	for ruleKey, ruleSet := range validatable.Rules {
		if !iterativeRuleRegex.MatchString(ruleKey) {
			newRules[ruleKey] = ruleSet

			continue
		}

		ruleUnprefixed := strings.TrimSuffix(ruleKey, ".*")

		for fieldKey := range validatable.JSON {
			match := iterativeRuleCompatibleRegex.FindStringSubmatch(fieldKey)
			if match == nil || len(match) <= 1 {
				continue
			}

			fieldUnprefixed := match[1]

			if ruleUnprefixed == fieldUnprefixed {
				newRules[fieldKey] = ruleSet
			}
		}
	}

	validatable.Rules = newRules
}

func camelFieldKeys(validatable *Validatable) {
	newFields := make(map[string]interface{})

	for key, value := range validatable.JSON {
		newKey := strings_thumbrise.ToCamel(key)
		newFields[newKey] = value
	}

	validatable.JSON = newFields
}

// Validate method processes validation of map by rules.
func Validate(validatable *Validatable) (*Error, error) {
	camelFieldKeys(validatable)
	unwrapIterativeRules(validatable)

	validationErrors := make(map[string]FieldValidationFail)

	for fieldKey, ruleSet := range validatable.Rules {
		fieldValue, fieldExists := validatable.JSON[fieldKey]

		// Handle empty or nil field
		if !fieldExists || fieldValue == nil {
			if slices.Contains(ruleSet, TagRequired) {
				// Else add required error
				validationErrors[fieldKey] = FieldValidationFail{
					Field: fieldKey,
					Rules: []string{TagRequired},
					Value: nil,
				}
			}

			continue
		}

		reflectedValue := reflect.ValueOf(fieldValue)

		ruleSet = slices.DeleteFunc(ruleSet, func(s string) bool {
			return s == TagRequired
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
