package appconfig

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigInfo(t *testing.T) {
	t.Parallel()
	const PFX = "TST"
	type ForInclude struct {
		Help    bool   `env:"e1" flag:"f1" help:"h1" default:"d1" use_as_show_help_flag:"yes"`
		Example bool   `env:"e1" flag:"f1" help:"h1" default:"d1" use_as_example_printing_flag:"true"`
		Config  string `env:"e1" flag:"f1" help:"h1" default:"d1" use_as_config_file_name:"+"`
	}

	tests := []struct {
		name        string
		cfgReceiver any
		expectedCI  *ConfigInfo
		wantErr     bool
	}{
		{
			name:        "wrong type int",
			cfgReceiver: int(0),
			wantErr:     true,
		},
		{
			name:        "wrong type ptr",
			cfgReceiver: func() any { i := 100; return &i }(),
			wantErr:     true,
		},
		{
			name: "struct arg with use_as_show_help_flag",
			cfgReceiver: struct {
				Param   int  `env:"p" flag:"f" help:"h" default:"d"`
				Help    bool `env:"e1" flag:"f1" help:"h1" default:"d1" use_as_show_help_flag:"yes"`
				private ForInclude
			}{},
			expectedCI: &ConfigInfo{
				helpFlagParamNumber: 2,
				params: ParamList{
					{Path: "Param", EnvName: PFX + "_P", FlagName: "--f", HelpText: "h", Default: "d", index: []int{0}},
					{Path: "Help", EnvName: PFX + "_E1", FlagName: "--f1", HelpText: "h1", Default: "d1", index: []int{1}},
				},
			},
		},
		{
			name: "struct arg with use_as_example_printing_flag",
			cfgReceiver: struct {
				Param   int  `env:"p" flag:"f" help:"h" default:"d"`
				Example bool `env:"e1" flag:"f1" help:"h1" default:"d1" use_as_example_printing_flag:"true"`
				private ForInclude
			}{},
			expectedCI: &ConfigInfo{
				exampleFlagParamNumber: 2,
				params: ParamList{
					{Path: "Param", EnvName: PFX + "_P", FlagName: "--f", HelpText: "h", Default: "d", index: []int{0}},
					{Path: "Example", EnvName: PFX + "_E1", FlagName: "--f1", HelpText: "h1", Default: "d1", index: []int{1}},
				},
			},
		},
		{
			name: "struct arg with use_as_config_file_name",
			cfgReceiver: struct {
				Param   int    `env:"p" flag:"f" help:"h" default:"d"`
				Config  string `env:"e1" flag:"f1" help:"h1" default:"d1" use_as_config_file_name:"+"`
				private ForInclude
			}{},
			expectedCI: &ConfigInfo{
				configNameParamNumber: 2,
				params: ParamList{
					{Path: "Param", EnvName: PFX + "_P", FlagName: "--f", HelpText: "h", Default: "d", index: []int{0}},
					{Path: "Config", EnvName: PFX + "_E1", FlagName: "--f1", HelpText: "h1", Default: "d1", index: []int{1}},
				},
			},
		},
		{
			name: "Full",
			cfgReceiver: struct {
				ForInclude
				Sub struct {
					Fld struct {
						Param   int `env:"p" flag:"f" help:"h" default:"d"`
						private ForInclude
					}
					Bool  bool    `env:"p1" flag:"f1" help:"h1" default:"d1"`
					Str   string  `env:"p2" flag:"f2" help:"h2" default:"d2"`
					Float float64 `env:"p3" flag:"f3" help:"h3" default:"d3"`
				} `env:"se" flag:"sf"`
			}{},
			expectedCI: &ConfigInfo{
				helpFlagParamNumber:    1,
				exampleFlagParamNumber: 2,
				configNameParamNumber:  3,
				params: ParamList{
					{Path: "ForInclude.Help", EnvName: PFX + "_E1", FlagName: "--f1", HelpText: "h1", Default: "d1", index: []int{0, 0}},
					{Path: "ForInclude.Example", EnvName: PFX + "_E1", FlagName: "--f1", HelpText: "h1", Default: "d1", index: []int{0, 1}},
					{Path: "ForInclude.Config", EnvName: PFX + "_E1", FlagName: "--f1", HelpText: "h1", Default: "d1", index: []int{0, 2}},
					{Path: "Sub.Fld.Param", EnvName: PFX + "_SE_FLD_P", FlagName: "--sf-fld-f", HelpText: "h", Default: "d", index: []int{1, 0, 0}},
					{Path: "Sub.Bool", EnvName: PFX + "_SE_P1", FlagName: "--sf-f1", HelpText: "h1", Default: "d1", index: []int{1, 1}},
					{Path: "Sub.Str", EnvName: PFX + "_SE_P2", FlagName: "--sf-f2", HelpText: "h2", Default: "d2", index: []int{1, 2}},
					{Path: "Sub.Float", EnvName: PFX + "_SE_P3", FlagName: "--sf-f3", HelpText: "h3", Default: "d3", index: []int{1, 3}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ci, err := NewConfigInfo(tt.cfgReceiver, "TST")
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedCI, ci)
		})
	}
}

func TestConfigInfo_Load(t *testing.T) {
	t.Parallel()
	type SubCfg struct {
		SubValue float64
	}
	type TestCfg struct {
		ConfigBase
		Name    string `env:"path"`
		Slice   []int
		Map     map[string]string
		Value   int
		Flag    bool
		Include SubCfg
	}
	osArgsSrc := os.Args
	defer func() { os.Args = osArgsSrc }()

	m := sync.Mutex{}

	tests := []struct {
		name        string
		setup       func()
		sourceCfg   any
		expectedCfg TestCfg
		ci          *ConfigInfo
		wantErr     bool
	}{
		{
			name:      "fail on not ptr",
			sourceCfg: TestCfg{},
			wantErr:   true,
		},
		{
			name:      "fail on not struct ptr",
			ci:        &ConfigInfo{},
			sourceCfg: new(string),
			wantErr:   true,
		},
		{
			name: "invalid cfg file data",
			setup: func() {
				os.Args = append(osArgsSrc, "--config=test_cfg.invalid")
			},
			wantErr: true,
		},
		{
			name: "invalid cfg file name",
			setup: func() {
				os.Args = append(osArgsSrc, "--config=test_cfg.not_exist")
			},
			wantErr: true,
		},
		{
			name: "valid cfg file",
			setup: func() {
				os.Args = append(osArgsSrc, "--config=test_cfg.valid")
			},
			expectedCfg: TestCfg{
				ConfigBase: ConfigBase{
					ConfigFile: "test_cfg.valid",
				},
				Name:  "Somebody",
				Slice: []int{1, 2},
				Map:   map[string]string{"one": "two"},
				Value: 101,
				Flag:  true,
				Include: SubCfg{
					SubValue: 100.1,
				},
			},
		},
		{
			name: "help flag + name",
			setup: func() {
				os.Args = append(osArgsSrc, "--help")
			},
			expectedCfg: TestCfg{
				ConfigBase: ConfigBase{
					ShowHelp: true,
				},
				Name: os.Getenv("PATH"),
			},
		},
		{
			name: "example flag",
			setup: func() {
				os.Args = append(osArgsSrc, "--example", "--value=99")
			},
			expectedCfg: TestCfg{
				ConfigBase: ConfigBase{
					PrintExample: true,
				},
				Name:  os.Getenv("PATH"),
				Value: 99,
			},
		},
		{
			name: "path to name",
			setup: func() {
				os.Args = osArgsSrc
			},
			expectedCfg: TestCfg{
				Name: os.Getenv("PATH"),
			},
		},
		{
			name: "bad default",
			setup: func() {
				os.Args = osArgsSrc
			},
			sourceCfg: &struct {
				Value int `default:"string"`
			}{},
			wantErr: true,
		},
		{
			name: "bad flag value",
			setup: func() {
				os.Args = append(osArgsSrc, "--value=string")
			},
			sourceCfg: &struct {
				Value int
			}{},
			wantErr: true,
		},
		{
			name: "bad env value",
			sourceCfg: &struct {
				Value int `env:"PATH"`
			}{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m.Lock()
			defer m.Unlock()
			cfg := tt.sourceCfg
			if tt.sourceCfg == nil {
				cfg = &TestCfg{}
			}

			ci, err := NewConfigInfo(cfg, "")
			if tt.ci == nil {
				require.NoError(t, err)
			}

			if tt.setup != nil {
				tt.setup()
			}

			err = ci.Load(cfg)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, &tt.expectedCfg, cfg)

			if ci.HasExampleFlag() {
				require.NoError(t, ci.ShowExample(cfg))
			}

			if ci.HasHelpFlag() {
				ci.ShowHelp()
			}
		})
	}
}
