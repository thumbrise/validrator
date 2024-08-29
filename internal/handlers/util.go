package handlers

import (
	"reflect"
	"strconv"
	"time"
)

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// asInt returns the parameter as a int64
// or panics if it can't convert.
func asInt(param string) int64 {
	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(err)

	return i
}

// asIntFromTimeDuration parses param as time.Duration and returns it as int64
// or panics on error.
func asIntFromTimeDuration(param string) int64 {
	dur, err := time.ParseDuration(param)
	if err != nil {
		// attempt parsing as an integer assuming nanosecond precision
		return asInt(param)
	}

	return int64(dur)
}

var timeDurationType = reflect.TypeOf(time.Duration(0))

// asIntFromType calls the proper function to parse param as int64,
// given a field's Type t.
func asIntFromType(t reflect.Type, param string) int64 {
	switch t {
	case timeDurationType:
		return asIntFromTimeDuration(param)
	default:
		return asInt(param)
	}
}

// asUint returns the parameter as a uint64
// or panics if it can't convert.
func asUint(param string) uint64 {
	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(err)

	return i
}

// asFloat64 returns the parameter as a float64
// or panics if it can't convert.
func asFloat64(param string) float64 {
	i, err := strconv.ParseFloat(param, 64)
	panicIf(err)

	return i
}

// asFloat64 returns the parameter as a float64
// or panics if it can't convert.
func asFloat32(param string) float64 {
	i, err := strconv.ParseFloat(param, 32)
	panicIf(err)

	return i
}

// asBool returns the parameter as a bool
// or panics if it can't convert.
func asBool(param string) bool { //nolint:unused
	i, err := strconv.ParseBool(param)
	panicIf(err)

	return i
}
