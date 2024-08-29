package handlers_test

import (
	"testing"

	"github.com/thumbrise/validrator"
	"github.com/thumbrise/validrator/internal/handlers"
)

func TestBool(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		data := map[string]interface{}{
			"field1": true,
			"field2": false,
			"field3": 1,
			"field4": 0,
			"field5": "true",
			"field6": "false",
			"field7": "0",
			"field8": "1",
		}
		rules := map[string][]string{
			"field1": {"bool"},
			"field2": {"bool"},
			"field3": {"bool"},
			"field4": {"bool"},
			"field5": {"bool"},
			"field6": {"bool"},
			"field7": {"bool"},
			"field8": {"bool"},
		}

		v := validrator.NewValidrator()
		v.AddRuleHandler("bool", handlers.Bool)

		err := v.Validate(data, rules)
		if err != nil {
			t.Errorf("Unexpected Validate() error\n%v", err)

			return
		}
	})
	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		data := map[string]interface{}{
			"field1": nil,
			"field2": 10,
			"field3": 0o1,
			"field4": "trues",
			"field5": "falsee",
			"field6": "01",
			"field7": "10",
		}
		rules := map[string][]string{
			"field1": {"bool"},
			"field2": {"bool"},
			"field3": {"bool"},
			"field4": {"bool"},
			"field5": {"bool"},
			"field6": {"bool"},
			"field7": {"bool"},
		}

		v := validrator.NewValidrator()
		v.AddRuleHandler("bool", handlers.Bool)

		err := v.Validate(data, rules)
		if err == nil {
			t.Error("Expected Validate() error")

			return
		}
	})
}
