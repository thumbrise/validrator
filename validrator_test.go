package validrator_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/thumbrise/validrator"
)

var mockRule = struct {
	name    string
	handler validrator.RuleHandlerFunc
}{
	name: "equals 1",
	handler: func(v reflect.Value, _ []string) bool {
		const cmp = 1
		r := (v.CanFloat() && v.Float() == cmp) ||
			(v.CanInt() && v.Int() == cmp) ||
			(v.CanComplex() && v.Complex() == cmp)

		return r
	},
}

func TestValidrator_ValidateMap(t *testing.T) {
	t.Parallel()

	t.Run("no error if fields are valid", func(t *testing.T) {
		t.Parallel()

		data := map[string]interface{}{
			"field": 1,
		}
		rules := map[string][]string{
			"field": {mockRule.name},
		}
		v := validrator.NewValidrator()
		v.AddRuleHandler(mockRule.name, mockRule.handler)

		err := v.ValidateMap(data, rules)
		if err != nil {
			t.Errorf("Unexpected ValidateMap() error\n%v", err)

			return
		}
	})

	t.Run("has error if invalid", func(t *testing.T) {
		t.Parallel()

		data := map[string]interface{}{
			"field": 2,
		}
		rules := map[string][]string{
			"field": {mockRule.name},
		}
		v := validrator.NewValidrator()
		v.AddRuleHandler(mockRule.name, mockRule.handler)

		err := v.ValidateMap(data, rules)
		if err == nil {
			t.Error("Expected ValidateMap() error but there is no\n")

			return
		}
	})

	t.Run("has rule error even if field missing", func(t *testing.T) {
		t.Parallel()

		data := map[string]interface{}{}
		rules := map[string][]string{
			"field": {mockRule.name},
		}
		v := validrator.NewValidrator()
		v.AddRuleHandler(mockRule.name, mockRule.handler)

		err := v.ValidateMap(data, rules)
		if err == nil {
			t.Error("Expected ValidateMap() error but there is no\n")

			return
		}
	})

	t.Run("no rule error if field missing with optional rule applied", func(t *testing.T) {
		t.Parallel()

		data := map[string]interface{}{}
		rules := map[string][]string{
			"field": {mockRule.name, "optional"},
		}
		v := validrator.NewValidrator()
		v.AddRuleHandler(mockRule.name, mockRule.handler)

		err := v.ValidateMap(data, rules)
		if err != nil {
			t.Errorf("Unexpected ValidateMap() error\n%v", err)

			return
		}
	})
}

func TestValidrator_ValidateJsonToStruct(t *testing.T) {
	t.Parallel()

	t.Run("", func(t *testing.T) {
		t.Parallel()

		type outputStruct struct {
			Field int
		}

		inputJSON := `{
"field": 1
}`

		data := strings.NewReader(inputJSON)

		rules := map[string][]string{
			"field": {mockRule.name},
		}

		vldtr := validrator.NewValidrator()
		vldtr.AddRuleHandler(mockRule.name, mockRule.handler)

		output := outputStruct{}

		err := vldtr.ValidateJSONReaderToStruct(data, rules, &output)
		if err != nil {
			t.Errorf("Unexpected ValidateMap() error\n%vldtr", err)

			return
		}

		if output.Field != 1 {
			t.Errorf("Unexpected field value '%vldtr' Expects 1", output.Field)
		}
	})
}
