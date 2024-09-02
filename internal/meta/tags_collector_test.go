package meta_test

import (
	"reflect"
	"testing"

	"github.com/thumbrise/validrator/internal/meta"
	"github.com/thumbrise/validrator/internal/testutil"
)

const tagKey = "validate"

func TestExtract(t *testing.T) {
	t.Parallel()

	type args struct {
		structure any
		tagKey    string
	}

	type innerStruct struct {
		NestedFieldA int `validate:"nested_field_a"`
		NestedFieldB int `validate:"nested_field_b"`
	}

	type nestedStruct struct {
		JustField     string      `validate:"just_field"`
		AnotherStruct innerStruct `validate:"another_struct"`
	}

	type testStruct struct {
		SimpleField            string `validate:"simple_field"`
		SimpleFieldWithOptions string `validate:"simple_field_with_options|omitempty"`

		NestedStruct innerStruct `validate:"nested_struct"`

		NestedPointerStruct *innerStruct `validate:"nested_pointer_struct"`

		DeepStruct nestedStruct `validate:"deep_struct"`

		AnonNestedStruct struct {
			JustField               int `validate:"just_field"`
			AnotherAnonNestedStruct struct {
				JustField int `validate:"just_field"`
			} `validate:"another_anon_nested_struct"`
		} `validate:"anon_nested_struct"`

		FieldWithMultipleTags bool `validate:"field_with_multiple_tags_json" xml:"field_with_multiple_tags_xml"`

		Slice []struct {
			JustBool        int    `validate:"just_bool"`
			NestedBoolSlice []bool `validate:"nested_bool_slice"`
		} `validate:"slice"`

		SliceWithIterativeRule []int `validate:"equals 1|[]equals 1"`

		NestedSelfReference struct {
			JustField string      `validate:"just_field"`
			SelfRef   *testStruct `validate:"self_ref"` // this field must be ignored
		} `validate:"nested_self_reference"`
		FieldWithSkippedOption int `validate:"|omitempty"`

		// next must be ignored completely...
		SelfReference                *testStruct `validate:"self_reference"`
		FieldWithoutTags             float64
		FieldWithoutNeededTag        float64 `xml:"field_without_needed_tag"`
		FieldWithEmptyTagValue       int     `validate:""`
		FieldWithCommaInsteadOfValue int     `validate:""`
		FieldWithWhiteSpaceValue     string  `validate:"    "`
		FieldWithPrivateValue        int     `validate:"-"`
		FieldWithInvalidTag          int     `validate:field_with_invalid_tag` //nolint:govet
	}

	expected := map[string][]string{
		"simpleField":            {"simple_field"},
		"simpleFieldWithOptions": {"simple_field_with_options", "omitempty"},

		"nestedStruct":              {"nested_struct"},
		"nestedStruct.nestedFieldA": {"nested_field_a"},
		"nestedStruct.nestedFieldB": {"nested_field_b"},

		"nestedPointerStruct":              {"nested_pointer_struct"},
		"nestedPointerStruct.nestedFieldA": {"nested_field_a"},
		"nestedPointerStruct.nestedFieldB": {"nested_field_b"},

		"deepStruct":                            {"deep_struct"},
		"deepStruct.justField":                  {"just_field"},
		"deepStruct.anotherStruct":              {"another_struct"},
		"deepStruct.anotherStruct.nestedFieldA": {"nested_field_a"},
		"deepStruct.anotherStruct.nestedFieldB": {"nested_field_b"},

		"anonNestedStruct":                                   {"anon_nested_struct"},
		"anonNestedStruct.justField":                         {"just_field"},
		"anonNestedStruct.anotherAnonNestedStruct":           {"another_anon_nested_struct"},
		"anonNestedStruct.anotherAnonNestedStruct.justField": {"just_field"},
		"fieldWithMultipleTags":                              {"field_with_multiple_tags_json"},

		"slice":                   {"slice"},
		"slice.*.justBool":        {"just_bool"},
		"slice.*.nestedBoolSlice": {"nested_bool_slice"},

		"sliceWithIterativeRule":   {"equals 1"},
		"sliceWithIterativeRule.*": {"equals 1"},

		"nestedSelfReference":           {"nested_self_reference"},
		"nestedSelfReference.justField": {"just_field"},

		"fieldWithSkippedOption": {"omitempty"},
	}

	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "should return correct tag set for value struct",
			args: args{
				structure: testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
		{
			name: "should return correct tag set for pointer struct",
			args: args{
				structure: &testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collector := meta.NewTagsCollector(tt.args.tagKey)
			got := collector.Extract(tt.args.structure)

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

func TestExtractInner(t *testing.T) {
	t.Parallel()

	type args struct {
		structure any
		tagKey    string
	}

	type innerStruct struct {
		NestedFieldA int `validate:"nested_field_a"`
		NestedFieldB int `validate:"nested_field_b"`
	}

	type nestedStruct struct {
		JustField     string      `validate:"just_field"`
		AnotherStruct innerStruct `validate:"another_struct"`
	}

	type testStruct struct {
		NestedStruct        innerStruct  `validate:"nested_struct"`
		NestedPointerStruct *innerStruct `validate:"nested_pointer_struct"`
		DeepStruct          nestedStruct `validate:"deep_struct"`
	}

	expected := map[string][]string{
		"nestedStruct":              {"nested_struct"},
		"nestedStruct.nestedFieldA": {"nested_field_a"},
		"nestedStruct.nestedFieldB": {"nested_field_b"},

		"nestedPointerStruct":              {"nested_pointer_struct"},
		"nestedPointerStruct.nestedFieldA": {"nested_field_a"},
		"nestedPointerStruct.nestedFieldB": {"nested_field_b"},

		"deepStruct":                            {"deep_struct"},
		"deepStruct.justField":                  {"just_field"},
		"deepStruct.anotherStruct":              {"another_struct"},
		"deepStruct.anotherStruct.nestedFieldA": {"nested_field_a"},
		"deepStruct.anotherStruct.nestedFieldB": {"nested_field_b"},
	}

	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "should return correct tag set for value struct",
			args: args{
				structure: testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
		{
			name: "should return correct tag set for pointer struct",
			args: args{
				structure: &testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collector := meta.NewTagsCollector(tt.args.tagKey)
			if got := collector.Extract(tt.args.structure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wrong result\nexpected:\n%+v\nactual:\n%+v", tt.want, got)
			}
		})
	}
}

func TestExtractSelfReference(t *testing.T) {
	t.Parallel()

	type args struct {
		structure any
		tagKey    string
	}

	type testStruct struct {
		SelfRef *testStruct `validate:"self_ref"` // this field must be ignored
	}

	expected := map[string][]string{}

	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "should return correct tag set for value struct",
			args: args{
				structure: testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
		{
			name: "should return correct tag set for pointer struct",
			args: args{
				structure: &testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collector := meta.NewTagsCollector(tt.args.tagKey)
			if got := collector.Extract(tt.args.structure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wrong result\nexpected:\n%+v\nactual:\n%+v", tt.want, got)
			}
		})
	}
}

func TestExtractNestedSelfReference(t *testing.T) {
	t.Parallel()

	type args struct {
		structure any
		tagKey    string
	}

	type testStruct struct {
		NestedSelfReference struct {
			JustField string      `validate:"just_field"`
			SelfRef   *testStruct `validate:"self_ref"` // this field must be ignored
		} `validate:"nested_self_reference"`
	}

	expected := map[string][]string{
		"nestedSelfReference":           {"nested_self_reference"},
		"nestedSelfReference.justField": {"just_field"},
	}

	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "should return correct tag set for value struct",
			args: args{
				structure: testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
		{
			name: "should return correct tag set for pointer struct",
			args: args{
				structure: &testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collector := meta.NewTagsCollector(tt.args.tagKey)
			if got := collector.Extract(tt.args.structure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wrong result\nexpected:\n%+v\nactual:\n%+v", tt.want, got)
			}
		})
	}
}

func TestExtractSlice(t *testing.T) {
	t.Parallel()

	type args struct {
		structure any
		tagKey    string
	}

	type testStruct struct {
		SomeField int `validate:"some_field"`
	}

	expected := map[string][]string{
		"someField": {"some_field"},
	}

	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "should return correct tag set for value struct",
			args: args{
				structure: testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
		{
			name: "should return correct tag set for pointer struct",
			args: args{
				structure: &testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collector := meta.NewTagsCollector(tt.args.tagKey)
			if got := collector.Extract(tt.args.structure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wrong result\nexpected:\n%+v\nactual:\n%+v", tt.want, got)
			}
		})
	}
}

func TestExtractInterface(t *testing.T) {
	t.Parallel()

	type args struct {
		structure any
		tagKey    string
	}

	type testStruct struct {
		Slice []struct {
			JustBool        int    `validate:"just_bool"`
			NestedBoolSlice []bool `validate:"nested_bool_slice"`
		} `validate:"slice"`
	}

	expected := map[string][]string{
		"slice":                   {"slice"},
		"slice.*.justBool":        {"just_bool"},
		"slice.*.nestedBoolSlice": {"nested_bool_slice"},
	}

	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "should return correct tag set for value struct",
			args: args{
				structure: testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
		{
			name: "should return correct tag set for pointer struct",
			args: args{
				structure: &testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collector := meta.NewTagsCollector(tt.args.tagKey)
			if got := collector.Extract(tt.args.structure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wrong result\nexpected:\n%+v\nactual:\n%+v", tt.want, got)
			}
		})
	}
}

func TestExtractIterativeSlice(t *testing.T) {
	t.Parallel()

	type args struct {
		structure any
		tagKey    string
	}

	type testStruct struct {
		SliceWithIterativeRule []int `validate:"equals 1|[]equals 1"`
	}

	expected := map[string][]string{
		"sliceWithIterativeRule":   {"equals 1"},
		"sliceWithIterativeRule.*": {"equals 1"},
	}

	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "should return correct tag set for value struct",
			args: args{
				structure: testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
		{
			name: "should return correct tag set for pointer struct",
			args: args{
				structure: &testStruct{},
				tagKey:    tagKey,
			},
			want: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collector := meta.NewTagsCollector(tt.args.tagKey)
			if got := collector.Extract(tt.args.structure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wrong result\nexpected:\n%+v\nactual:\n%+v", tt.want, got)
			}
		})
	}
}
