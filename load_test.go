package appconfig_test

import (
	"errors"
	"fmt"
	"os"

	"github.com/mirrorru/appconfig"
)

func ExampleLoad() {
	type httpServiceCfg struct {
		Address string `default:":8080" flag:"addr" help:"Address to listen HTTP requests"`
		UseTLS  bool   `help:"Use TLS (HTTPS)"`
	}

	type appConfig struct {
		// Include appconfig.ConfigBase, if you need to process --help, --example or --config=file_name
		Title string `default:"My App" env:"name" flag:"name" help:"Name of application"`
		HTTP  httpServiceCfg
	}

	cfg := appConfig{}
	err := appconfig.Load(&cfg, "APP")

	if errors.Is(err, appconfig.ErrStopExpected) {
		os.Exit(-1)
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", cfg)

	// Output:
	// appconfig_test.appConfig{Title:"My App", HTTP:appconfig_test.httpServiceCfg{Address:":8080", UseTLS:false}}
}
