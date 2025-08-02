package appconfig

import (
	"reflect"
)

func addPrefix(name string, prefix string, separator string) string {
	if name == "" || prefix == "" {
		return name
	}
	return prefix + separator + name
}

func getTagOrName(tag string, field *reflect.StructField) string {
	result := field.Tag.Get(tag)
	switch result {
	case "":
		return field.Name
	case "-":
		return ""
	default:
		return result
	}
}
