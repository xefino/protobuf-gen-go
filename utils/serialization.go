package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// None is an empty map that can be used in place of alternate values for the unmarshaller
var None = make(map[string]int32)

// Ignore is an empty map that can be used in place of alternate values for the marshaller
var Ignore = make(map[int32]string)

// NoValue returns the equivalent of no-value for one of the common financial enums. This function can
// be used when the user wants to deactivate a filter that relies on an enum message type
func NoValue[TEnum ~int32]() TEnum {
	return TEnum(-1)
}

// HasValue returns true if the value provided represents an actual enum value, or false if the value is
// equivalent to what would be produced by calling NoValue for that enum
func HasValue[TEnum ~int32](value TEnum) bool {
	return value != NoValue[TEnum]()
}

// MarshalString converts an enum value to a string based on the _value mapping and a possible alternative mapping
func MarshalString[TIn ~int32, TMap ~int32, TAlt ~int32](data TIn, mapping map[TMap]string, alternates map[TAlt]string, quote bool) string {
	var asStr string
	if alternate, ok := alternates[TAlt(data)]; ok {
		asStr = alternate
	} else if mapped, ok := mapping[TMap(data)]; ok {
		asStr = mapped
	} else {
		asStr = fmt.Sprintf("%d", data)
	}

	if quote {
		return "\"" + asStr + "\""
	} else {
		return asStr
	}
}

// UnmarshalValue converts JSON to a protobuf enum value
func UnmarshalValue[TMap ~int32, TAlt ~int32, TOut ~int32](raw []byte, mapping map[string]TMap,
	alternates map[string]TAlt, data *TOut) error {

	// Attempt to deserialize the value to a string to remove any escapes or
	// quotes that aren't needed; if this fails then return an error. If the
	// string isn't already quoted then we probably don't have any work to do
	// here so just trim the whitespace off and set it directly
	var asStr string
	if runes := []rune(string(raw)); len(runes) >= 2 && runes[0] == '"' && runes[len(runes)-1] == '"' {
		if err := json.Unmarshal(raw, &asStr); err != nil {
			return err
		}
	} else {
		asStr = string(raw)
	}

	// Convert the value from a string to its enum equivalent and return the result
	return UnmarshalString(asStr, mapping, alternates, data)
}

// ScanValue is intended to be used by functions that want to implement the Scanner
// interface for converting SQL values to enums
func ScanValue[TMap ~int32, TAlt ~int32, TOut ~int32](value interface{}, mapping map[string]TMap,
	alternates map[string]TAlt, data *TOut) error {

	// The type of the interface will define how we translate the value
	switch casted := value.(type) {

	// We have a string so we'll attempt to convert it to its code using either the primary
	// mapping or the alternate mapping as available; if the value doesn't map to either of
	// these cases then throw an error
	case string:
		if inner, ok := mapping[casted]; ok {
			*data = TOut(inner)
		} else if inner, ok := alternates[casted]; ok {
			*data = TOut(inner)
		} else {
			return fmt.Errorf("value of %v cannot be mapped to a %T", casted, *data)
		}

	// We have an integer value so cast it to a 32-bit integer
	case int:
		*data = TOut(int32(casted))
	case int16:
		*data = TOut(int32(casted))
	case int32:
		*data = TOut(casted)
	case int64:
		*data = TOut(int32(casted))

	// We aren't sure what type this is so throw an error
	default:
		return fmt.Errorf("value of %q had an invalid type of %T", casted, casted)
	}

	return nil
}

// UnmarshalString unmarshals a value from a string into an enum value
func UnmarshalString[TMap ~int32, TAlt ~int32, TOut ~int32](value string, mapping map[string]TMap,
	alternates map[string]TAlt, data *TOut) error {

	// First, trim all whitespace from the value
	asStr := strings.TrimSpace(value)

	// Next, attempt to lookup the value in the mapping; if it can be found
	// then return it; otherwise we'll continue on
	if value, ok := mapping[asStr]; ok {
		*data = TOut(value)
		return nil
	} else if value, ok := alternates[asStr]; ok {
		*data = TOut(value)
		return nil
	}

	// Finally, the value didn't map to one of the mapping options so attempt
	// to cover it to an integer directly; if this fails then return an error
	if value, err := strconv.ParseInt(asStr, 10, 32); err == nil {
		*data = TOut(int32(value))
	} else {
		return fmt.Errorf("value of %q cannot be mapped to a %T", asStr, *data)
	}

	return nil
}
