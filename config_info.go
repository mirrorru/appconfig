package appconfig

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type ConfigInfo struct {
	params                 ParamList
	helpFlagParamNumber    int
	helpFlagParamValue     bool
	exampleFlagParamNumber int
	exampleFlagParamValue  bool
	configNameParamNumber  int
	configNameParamValue   string
}

// NewConfigInfo creates new item on ConfigInfo and fills it with information of config parameters from `config`
//   - config - any structure or a pointer to it where the configuration is planned to be loaded
//   - envPrefix - a common prefix for environment variables from which configuration values can be taken
func NewConfigInfo(config any, envPrefix string) (result *ConfigInfo, err error) {
	rv := reflect.ValueOf(config)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, errors.New("value is not a struct or pointer to struct")
	}

	result = new(ConfigInfo)
	result.processType(rv.Type(), "", envPrefix, "", nil)
	for idx := range result.params {
		if result.params[idx].EnvName != "" {
			result.params[idx].EnvName = strings.ToUpper(result.params[idx].EnvName)
		}
		if result.params[idx].FlagName != "" {
			result.params[idx].FlagName = "--" + strings.ToLower(result.params[idx].FlagName)
		}
	}

	return
}

func (ci *ConfigInfo) processType(t reflect.Type, pathPrefix string, envPrefix string, flagPrefix string, indexes []int) {
fieldsLoop:
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue // Пропускаем неэкспортируемые поля
		}

		if field.Type.Kind() == reflect.Struct {
			subEnvPrefix := envPrefix
			subFlagPrefix := flagPrefix
			if !field.Anonymous {
				subEnvPrefix = addPrefix(getTagOrName("env", &field), envPrefix, "_")
				subFlagPrefix = addPrefix(getTagOrName("flag", &field), flagPrefix, "-")
			}
			subPathPrefix := addPrefix(field.Name, pathPrefix, ".")
			ci.processType(field.Type, subPathPrefix, subEnvPrefix, subFlagPrefix, append(indexes, field.Index...))

			continue fieldsLoop
		}

		pi := ParamInfo{
			Path:     addPrefix(field.Name, pathPrefix, "."),
			EnvName:  addPrefix(getTagOrName("env", &field), envPrefix, "_"),
			FlagName: addPrefix(getTagOrName("flag", &field), flagPrefix, "-"),
			HelpText: getTagOrName("help", &field),
			Default:  field.Tag.Get("default"),
			index:    append(indexes, field.Index...),
		}

		ci.params = append(ci.params, pi)
		if field.Tag.Get("use_as_show_help_flag") != "" && field.Type.Kind() == reflect.Bool {
			ci.helpFlagParamNumber = len(ci.params) // after append
		}
		if field.Tag.Get("use_as_example_printing_flag") != "" && field.Type.Kind() == reflect.Bool {
			ci.exampleFlagParamNumber = len(ci.params) // after append
		}
		if field.Tag.Get("use_as_config_file_name") != "" && field.Type.Kind() == reflect.String {
			ci.configNameParamNumber = len(ci.params) // after append
		}
	}
}

type loadSource byte

const (
	LoadSourceDefaults loadSource = iota
	LoadSourceFlags
	LoadSourceEnvs
)

// LoadInOrder - loads field values from specified source order, can be used for init default config.
//   - config - a pointer to structure where the configuration is planned to be loaded
func (ci *ConfigInfo) LoadInOrder(config any, order ...loadSource) error {
	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Ptr {
		return errors.New("value is not a pointer to struct")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("value is not a pointer to struct")
	}

	var flags map[string]string
	if slices.Contains(order, LoadSourceFlags) {
		flags = parseFlags(os.Args[1:])
	}

	for idx, param := range ci.params {
		field := rv.FieldByIndex(param.index)
		for _, source := range order {
			switch source {
			case LoadSourceDefaults:
				if param.Default != "" {
					if err := parseFieldValue(field, param.Default); err != nil {
						return fmt.Errorf("can't parse default value `%s` for %s: %w", param.Default, param.Path, err)
					}
				}
			case LoadSourceEnvs:
				if param.EnvName != "" {
					if envValue, exists := os.LookupEnv(param.EnvName); exists && envValue != "" {
						if err := parseFieldValue(field, envValue); err != nil {
							return fmt.Errorf("can't parse env value `%s` for %s: %w", envValue, param.Path, err)
						}
					}
				}
			case LoadSourceFlags:
				if param.FlagName != "" {
					if flagValue, exists := flags[param.FlagName]; exists {
						if err := parseFieldValue(field, flagValue); err != nil {
							return fmt.Errorf("can't parse flag value `%s` for %s: %w", flagValue, param.Path, err)
						}
					}
				}
			}
		}

		if idx+1 == ci.helpFlagParamNumber {
			ci.helpFlagParamValue = field.Bool()
		}
		if idx+1 == ci.exampleFlagParamNumber {
			ci.exampleFlagParamValue = field.Bool()
		}
		if idx+1 == ci.configNameParamNumber {
			ci.configNameParamValue = field.String()
		}
	}

	return nil
}

// TryLoadConfigFile - loads field values from config-file, if specified in ConfigInfo
//   - config - a pointer to structure where the configuration is planned to be loaded
func (ci *ConfigInfo) TryLoadConfigFile(config any) error {
	if ci.configNameParamValue == "" {
		return nil
	}

	// Load config from file
	// Читаем файл
	data, err := os.ReadFile(ci.configNameParamValue)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	if err = yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to unmarshal config file: %v", err)
	}

	return nil
}

// DefaultLoadOrder - default param-source order for loading in Load method
var DefaultLoadOrder = []loadSource{LoadSourceDefaults, LoadSourceFlags, LoadSourceEnvs}

// Load - loads field values from defaults, then from environment, when from flags, when from config, if specified
//   - config - a pointer to structure where the configuration is planned to be loaded
func (ci *ConfigInfo) Load(config any) error {
	if err := ci.LoadInOrder(config, DefaultLoadOrder...); err != nil {
		return err
	}

	return ci.TryLoadConfigFile(config)
}

// HasHelpFlag checks that the "help" flag is set
func (ci *ConfigInfo) HasHelpFlag() bool {
	return ci.helpFlagParamValue
}

// HasExampleFlag checks that the "example" flag is set
func (ci *ConfigInfo) HasExampleFlag() bool {
	return ci.exampleFlagParamValue
}

// ShowHelp showing help
func (ci *ConfigInfo) ShowHelp() {
	// Showing parameters help
	const lineFormat = "%-30s %-30s %-15s %s\n"
	fmt.Println("List or program parameters")
	_, _ = fmt.Printf(lineFormat, "Environment param", "command-line flag", "default value", "description")
	for _, param := range ci.params {
		fmt.Printf(lineFormat, param.EnvName, param.FlagName, param.Default, param.HelpText)
	}
}

// ShowExample showing config example based on `config` data
func (ci *ConfigInfo) ShowExample(config any) error {
	// printing config file example
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config file for printing: %v", err)
	}
	fmt.Printf("Config file example:\n## >>>>> config file starts here >>>>>\n%s## >>>>> config file ends here <<<<<<\n", string(data))

	return nil
}
