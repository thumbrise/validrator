// Package validrator provides validation of golang structure, golang map, raw json.
// At the same time, it makes it possible to convert json to a golang structure.
// Provides a flexible interaction interface for the possibilities of internationalization,
// redefinition and addition of validation error messages, redefinition and addition of validation handler handlers.
// Unlike most other golang validators, it makes it possible to use normal json data,
// as it would be in other programming languages in other frameworks.
package validrator

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/thumbrise/validrator/internal/dot"
	"github.com/thumbrise/validrator/internal/meta"
	"github.com/thumbrise/validrator/internal/validation"
)

var inBuiltHandlers = map[string]validation.RuleHandlerFunc{
	validation.TagRequired: func(_ reflect.Value, _ []string) bool {
		return true
	},
}
var errInvalidJSON = errors.New("invalid json")

// Validrator is main struct of package. Create via constructor.
type Validrator struct {
	handlers map[string]validation.RuleHandlerFunc
}

// NewValidrator constructor.
func NewValidrator() *Validrator {
	r := &Validrator{
		handlers: make(map[string]validation.RuleHandlerFunc),
	}
	r.AddRuleHandlers(inBuiltHandlers)

	return r
}

// Validate method processes validation by structure tags and marshall to that struct.
func (v *Validrator) Validate(input []byte, output any) (*validation.Error, error) {
	if !json.Valid(input) {
		return nil, errInvalidJSON
	}

	// Preparing validation. Need handlers map and jsonInput map
	rules := collectRules(output)

	jsonInput, err := collectJSONMap(input)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	validationErrors, err := validateReal(jsonInput, rules, v.handlers)
	if validationErrors != nil || err != nil {
		return validationErrors, errors.Unwrap(err)
	}

	// Mapping to struct
	err = jsonToStruct(input, output)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	return nil, nil //nolint:nilnil
}

// AddRuleHandler register new custom rule with handler function.
func (v *Validrator) AddRuleHandler(rule string, handlerFunc validation.RuleHandlerFunc) {
	v.handlers[rule] = handlerFunc
}

// AddRuleHandlers register new custom rules with handler functions.
func (v *Validrator) AddRuleHandlers(handlers map[string]validation.RuleHandlerFunc) {
	for rule, handlerFunc := range handlers {
		v.AddRuleHandler(rule, handlerFunc)
	}
}

func collectJSONMap(input []byte) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})

	err := jsonToMap(input, jsonMap)
	if err != nil {
		return nil, err
	}

	result := dot.Map(jsonMap)

	return result, nil
}

// ValidateJSON method processes validation of map by handlers.
func validateReal(data map[string]interface{}, rules map[string][]string, handlers map[string]validation.RuleHandlerFunc) (*validation.Error, error) {
	input := &validation.Validatable{
		JSON:     data,
		Rules:    rules,
		Handlers: handlers,
	}

	return validation.Validate(input) //nolint:wrapcheck
}

func collectRules(output any) map[string][]string {
	tagCollector := meta.NewTagsCollector("validate")

	return tagCollector.Extract(output)
}
