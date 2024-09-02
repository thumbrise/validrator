package validation

import (
	"errors"
	"fmt"
	"strings"
)

var errInvalidRule = errors.New("invalid rule")

// Special tags/rules.
const (
	// tagOptional allow not apply "required" error if field does not exist in input.
	tagOptional = "optional"
	// tagOptional allow not apply "required" error if field equals nil.
	tagNullable = "nullable"

	// tagRequired define rule which returns in validation error if optional/nullable not applied but field empty or even does not exist.
	tagRequired = "required"
)

// FieldValidationFail is fail entry of field.
type FieldValidationFail struct {
	Field string
	Rules []string
	Value interface{}
}

// Error are set of fail entries.
type Error struct {
	Failed map[string]FieldValidationFail
}

// ToMap godoc.
func (v *Error) ToMap() map[string][]string {
	result := make(map[string][]string, len(v.Failed))

	for key, fail := range v.Failed {
		result[key] = fail.Rules
	}

	return result
}

func (v *Error) Error() string {
	builder := strings.Builder{}

	for fieldName, fail := range v.Failed {
		rulesStr := strings.Join(fail.Rules, ", ")
		builder.WriteString(fmt.Sprintf("field=%s rules=%s value=%+v\n", fieldName, rulesStr, fail.Value))
	}

	return builder.String()
}
