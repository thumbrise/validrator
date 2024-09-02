package dot_test

import (
	"testing"

	"github.com/thumbrise/validrator/internal/dot"
	"github.com/thumbrise/validrator/internal/testutil"
)

func TestMap(t *testing.T) {
	t.Parallel()

	type args struct {
		input map[string]interface{}
	}

	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "base",
			args: args{
				input: map[string]interface{}{
					"someInt":   1,
					"someFloat": 1234.1234,
					"someBool":  true,
					"someObject": map[string]interface{}{
						"a": "a",
						"b": "b",
					},
					"someArray": []interface{}{
						1,
						2,
						3,
						4,
						5,
						6,
						7,
						8,
					},
				},
			},
			want: map[string]interface{}{
				"someInt":   1,
				"someFloat": 1234.1234,
				"someBool":  true,
				"someObject": map[string]interface{}{
					"a": "a",
					"b": "b",
				},
				"someObject.a": "a",
				"someObject.b": "b",
				"someArray": []interface{}{
					1,
					2,
					3,
					4,
					5,
					6,
					7,
					8,
				},
				"someArray.0": 1,
				"someArray.1": 2,
				"someArray.2": 3,
				"someArray.3": 4,
				"someArray.4": 5,
				"someArray.5": 6,
				"someArray.6": 7,
				"someArray.7": 8,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := dot.Map(tt.args.input)

			diff := testutil.DiffAsJSON(tt.want, got)
			if diff != "" {
				t.Errorf(
					"validation errors not match\nexpected:\n%#v\nactual:\n%#v\ndiff:\n%s\n",
					tt.want,
					got,
					diff,
				)
			}
		})
	}
}
