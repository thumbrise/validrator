package convert_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/thumbrise/validrator/internal/convert"
)

type mock struct {
	SomeText   string
	SomeInt    int
	SomeFloat  float64
	SomeBool   bool
	SomeNull   interface{}
	SomeObject struct {
		A string
		B string
	}
	SomeArray []int
}

func TestJsonToStruct(t *testing.T) {
	t.Parallel()

	type args struct {
		jsonString  string
		objExpected any
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				jsonString: `{
"someText": "text",	
"someInt": 1234,	
"someFloat": 1234.1234,	
"someBool": true,	
"someNull": null,
"someObject": {"a": "a", "b": "b"},
"someArray": [1,2,3,4,5]
}`,
				objExpected: mock{
					SomeText:  "text",
					SomeInt:   1234,
					SomeFloat: 1234.1234,
					SomeBool:  true,
					SomeNull:  nil,
					SomeObject: struct {
						A string
						B string
					}{
						A: "a",
						B: "b",
					},
					SomeArray: []int{1, 2, 3, 4, 5},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := strings.NewReader(tt.args.jsonString)
			objActual := mock{}

			if err := convert.JSONToStruct(r, &objActual); (err != nil) != tt.wantErr {
				t.Errorf("JSONToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(objActual, tt.args.objExpected) {
				t.Errorf("JSONToStruct() expected %+v actual %v", tt.args.objExpected, objActual)
			}
		})
	}
}

func TestMapToStruct(t *testing.T) {
	t.Parallel()

	type args struct {
		input       map[string]interface{}
		objExpected any
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				input: map[string]any{
					"someText":   "text",
					"someInt":    1234,
					"someFloat":  1234.1234,
					"someBool":   true,
					"someNull":   nil,
					"someObject": map[string]any{"a": "a", "b": "b"},
					"someArray":  []int{1, 2, 3, 4, 5},
				},
				objExpected: mock{
					SomeText:  "text",
					SomeInt:   1234,
					SomeFloat: 1234.1234,
					SomeBool:  true,
					SomeNull:  nil,
					SomeObject: struct {
						A string
						B string
					}{
						A: "a",
						B: "b",
					},
					SomeArray: []int{1, 2, 3, 4, 5},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			objActual := mock{}

			if err := convert.MapToStruct(tt.args.input, &objActual); (err != nil) != tt.wantErr {
				t.Errorf("MapToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(objActual, tt.args.objExpected) {
				t.Errorf("MapToStruct() expected %+v actual %v", tt.args.objExpected, objActual)
			}
		})
	}
}
