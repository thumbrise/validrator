package convert_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/thumbrise/validrator/internal/convert"
)

type mock struct {
	SomeText   string      `json:"someText,omitempty"`
	SomeInt    int         `json:"someInt,omitempty"`
	SomeFloat  float64     `json:"someFloat,omitempty"`
	SomeBool   bool        `json:"someBool,omitempty"`
	SomeNull   interface{} `json:"someNull,omitempty"`
	SomeObject struct {
		A string `json:"a,omitempty"`
		B string `json:"b,omitempty"`
	} `json:"someObject"`
	SomeArray []int `json:"someArray,omitempty"`
}

var mockJSON = `{
"someText": "text",	
"someInt": 1234,	
"someFloat": 1234.1234,	
"someBool": true,	
"someNull": null,
"someObject": {"a": "a", "b": "b"},
"someArray": [1,2,3,4,5,6,7,8]
}`

var mockMap = map[string]any{
	"someText":   "text",
	"someInt":    1234,
	"someFloat":  1234.1234,
	"someBool":   true,
	"someNull":   nil,
	"someObject": map[string]any{"a": "a", "b": "b"},
	"someArray":  []any{1, 2, 3, 4, 5, 6, 7, 8},
}

var mockStruct = mock{
	SomeText:  "text",
	SomeInt:   1234,
	SomeFloat: 1234.1234,
	SomeBool:  true,
	SomeNull:  nil,
	SomeObject: struct {
		A string `json:"a,omitempty"`
		B string `json:"b,omitempty"`
	}(struct {
		A string
		B string
	}{
		A: "a",
		B: "b",
	}),
	SomeArray: []int{1, 2, 3, 4, 5, 6, 7, 8},
}

func TestJSONToStruct(t *testing.T) {
	t.Parallel()

	type args struct {
		jsonString string
		expected   any
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				jsonString: mockJSON,
				expected:   mockStruct,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := mock{}

			if err := convert.JSONToStruct([]byte(tt.args.jsonString), &actual); (err != nil) != tt.wantErr {
				t.Errorf("JSONToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(actual, tt.args.expected) {
				t.Errorf("JSONToStruct() expected %+v actual %v", tt.args.expected, actual)
			}
		})
	}
}

func TestJSONToMap(t *testing.T) {
	t.Parallel()

	type args struct {
		jsonString string
		expected   map[string]interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				jsonString: mockJSON,
				expected:   mockMap,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := make(map[string]interface{})
			if err := convert.JSONToMap([]byte(tt.args.jsonString), actual); (err != nil) != tt.wantErr {
				t.Errorf("JSONToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}

			actualJSONBytes, _ := json.Marshal(actual)
			expectedJSONBytes, _ := json.Marshal(tt.args.expected)
			actualJSON := string(actualJSONBytes)
			expectedJSON := string(expectedJSONBytes)

			if actualJSON != expectedJSON {
				t.Errorf("JSONToStruct() expected %s actual %s", expectedJSON, actualJSON)
			}
		})
	}
}

func TestMapToStruct(t *testing.T) {
	t.Parallel()

	type args struct {
		input    map[string]interface{}
		expected any
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				input:    mockMap,
				expected: mockStruct,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := mock{}

			if err := convert.MapToStruct(tt.args.input, &actual); (err != nil) != tt.wantErr {
				t.Errorf("MapToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(actual, tt.args.expected) {
				t.Errorf("MapToStruct() expected %+v actual %v", tt.args.expected, actual)
			}
		})
	}
}
