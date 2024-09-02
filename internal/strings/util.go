// Package strings is utility
//
//nolint:varnamelen,gocritic,cyclop
package strings

import "strings"

// ToCamel Converts a string to CamelCase.
func ToCamel(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))

	capNext := false
	prevIsCap := false

	for i, v := range []byte(s) {
		if v == '.' || v == '*' {
			capNext = false

			n.WriteByte(v)

			continue
		}

		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'

		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		} else if prevIsCap && vIsCap {
			v += 'a'
			v -= 'A'
		}

		prevIsCap = vIsCap

		if vIsCap || vIsLow {
			n.WriteByte(v)

			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)

			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-'
		}
	}

	return n.String()
}
