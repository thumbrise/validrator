// Package dot provides dot-notation functionality
package dot

import "strconv"

// Map converts nested map to flat dot notation projection of map. You want use this when input is result of json unmarshalling.
func Map(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	mapRecursive(input, result, "")

	return result
}

func mapRecursive(input any, output map[string]interface{}, outputKey string) {
	nextNodes := make(map[string]interface{})
	nextPrefix := ""

	if outputKey != "" {
		output[outputKey] = input
		nextPrefix = outputKey + "."
	}

	switch castedInput := input.(type) {
	case map[string]interface{}:
		for key, value := range castedInput {
			nextNodes[key] = value
		}
	case []interface{}:
		for key, value := range castedInput {
			nextNodes[strconv.Itoa(key)] = value
		}
	}

	for key, value := range nextNodes {
		nextOutputKey := nextPrefix + key

		mapRecursive(value, output, nextOutputKey)
	}
}
