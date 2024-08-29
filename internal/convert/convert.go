// Package convert contains methods for conversion some data types to golang struct
package convert

import (
	"encoding/json"
	"errors"
	"io"
)

// JSONToStruct is converts io.Reader to golang struct.
func JSONToStruct(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(obj)

	return errors.Unwrap(err)
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
