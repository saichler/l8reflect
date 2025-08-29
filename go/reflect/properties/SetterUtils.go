package properties

import (
	"reflect"
	"strconv"
)

func ConvertValue(target, source reflect.Value) reflect.Value {
	targetKind := target.Kind()
	sourceKind := source.Kind()

	// Handle numeric conversions
	if isNumeric(targetKind) && isNumeric(sourceKind) {
		return convertNumeric(target, source)
	}

	// Handle string conversion
	if targetKind == reflect.String {
		return reflect.ValueOf(convertToString(source))
	}

	// Return source unchanged if no conversion needed
	return source
}

func isNumeric(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func convertNumeric(target, source reflect.Value) reflect.Value {
	targetKind := target.Kind()
	sourceKind := source.Kind()

	// Get the numeric value as float64 for conversion
	var numValue float64
	switch sourceKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		numValue = float64(source.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		numValue = float64(source.Uint())
	case reflect.Float32, reflect.Float64:
		numValue = source.Float()
	default:
		return source
	}

	// Convert to target type
	switch targetKind {
	case reflect.Int:
		return reflect.ValueOf(int(numValue))
	case reflect.Int8:
		return reflect.ValueOf(int8(numValue))
	case reflect.Int16:
		return reflect.ValueOf(int16(numValue))
	case reflect.Int32:
		return reflect.ValueOf(int32(numValue))
	case reflect.Int64:
		return reflect.ValueOf(int64(numValue))
	case reflect.Uint:
		return reflect.ValueOf(uint(numValue))
	case reflect.Uint8:
		return reflect.ValueOf(uint8(numValue))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(numValue))
	case reflect.Uint32:
		return reflect.ValueOf(uint32(numValue))
	case reflect.Uint64:
		return reflect.ValueOf(uint64(numValue))
	case reflect.Float32:
		return reflect.ValueOf(float32(numValue))
	case reflect.Float64:
		return reflect.ValueOf(numValue)
	}

	return source
}

func convertToString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Slice:
		return string(value.Interface().([]byte))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	default:
		return value.String()
	}
}
