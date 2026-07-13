package appconfig

import (
	"fmt"
	"reflect"
	"strings"
)

func parseFieldValue(field reflect.Value, values []string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(values[0])
	case reflect.Bool:
		if values[0] == "" {
			// in case then flag is "--enableSomething"
			field.SetBool(true)
			return nil
		}
		boolValue, err := parseBool(values[0])
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := parseInt(values[0])
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := parseUint(values[0])
		if err != nil {
			return err
		}
		field.SetUint(uintValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := parseFloat(values[0])
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)
	case reflect.Slice:
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

func parseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", s)
	}
}

func parseInt(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func parseUint(s string) (uint64, error) {
	var u uint64
	_, err := fmt.Sscanf(s, "%d", &u)
	return u, err
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

func parseFlags(args1toN []string) map[string][]string {
	result := map[string][]string{}
	for _, arg := range args1toN {
		arr := strings.SplitN(arg, "=", 2)
		key := arr[0]
		val := ""
		if len(arr) > 1 {
			val = arr[1]
		}
		vals, ok := result[key]
		if !ok {
			result[key] = []string{val}

			continue
		}
		result[key] = append(vals, val)
	}

	return result
}
