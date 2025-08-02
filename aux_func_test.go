package appconfig

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddPrefix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		inputName string
		inputPref string
		inputSep  string
		expected  string
	}{
		{
			name:      "empty name",
			inputName: "",
			inputPref: "pref",
			inputSep:  ":",
			expected:  "",
		},
		{
			name:      "empty prefix",
			inputName: "name",
			inputPref: "",
			inputSep:  ":",
			expected:  "name",
		},
		{
			name:      "both non-empty",
			inputName: "name",
			inputPref: "pref",
			inputSep:  ":",
			expected:  "pref:name",
		},
		{
			name:      "custom separator",
			inputName: "name",
			inputPref: "pref",
			inputSep:  "-",
			expected:  "pref-name",
		},
		{
			name:      "empty separator",
			inputName: "name",
			inputPref: "pref",
			inputSep:  "",
			expected:  "prefname",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := addPrefix(tt.inputName, tt.inputPref, tt.inputSep)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetTagOrName(t *testing.T) {
	t.Parallel()
	type testStruct struct {
		Field1 string `test:"field1"`
		Field2 string `test:"-"`
		Field3 string `test:""`
		Field4 string `test:"custom"`
		Field5 string
	}

	fields := []reflect.StructField{}
	rt := reflect.TypeOf(testStruct{})
	for i := 0; i < rt.NumField(); i++ {
		fields = append(fields, rt.Field(i))
	}

	tests := []struct {
		name     string
		tag      string
		field    *reflect.StructField
		expected string
	}{
		{
			name:     "tag exists",
			tag:      "test",
			field:    &fields[0],
			expected: "field1",
		},
		{
			name:     "tag with minus",
			tag:      "test",
			field:    &fields[1],
			expected: "",
		},
		{
			name:     "empty tag value",
			tag:      "test",
			field:    &fields[2],
			expected: "Field3",
		},
		{
			name:     "custom tag value",
			tag:      "test",
			field:    &fields[3],
			expected: "custom",
		},
		{
			name:     "no tag",
			tag:      "test",
			field:    &fields[4],
			expected: "Field5",
		},
		{
			name:     "non-existent tag",
			tag:      "nonexistent",
			field:    &fields[0],
			expected: "Field1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := getTagOrName(tt.tag, tt.field)
			require.Equal(t, tt.expected, result)
		})
	}
}
