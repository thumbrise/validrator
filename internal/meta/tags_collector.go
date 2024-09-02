// Package meta help collect meta information from go types
package meta

import (
	"errors"
	"reflect"
	"strings"

	strings2 "github.com/thumbrise/validrator/internal/strings"
)

const (
	privateFieldVal = "-"
	iterativePrefix = "[]"
)

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

			// handle iterative tag
			if strings.HasPrefix(tagPart, iterativePrefix) {
				realTag := strings.TrimPrefix(tagPart, iterativePrefix)
				if realTag == "" {
					continue
				}

				// iterative tag applying to underlying values, so rewrite key with ".*" notation
				tagPart = realTag
				key += ".*"
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

	// Preventing self reference
	stringType := generateStringType(typ)
	if stringType != "" && typesChain[stringType] {
		return errHierarchyFinished
	}

	typesChain[stringType] = true
	defer delete(typesChain, stringType)

	outputKey = hierarchyKeyPrefix + fieldKey

	var nextUnits []interface{}

	var nextPrefix string

	switch typ.Kind() { //nolint:exhaustive
	case reflect.Struct:
		if fieldKey != "" {
			nextPrefix = outputKey + "."
		} else {
			nextPrefix = hierarchyKeyPrefix
		}

		for i := range typ.NumField() {
			field := typ.Field(i)
			nextUnits = append(nextUnits, field)
		}
	case reflect.Slice:
		if outputKey != "" {
			nextPrefix = outputKey + ".*."
		}

		nextUnits = append(nextUnits, typ.Elem())
	default:
	}

	if outputKey != "" {
		output[outputKey] = outputValue
	}

	for _, nextUnit := range nextUnits {
		_ = computeTraverseTree(nextUnit, output, nextPrefix, typesChain)
	}

	return errHierarchyFinished
}

func generateStringType(typ reflect.Type) string {
	if typ.PkgPath() == "" || typ.Name() == "" {
		return ""
	}

	return typ.PkgPath() + "." + typ.Name()
}
