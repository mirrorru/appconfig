
Application configuration structure loader `appconfig`
======================================================

This library is designed to simplify loading the application configuration into structures.

Allows initializing simple fields via default values, command line flags, or environment variables. Complex structures (slices and maps) can be initialized via a configuration file.

It is possible to organize the output of a hint and an example configuration file.

General usage example:
```GO
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/mirrorru/appconfig"
)

type httpCfg struct {
	Address string `default:":8080" help:"Address to listen HTTP requests"`
	UseTLS  bool   `flag:"tls" help:"Use TLS (HTTPS)"`
}

type appCfg struct {
	appconfig.ConfigBase        // Including fields with "magic" tags, if you need to process --help, --example or --config=file_name
	Title                string `default:"My App" env:"name" flag:"name" help:"Name of application"`
	HTTP                 httpCfg
}

func main() {
	cfg := appCfg{}
	err := appconfig.Load(&cfg, "APP")

	if err != nil {
		if errors.Is(err, appconfig.ErrStopExpected) {
			os.Exit(0)
		}
		panic(err)
	}

	fmt.Printf("%#v\n", cfg)
}

```

#####  Just run      
    $ go run main.go
    main.appCfg{Title:"My App", HTTP:main.httpCfg{Address:":8080", UseTLS:false}, ConfigBase:appconfig.ConfigBase{ShowHelp:false, PrintExample:false, ConfigFile:""}}

#####  Showing help
    $ go run main.go --help
    List or program parameters
    Environment param              command-line flag              default value   description
    APP_NAME                       --name                         My App          Name of application
    APP_HTTP_ADDRESS               --http-address                 :8080           Address to listen HTTP requests
    APP_HTTP_USETLS                --http-tls                                     Use TLS (HTTPS)
                                   --help                         false           show this help
                                   --example                      false           show config example
                                   --config                                       config file to load

#####  Showing config file example
    $ go run main.go --example
    Config file example:
    ## >>>>> config file starts here >>>>>
    title: Best APP
    http:
    address: :8080
    usetls: false
    configbase: {}
    ## >>>>> config file ends here <<<<<<

#####  Load config from flags and environment
    $ APP_NAME="Best APP" go run main.go --http-address=:8888 --http-tls
    main.appCfg{Title:"Best APP", HTTP:main.httpCfg{Address:":8888", UseTLS:true}, ConfigBase:appconfig.ConfigBase{ShowHelp:false, PrintExample:false, ConfigFile:""}}


