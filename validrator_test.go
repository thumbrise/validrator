//nolint:cyclop
package validrator_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/thumbrise/validrator"
	"github.com/thumbrise/validrator/internal/validation"
)

func TestValidrator_Validate(t *testing.T) {
	t.Parallel()

	ruleEquals1 := struct {
		name    string
		handler validation.RuleHandlerFunc
	}{
		name: "equals 1",
		handler: func(v reflect.Value, _ []string) bool {
			const validValue = 1
			r := (v.CanFloat() && v.Float() == validValue) ||
				(v.CanInt() && v.Int() == validValue) ||
				(v.CanComplex() && v.Complex() == validValue)

			return r
		},
	}
	handlers := map[string]validation.RuleHandlerFunc{
		ruleEquals1.name: ruleEquals1.handler,
	}

	type testStructNumbersWithEquals1 struct {
		Field1 int     `validate:"equals 1"`
		Field2 int8    `validate:"equals 1"`
		Field3 int16   `validate:"equals 1"`
		Field4 int32   `validate:"equals 1"`
		Field5 int64   `validate:"equals 1"`
		Field6 float32 `validate:"equals 1"`
		Field7 float64 `validate:"equals 1"`
	}

	type testStructNumbersWithEquals1AndOptional struct {
		Field1 int     `json:"field1"                 validate:"equals 1,optional"`
		Field2 int8    `validate:"equals 1,optional"`
		Field3 int16   `validate:"equals 1, optional"`
		Field4 int32   `validate:"equals 1,optional"`
		Field5 int64   `validate:"equals 1, optional"`
		Field6 float32 `validate:"equals 1,optional"`
		Field7 float64 `validate:"equals 1, optional"`
	}

	type testStructMixed struct {
		SomeInt    int         `validate:"equals 1"`
		SomeFloat  float64     `validate:"equals 1"`
		SomeBool   bool        `validate:"equals 1"`
		SomeNull   interface{} `validate:"equals 1"`
		SomeObject struct {
			A string `validate:"equals 1"`
			B string `validate:"equals 1"`
		} `validate:"equals 1"`
		SomeArray []int `validate:"equals 1"`
	}

	tests := []struct {
		name           string
		inputJSON      string
		expectedOutput any
		actualOutput   any
		expectedErrors map[string][]string
		wantErr        bool
	}{
		{
			name: "testStructNumbersWithEquals1 should valid",
			inputJSON: `{
				"field1": 1,
				"field2": 1,
				"field3": 1,
				"field4": 1,
				"field5": 1,
				"field6": 1,
				"field7": 1
			}`,
			expectedOutput: &testStructNumbersWithEquals1{
				Field1: 1,
				Field2: 1,
				Field3: 1,
				Field4: 1,
				Field5: 1,
				Field6: 1,
				Field7: 1,
			},
			actualOutput:   &testStructNumbersWithEquals1{},
			expectedErrors: map[string][]string{},
			wantErr:        false,
		},
		{
			name: "testStructNumbersWithEquals1 should invalid",
			inputJSON: `{
				"field1": 1,
				"field2": 1,
				"field3": 1,
				"field4": 2
			}`,
			expectedOutput: &testStructNumbersWithEquals1{},
			actualOutput:   &testStructNumbersWithEquals1{},
			expectedErrors: map[string][]string{
				"field4": {"equals 1"},
				"field5": {"required"},
				"field6": {"required"},
				"field7": {"required"},
			},
			wantErr: false,
		},
		{
			name: "testStructNumbersWithEquals1AndOptional should valid",
			inputJSON: `{
				"field1": 1,
				"field2": 1,
				"field3": 1,
				"field4": 1
			}`,
			expectedOutput: &testStructNumbersWithEquals1AndOptional{
				Field1: 1,
				Field2: 1,
				Field3: 1,
				Field4: 1,
				Field5: 0,
				Field6: 0,
				Field7: 0,
			},
			actualOutput:   &testStructNumbersWithEquals1AndOptional{},
			expectedErrors: map[string][]string{},
			wantErr:        false,
		},
		{
			name: "testStructNumbersWithEquals1AndOptional should invalid only by equals 1 rule but without required rule",
			inputJSON: `{
				"field1": 1,
				"field2": 1,
				"field3": 1,
				"field4": 2
			}`,
			expectedOutput: &testStructNumbersWithEquals1AndOptional{},
			actualOutput:   &testStructNumbersWithEquals1AndOptional{},
			expectedErrors: map[string][]string{
				"field4": {"equals 1"},
			},
			wantErr: false,
		},
		{
			name: "testStructNumbersWithEquals1AndOptional should error invalid json",
			// in this json trailing comma at the end. Its invalid in json
			inputJSON: `{
				"field1": 1,
				"field2": 1,
				"field3": 1,
				"field4": 1,
			}`,
			expectedOutput: &testStructNumbersWithEquals1AndOptional{},
			actualOutput:   &testStructNumbersWithEquals1AndOptional{},
			expectedErrors: nil,
			wantErr:        true,
		},
		{
			name: "testStructMixed should invalid",
			inputJSON: `{
				"someInt": 1,
				"someFloat": 1234.1234,
				"someBool": true,
				"someNull": null,
				"someObject": {"a": "a", "b": "b"},
				"someArray": [1,2,3,4,5,6,7,8]
			}`,

			expectedOutput: &testStructMixed{},
			actualOutput:   &testStructMixed{},
			expectedErrors: map[string][]string{
				"someFloat":  {ruleEquals1.name},
				"someBool":   {ruleEquals1.name},
				"someNull":   {ruleEquals1.name},
				"someObject": {ruleEquals1.name},
				"someArray":  {ruleEquals1.name},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			validator := validrator.NewValidrator()

			validator.AddRuleHandlers(handlers)

			actualOutput := tt.actualOutput

			actualValidationErrors, err := validator.Validate([]byte(tt.inputJSON), actualOutput)

			// logic error check
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %validator, wantErr %validator", err, tt.wantErr)

				return
			}

			// validation errors check
			if len(tt.expectedErrors) > 0 {
				if actualValidationErrors == nil {
					t.Errorf("Validation errors are missing, but wanted")

					return
				}

				expectedValidationErrors := tt.expectedErrors

				expectedJSON, _ := json.MarshalIndent(expectedValidationErrors, "", "    ")
				actualJSON, _ := json.MarshalIndent(actualValidationErrors.ToMap(), "", "    ")

				diff := cmp.Diff(expectedJSON, actualJSON)
				if diff != "" {
					t.Errorf(
						"validation errors not match\nexpected:\n%#validator\nactual:\n%#validator\ndiff:\n%s\n",
						expectedValidationErrors,
						actualValidationErrors.ToMap(),
						diff,
					)
				}

				return
			}

			// output check
			if err == nil {
				if !reflect.DeepEqual(actualOutput, tt.expectedOutput) {
					t.Errorf("Validate() got = %validator, want %validator", actualOutput, tt.expectedOutput)
				}

				return
			}
		})
	}
}
