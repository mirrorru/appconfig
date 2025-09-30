package appconfig

import (
	"strings"
	"unicode"
)

// splitCamelCase разбивает строку по смене регистра
func splitCamelCase(s string) []string {
	if s == "" {
		return nil
	}

	var result []string
	runes := []rune(s)
	start := 0

	for i := 1; i < len(runes); i++ {
		if unicode.IsUpper(runes[i]) {
			// если текущая буква заглавная
			if unicode.IsLower(runes[i-1]) {
				// переход строчная -> заглавная
				result = append(result, string(runes[start:i]))
				start = i
			} else if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				// случай: несколько заглавных, затем строчная
				// пример: "DBMSKey" -> "DBMS", "Key"
				result = append(result, string(runes[start:i]))
				start = i
			}
		}
	}

	// добавляем последний кусок
	result = append(result, string(runes[start:]))

	return result
}

func toSnakeCase(s string) string {
	return strings.Join(splitCamelCase(s), EnvSeparator)
}

func toKebabCase(s string) string {
	return strings.Join(splitCamelCase(s), FlagSeparator)
}
