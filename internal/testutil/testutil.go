// Package testutil godoc
//
//nolint:errchkjson
package testutil

import (
	"encoding/json"

	"github.com/google/go-cmp/cmp"
)

// DiffAsJSON converts to indented json and return diff between a and b.
func DiffAsJSON(expected any, actual any) string {
	expectedJSON, _ := json.MarshalIndent(expected, "", "    ")
	actualJSON, _ := json.MarshalIndent(actual, "", "    ")

	return cmp.Diff(expectedJSON, actualJSON)
}
