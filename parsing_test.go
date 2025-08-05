package appconfig

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBool(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{"true", "true", true, false},
		{"1", "1", true, false},
		{"yes", "yes", true, false},
		{"on", "on", true, false},
		{"false", "false", false, false},
		{"0", "0", false, false},
		{"no", "no", false, false},
		{"off", "off", false, false},
		{"invalid", "invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseBool(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{"positive", "123", 123, false},
		{"negative", "-456", -456, false},
		{"zero", "0", 0, false},
		{"invalid", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseInt(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseUint(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    uint64
		wantErr bool
	}{
		{"positive", "123", 123, false},
		{"zero", "0", 0, false},
		{"invalid_negative", "-456", 0, true},
		{"invalid_text", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseUint(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseFloat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{"positive", "123.45", 123.45, false},
		{"negative", "-456.78", -456.78, false},
		{"zero", "0", 0, false},
		{"invalid", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseFloat(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.InDelta(t, tt.want, got, 0.0001)
		})
	}
}

func TestParseFieldValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		field    any
		value    string
		expected any
		wantErr  bool
	}{
		{"string", "initial", "test", "test", false},
		{"bool_true", false, "true", true, false},
		{"bool_false", true, "false", false, false},
		{"bool_invalid", false, "invalid", false, true},
		{"bool_empty", false, "", true, false},
		{"int", int(0), "123", int(123), false},
		{"int8", int8(0), "123", int8(123), false},
		{"int_invalid", int(0), "abc", int(0), true},
		{"uint", uint(0), "123", uint(123), false},
		{"uint8", uint8(0), "123", uint8(123), false},
		{"uint_invalid", uint(0), "abc", uint(0), true},
		{"float32", float32(0), "123.45", float32(123.45), false},
		{"float64", float64(0), "123.45", float64(123.45), false},
		{"float_invalid", float64(0), "abc", float64(0), true},
		{"unsupported_type", []string{}, "test", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			field := reflect.New(reflect.TypeOf(tt.field)).Elem()
			err := parseFieldValue(field, tt.value)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, field.Interface())
		})
	}
}

func TestParseFlags(t *testing.T) {
	t.Parallel()
	// Сохраняем оригинальные аргументы и восстанавливаем их после теста
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name     string
		args     []string
		expected map[string]string
	}{
		{
			"no_args",
			[]string{},
			map[string]string{},
		},
		{
			"single_flag",
			[]string{"--flag=value"},
			map[string]string{"--flag": "value"},
		},
		{
			"multiple_flags",
			[]string{"--flag1=value1", "--flag2=value2"},
			map[string]string{"--flag1": "value1", "--flag2": "value2"},
		},
		{
			"flag_without_value",
			[]string{"--flag"},
			map[string]string{"--flag": ""},
		},
		{
			"mixed_flags",
			[]string{"--flag1=value", "--flag2"},
			map[string]string{"--flag1": "value", "--flag2": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			os.Args = tt.args
			result := parseFlags(tt.args)
			assert.Equal(t, tt.expected, result)
		})
	}
}
