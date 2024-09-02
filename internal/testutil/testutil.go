// Package testutil godoc
//
//nolint:errchkjson
package testutil

import (
	"encoding/json"

	"github.com/google/go-cmp/cmp"
)

// DiffAsJSON converts to indented json and return diff between a and b.
func DiffAsJSON(a any, b any) string {
	expectedJSON, _ := json.MarshalIndent(a, "", "    ")
	actualJSON, _ := json.MarshalIndent(b, "", "    ")

	return cmp.Diff(expectedJSON, actualJSON)
}
