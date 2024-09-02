//nolint:cyclop,maintidx
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
		Field1 int     `validate:"required|equals 1"`
		Field2 int8    `validate:"required|equals 1"`
		Field3 int16   `validate:"required|equals 1"`
		Field4 int32   `validate:"required|equals 1"`
		Field5 int64   `validate:"required|equals 1"`
		Field6 float32 `validate:"required|equals 1"`
		Field7 float64 `validate:"required|equals 1"`
	}

	type testStructNumbersWithEquals1AndNotRequired struct {
		Field1 int     `json:"field1"       validate:"equals 1"`
		Field2 int8    `validate:"equals 1"`
		Field3 int16   `validate:"equals 1"`
		Field4 int32   `validate:"equals 1"`
		Field5 int64   `validate:"equals 1"`
		Field6 float32 `validate:"equals 1"`
		Field7 float64 `validate:"equals 1"`
	}

	type testStructMixed struct {
		SomeInt    int         `validate:"required|equals 1"`
		SomeFloat  float64     `validate:"required|equals 1"`
		SomeBool   bool        `validate:"required|equals 1"`
		SomeNull   interface{} `validate:"required|equals 1"`
		SomeObject struct {
			A string `validate:"required|equals 1"`
			B string `validate:"required|equals 1"`
		} `validate:"required|equals 1"`
		SomeArray []int `validate:"required|equals 1"`
	}

	type testStructWithObject struct {
		SomeObject struct {
			A1 struct {
				A2 int `validate:"equals 1"`
			} `validate:"equals 1"`
			B1 struct {
				B2 int `validate:"equals 1"`
			} `validate:"equals 1"`
		} `validate:"equals 1"`
	}

	type testStructWithRequiredTags struct {
		SomeNull       int
		SomeZero       int
		SomeNotExists1 int `validate:"required|equals 1"`
		SomeNotExists2 int `validate:"required|equals 1"`
		SomeNotExists3 int `validate:"equals 1"`
	}

	type testStructWithArray struct {
		SomeArrayNull      []int `validate:"nullable"`
		SomeArrayNotExists []int `validate:"optional"`
		SomeArrayEmpty     []int `validate:"equals 1"`
	}

	type testStructWithIterativeTag struct {
		SomeArray          []int `validate:"required|equals 1|[]equals 1"`
		SomeArrayEmpty     []int `validate:"required|equals 1|[]equals 1"`
		SomeArrayNull      []int `validate:"required|equals 1|[]equals 1"`
		SomeArrayNotExists []int `validate:"required|equals 1|[]equals 1"`
	}

	type testStructWithoutTags struct {
		SomeField1 string
		SomeField2 string
		SomeField3 string
		SomeField4 string
	}

	type testStructWithIterativeRequiredTag struct {
		SomeField []int `validate:"[]required"`
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
			name: "testStructNumbersWithEquals1AndNotRequired should valid",
			inputJSON: `{
				"field1": 1,
				"field2": 1,
				"field3": 1,
				"field4": 1
			}`,
			expectedOutput: &testStructNumbersWithEquals1AndNotRequired{
				Field1: 1,
				Field2: 1,
				Field3: 1,
				Field4: 1,
				Field5: 0,
				Field6: 0,
				Field7: 0,
			},
			actualOutput:   &testStructNumbersWithEquals1AndNotRequired{},
			expectedErrors: map[string][]string{},
			wantErr:        false,
		},
		{
			name: "testStructNumbersWithEquals1AndNotRequired should invalid only by equals 1 rule but without required rule",
			inputJSON: `{
				"field1": 1,
				"field2": 1,
				"field3": 1,
				"field4": 2
			}`,
			expectedOutput: &testStructNumbersWithEquals1AndNotRequired{},
			actualOutput:   &testStructNumbersWithEquals1AndNotRequired{},
			expectedErrors: map[string][]string{
				"field4": {"equals 1"},
			},
			wantErr: false,
		},
		{
			name: "testStructNumbersWithEquals1AndNotRequired should error invalid json",
			// in this json trailing comma at the end. Its invalid in json
			inputJSON: `{
				"field1": 1,
				"field2": 1,
				"field3": 1,
				"field4": 1,
			}`,
			expectedOutput: &testStructNumbersWithEquals1AndNotRequired{},
			actualOutput:   &testStructNumbersWithEquals1AndNotRequired{},
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
				"someFloat":    {ruleEquals1.name},
				"someBool":     {ruleEquals1.name},
				"someNull":     {"required"},
				"someObject":   {ruleEquals1.name},
				"someObject.a": {ruleEquals1.name},
				"someObject.b": {ruleEquals1.name},
				"someArray":    {ruleEquals1.name},
			},
			wantErr: false,
		},
		{
			name: "testStructWithObject should invalid",
			inputJSON: `{
				"someObject": {
					"a1": {"a2": 2}, 
					"b1": {"b2": 2}
				}
			}`,

			expectedOutput: &testStructWithObject{},
			actualOutput:   &testStructWithObject{},
			expectedErrors: map[string][]string{
				"someObject":       {ruleEquals1.name},
				"someObject.a1":    {ruleEquals1.name},
				"someObject.b1":    {ruleEquals1.name},
				"someObject.a1.a2": {ruleEquals1.name},
				"someObject.b1.b2": {ruleEquals1.name},
			},
			wantErr: false,
		},
		{
			name: "testStructWithRequiredTags should some valid some invalid",
			inputJSON: `{
				"someNull": null,
				"someZero": 0
			}`,

			expectedOutput: &testStructWithRequiredTags{
				SomeNull:       0,
				SomeZero:       0,
				SomeNotExists1: 0,
				SomeNotExists2: 0,
				SomeNotExists3: 0,
			},
			actualOutput: &testStructWithRequiredTags{},
			expectedErrors: map[string][]string{
				"someNotExists1": {"required"},
				"someNotExists2": {"required"},
			},
			wantErr: false,
		},
		{
			name: "testStructWithArray should some valid some invalid",
			inputJSON: `{
				"someArrayNull": null,
				"someArrayEmpty": []
			}`,

			expectedOutput: &testStructWithArray{
				SomeArrayNull:      nil,
				SomeArrayNotExists: nil,
				SomeArrayEmpty:     nil,
			},
			actualOutput: &testStructWithArray{},
			expectedErrors: map[string][]string{
				"someArrayEmpty": {ruleEquals1.name},
			},
			wantErr: false,
		},
		{
			name: "testStructWithIterativeTag should some valid some invalid",
			inputJSON: `{
				"someArray": [1,1,1,1,0,1,1,1,0,1],
				"someArrayEmpty": [],
				"someArrayNull": null
			}`,

			expectedOutput: &testStructWithIterativeTag{},
			actualOutput:   &testStructWithIterativeTag{},
			expectedErrors: map[string][]string{
				"someArray":          {ruleEquals1.name},
				"someArray.4":        {ruleEquals1.name},
				"someArray.8":        {ruleEquals1.name},
				"someArrayEmpty":     {ruleEquals1.name},
				"someArrayNull":      {"required"},
				"someArrayNotExists": {"required"},
			},
			wantErr: false,
		},
		{
			name: "testStructWithoutTags should invalid required even without tags",
			inputJSON: `{
			}`,

			expectedOutput: &testStructWithoutTags{},
			actualOutput:   &testStructWithoutTags{},
			expectedErrors: map[string][]string{},
			wantErr:        false,
		},
		{
			name: "testStructWithIterativeRequiredTag should not return error",
			inputJSON: `{
				"someField": [1, null]
			}`,

			expectedOutput: &testStructWithIterativeRequiredTag{SomeField: []int{1, 0}},
			actualOutput:   &testStructWithIterativeRequiredTag{},
			expectedErrors: map[string][]string{},
			wantErr:        false,
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
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)

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
						"validation errors not match\nexpected:\n%#v\nactual:\n%#v\ndiff:\n%s\n",
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
					t.Errorf("Validate() got = %v, want %v", actualOutput, tt.expectedOutput)
				}

				return
			}
		})
	}
}
