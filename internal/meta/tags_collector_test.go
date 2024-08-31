package meta_test

import (
	"reflect"
	"testing"

	"github.com/thumbrise/validrator/internal/meta"
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
		SimpleField            string       `validate:"simple_field"`
		SimpleFieldWithOptions string       `validate:"simple_field_with_options,omitempty"`
		NestedStruct           innerStruct  `validate:"nested_struct"`
		NestedPointerStruct    *innerStruct `validate:"nested_pointer_struct"`
		DeepStruct             nestedStruct `validate:"deep_struct"`
		AnonNestedStruct       struct {
			JustField               int `validate:"just_field"`
			AnotherAnonNestedStruct struct {
				JustField int `validate:"just_field"`
			} `validate:"another_anon_nested_struct"`
		} `validate:"anon_nested_struct"`
		FieldWithMultipleTags bool `validate:"field_with_multiple_tags_json" xml:"field_with_multiple_tags_xml"`
		Slice                 []struct {
			JustBool        int    `validate:"just_bool"`
			NestedBoolSlice []bool `validate:"nested_bool_slice"`
		} `validate:"slice"`
		NestedSelfReference struct {
			JustField string      `validate:"just_field"`
			SelfRef   *testStruct `validate:"self_ref"` // this field must be ignored
		} `validate:"nested_self_reference"`

		// next must be ignored completely...
		SelfReference                *testStruct `validate:"self_reference"`
		FieldWithoutTags             float64
		FieldWithoutNeededTag        float64 `xml:"field_without_needed_tag"`
		FieldWithEmptyTagValue       int     `validate:""`
		FieldWithCommaInsteadOfValue int     `validate:""`
		FieldWithWhiteSpaceValue     string  `validate:"    "`
		FieldWithPrivateValue        int     `validate:"-"`
		FieldWithJustOptions         int     `validate:",omitempty"`
		FieldWithInvalidTag          int     `validate:field_with_invalid_tag` //nolint:govet
	}

	expected := map[string][]string{
		"SimpleField":            {"simple_field"},
		"SimpleFieldWithOptions": {"simple_field_with_options"},

		"NestedStruct":              {"nested_struct"},
		"NestedStruct.NestedFieldA": {"nested_field_a"},
		"NestedStruct.NestedFieldB": {"nested_field_b"},

		"NestedPointerStruct":              {"nested_pointer_struct"},
		"NestedPointerStruct.NestedFieldA": {"nested_field_a"},
		"NestedPointerStruct.NestedFieldB": {"nested_field_b"},

		"DeepStruct":                            {"deep_struct"},
		"DeepStruct.JustField":                  {"just_field"},
		"DeepStruct.AnotherStruct":              {"another_struct"},
		"DeepStruct.AnotherStruct.NestedFieldA": {"nested_field_a"},
		"DeepStruct.AnotherStruct.NestedFieldB": {"nested_field_b"},

		"AnonNestedStruct":                                   {"anon_nested_struct"},
		"AnonNestedStruct.JustField":                         {"just_field"},
		"AnonNestedStruct.AnotherAnonNestedStruct":           {"another_anon_nested_struct"},
		"AnonNestedStruct.AnotherAnonNestedStruct.JustField": {"just_field"},
		"FieldWithMultipleTags":                              {"field_with_multiple_tags_json"},

		"Slice":                   {"slice"},
		"Slice.*.JustBool":        {"just_bool"},
		"Slice.*.NestedBoolSlice": {"nested_bool_slice"},

		"NestedSelfReference":           {"nested_self_reference"},
		"NestedSelfReference.JustField": {"just_field"},
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
		"NestedStruct":              {"nested_struct"},
		"NestedStruct.NestedFieldA": {"nested_field_a"},
		"NestedStruct.NestedFieldB": {"nested_field_b"},

		"NestedPointerStruct":              {"nested_pointer_struct"},
		"NestedPointerStruct.NestedFieldA": {"nested_field_a"},
		"NestedPointerStruct.NestedFieldB": {"nested_field_b"},

		"DeepStruct":                            {"deep_struct"},
		"DeepStruct.JustField":                  {"just_field"},
		"DeepStruct.AnotherStruct":              {"another_struct"},
		"DeepStruct.AnotherStruct.NestedFieldA": {"nested_field_a"},
		"DeepStruct.AnotherStruct.NestedFieldB": {"nested_field_b"},
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
		"NestedSelfReference":           {"nested_self_reference"},
		"NestedSelfReference.JustField": {"just_field"},
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
		Slice []struct {
			JustBool        int    `validate:"just_bool"`
			NestedBoolSlice []bool `validate:"nested_bool_slice"`
		} `validate:"slice"`
	}

	expected := map[string][]string{
		"Slice":                   {"slice"},
		"Slice.*.JustBool":        {"just_bool"},
		"Slice.*.NestedBoolSlice": {"nested_bool_slice"},
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
