package appconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCamelCase(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  []string
	}{
		{"", nil},
		{"Any", []string{"Any"}},
		{"AnyKey", []string{"Any", "Key"}},
		{"AnyDBMS", []string{"Any", "DBMS"}},
		{"anyDBMS", []string{"any", "DBMS"}},
		{"anyDBMSKey", []string{"any", "DBMS", "Key"}},
		{"MyLongDBNameForSQL", []string{"My", "Long", "DB", "Name", "For", "SQL"}},
		{"", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got := splitCamelCase(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"Any", "Any"},
		{"OneTwo", "One_Two"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got := toSnakeCase(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToKebabCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"Any", "Any"},
		{"OneTwo", "One-Two"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got := toKebabCase(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
