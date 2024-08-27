// Package validrator provides validation of golang structure, golang map, raw json.
// At the same time, it makes it possible to convert json to a golang structure.
// Provides a flexible interaction interface for the possibilities of internationalization,
// redefinition and addition of validation error messages, redefinition and addition of validation handler rules.
// Unlike most other golang validators, it makes it possible to use normal json data,
// as it would be in other programming languages in other frameworks.
package validrator

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var errInvalidRule = errors.New("invalid rule")

const tagOptional = "optional"

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
type RuleHandlerFunc func(v interface{}, ruleArgs string) bool

func (v *ValidationError) Error() string {
	builder := strings.Builder{}

	for fieldName, fail := range v.Failed {
		rulesStr := strings.Join(fail.Rules, ", ")
		builder.WriteString(fmt.Sprintf("field=%s\nrules=%s\nvalue=%+v\n", fieldName, rulesStr, fail.Value))
	}

	return builder.String()
}

// NewValidrator constructor.
func NewValidrator() *Validrator {
	return &Validrator{
		handlers: make(map[string]RuleHandlerFunc),
	}
}

func parseRuleArgs(tag string) string {
	colonIndex := strings.Index(tag, ":")
	if colonIndex == -1 {
		return ""
	}

	return tag[colonIndex+1:]
}

// Validate method processes validation of map by rules.
func (v *Validrator) Validate(data map[string]interface{}, rules map[string][]string) error {
	errs := make(map[string]FieldValidationFail)

	for fieldKey, ruleSet := range rules {
		fieldValue, fieldExists := data[fieldKey]

		if !fieldExists && slices.Contains(ruleSet, tagOptional) {
			continue
		}

		ruleSet = slices.DeleteFunc(ruleSet, func(s string) bool {
			return s == tagOptional
		})

		fieldErrs := make([]string, 0, 10)

		for _, rule := range ruleSet {
			handler, ok := v.handlers[rule]
			if !ok {
				return fmt.Errorf("%w: %s", errInvalidRule, rule)
			}

			if handler(fieldValue, parseRuleArgs(rule)) {
				continue
			}

			fieldErrs = append(fieldErrs, rule)
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

// AddRuleHandler register new custom rule with handler function.
func (v *Validrator) AddRuleHandler(rule string, handlerFunc RuleHandlerFunc) {
	v.handlers[rule] = handlerFunc
}
