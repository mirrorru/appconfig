package appconfig

// ConfigBase can be used as embedded field in configuration structure with predefined parameters with autoprocessing:
//
// - `help` to use as showing help flag
//
// - `example`to use as printing config example flag
//
// - `config` to specify yaml-config file for loading
type ConfigBase struct {
	ShowHelp     bool   `yaml:"-" json:"-" env:"-" flag:"help"    default:"false" help:"show this help"      use_as_show_help_flag:"yes"`
	PrintExample bool   `yaml:"-" json:"-" env:"-" flag:"example" default:"false" help:"show config example" use_as_example_printing_flag:"yes"`
	ConfigFile   string `yaml:"-" json:"-" env:"-" flag:"config"  default:""     help:"config file to load"  use_as_config_file_name:"yes"`
}

type ParamInfo struct {
	Path     string
	EnvName  string
	FlagName string
	HelpText string
	Default  string
	index    []int
}

type ParamList []ParamInfo
