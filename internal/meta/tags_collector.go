// Package meta help collect meta information from go types
package meta

import (
	"errors"
	"reflect"
	"strings"

	strings2 "github.com/thumbrise/validrator/internal/strings"
)

const privateFieldVal = "-"

var errHierarchyFinished = errors.New("hierarchy finished")

// TagsCollector godoc.
type TagsCollector struct {
	tagKey string
}

// NewTagsCollector constructor.
func NewTagsCollector(tagKey string) *TagsCollector {
	return &TagsCollector{tagKey: tagKey}
}

// Extract returns flat map of founded tags with dot and star notation (field.nestedField: someTag, sliceField.*.someType: someAnotherTag).
func (t *TagsCollector) Extract(structure any) map[string][]string {
	return t.traverseHierarchy(structure)
}

func (t *TagsCollector) traverseHierarchy(structure any) map[string][]string {
	result := make(map[string][]string)

	toTraverse := make(map[string]reflect.StructField)
	typesChain := make(map[string]bool)
	_ = computeTraverseTree(structure, toTraverse, "", typesChain)

	for key, field := range toTraverse {
		rawTag := field.Tag.Get(t.tagKey)
		tagParts := strings.Split(rawTag, ",")

		if len(tagParts) == 0 {
			continue
		}

		for _, tagPart := range tagParts {
			if tagPart == privateFieldVal {
				continue
			}

			tagPart = strings.TrimSpace(tagPart)
			if tagPart == "" {
				continue
			}

			result[key] = append(result[key], tagPart)
		}
	}

	return result
}

func computeTraverseTree(unit interface{}, output map[string]reflect.StructField, hierarchyKeyPrefix string, typesChain map[string]bool) error { //nolint: cyclop // TODO: refactor
	var typ reflect.Type

	var fieldKey string

	var outputValue reflect.StructField

	var outputKey string

	switch val := unit.(type) {
	case reflect.StructField:
		outputValue = val
		fieldKey = strings2.ToCamel(outputValue.Name)
		typ = outputValue.Type
	case reflect.Type:
		typ = val

	default:
		typ = reflect.TypeOf(val)
	}

	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	stringType := generateStringType(typ)
	if stringType != "" && typesChain[stringType] {
		return errHierarchyFinished
	}

	typesChain[stringType] = true
	outputKey = hierarchyKeyPrefix + fieldKey

	switch typ.Kind() { //nolint:exhaustive
	case reflect.Struct:
		var newPrefix string
		if fieldKey != "" {
			newPrefix = outputKey + "."
		} else {
			newPrefix = hierarchyKeyPrefix
		}

		for i := range typ.NumField() {
			field := typ.Field(i)
			_ = computeTraverseTree(field, output, newPrefix, typesChain)
		}
	case reflect.Slice:
		newPrefix := ""
		if outputKey != "" {
			newPrefix = outputKey + ".*."
		}

		_ = computeTraverseTree(typ.Elem(), output, newPrefix, typesChain)

	default:
	}

	if outputKey != "" {
		output[outputKey] = outputValue
	}

	delete(typesChain, stringType)

	return errHierarchyFinished
}

func generateStringType(typ reflect.Type) string {
	if typ.PkgPath() == "" || typ.Name() == "" {
		return ""
	}

	return typ.PkgPath() + "." + typ.Name()
}
