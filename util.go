package validrator

import (
	"bytes"
	"encoding/json"
	"errors"
)

// JSONToStruct is converts io.Reader to golang struct.
func jsonToStruct(input []byte, obj any) error {
	decoder := json.NewDecoder(bytes.NewReader(input))
	err := decoder.Decode(obj)

	return errors.Unwrap(err)
}

// JSONToMap is converts io.Reader to golang map.
func jsonToMap(input []byte, output map[string]interface{}) error {
	err := json.Unmarshal(input, &output)
	if err != nil {
		return errors.Unwrap(err)
	}

	return nil
}
