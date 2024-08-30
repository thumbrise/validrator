// Package convert contains methods for conversion some data types
package convert

import (
	"bytes"
	"encoding/json"
	"errors"
)

// JSONToStruct is converts io.Reader to golang struct.
func JSONToStruct(input []byte, obj any) error {
	decoder := json.NewDecoder(bytes.NewReader(input))
	err := decoder.Decode(obj)

	return errors.Unwrap(err)
}

// JSONToMap is converts io.Reader to golang map.
func JSONToMap(input []byte, output map[string]interface{}) error {
	err := json.Unmarshal(input, &output)
	if err != nil {
		return errors.Unwrap(err)
	}

	return nil
}

// MapToStruct is converts golang map to golang struct.
func MapToStruct(input map[string]interface{}, obj any) error {
	marshaled, err := json.Marshal(input)
	if err != nil {
		return errors.Unwrap(err)
	}

	err = json.Unmarshal(marshaled, obj)
	if err != nil {
		return errors.Unwrap(err)
	}

	return nil
}
