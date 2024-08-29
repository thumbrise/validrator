package validrator_test

import (
	"testing"

	"github.com/thumbrise/validrator"
)

func TestValidrator_Validate(t *testing.T) {
	t.Parallel()

	t.Run("no error if fields are valid", func(t *testing.T) {
		t.Parallel()

		data := map[string]interface{}{
			"field": 1,
		}
		rules := map[string][]string{
			"field": {"equals 1"},
		}
		v := validrator.NewValidrator()
		v.AddRuleHandler("equals 1", func(v interface{}, _ []string) bool {
			return v == 1
		})

		err := v.Validate(data, rules)
		if err != nil {
			t.Errorf("Unexpected Validate() error\n%v", err)

			return
		}
	})

	t.Run("has error if invalid", func(t *testing.T) {
		t.Parallel()

		data := map[string]interface{}{
			"field": 2,
		}
		rules := map[string][]string{
			"field": {"equals 1"},
		}
		v := validrator.NewValidrator()
		v.AddRuleHandler("equals 1", func(v interface{}, _ []string) bool {
			return v == 1
		})

		err := v.Validate(data, rules)
		if err == nil {
			t.Error("Expected Validate() error but there is no\n")

			return
		}
	})

	t.Run("has rule error even if field missing", func(t *testing.T) {
		t.Parallel()

		data := map[string]interface{}{}
		rules := map[string][]string{
			"field": {"equals 1"},
		}
		v := validrator.NewValidrator()
		v.AddRuleHandler("equals 1", func(v interface{}, _ []string) bool {
			return v == 1
		})

		err := v.Validate(data, rules)
		if err == nil {
			t.Error("Expected Validate() error but there is no\n")

			return
		}
	})

	t.Run("no rule error if field missing with optional rule applied", func(t *testing.T) {
		t.Parallel()

		data := map[string]interface{}{}
		rules := map[string][]string{
			"field": {"equals 1", "optional"},
		}
		v := validrator.NewValidrator()
		v.AddRuleHandler("equals 1", func(v interface{}, _ []string) bool {
			return v == 1
		})

		err := v.Validate(data, rules)
		if err != nil {
			t.Errorf("Unexpected Validate() error\n%v", err)

			return
		}
	})
}
