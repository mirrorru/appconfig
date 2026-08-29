
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

	"github.com/delta-five/appconfig"
)

type httpCfg struct {
	Address string `default:":8080" flag:"addr" help:"Address to listen HTTP requests"`
	UseTLS  bool   `help:"Use TLS (HTTPS)"`
}

type appCfg struct {
	appconfig.ConfigBase `yaml:"-"` // Including fields with "magic" tags, if you need to process --help, --example or --config=file_name
	HTTP                 httpCfg
	Title                string `default:"My App" env:"name" flag:"name" help:"Name of application"`
	MagicNums            []int  `flag:"num" help:"Seed numbers"`
}

func main() {
	cfg := appCfg{}
	err := appconfig.Load(&cfg, appconfig.Params{EnvPrefix: "APP", FlagPrefix: "--"})

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
    main.appCfg{ConfigBase:appconfig.ConfigBase{ShowHelp:false, PrintExample:false, ConfigFile:""}, HTTP:main.httpCfg{Address:":8080", UseTLS:false}, Title:"My App", MagicNums:[]int(nil)}

#####  Showing help
    $ go run main.go --help
    List or program parameters
    Environment param              command-line flag              default value   description
    --help                         false           show this help
    --example                      false           show config example
    --config                                       config file to load
    APP_HTTP_ADDRESS               --http-addr                    :8080           Address to listen HTTP requests
    APP_HTTP_USE_TLS               --http-use-tls                                 Use TLS (HTTPS)
    APP_NAME                       --name                         My App          Name of application
    APP_MAGIC_NUMS                 --num                                          Seed numbers

#####  Showing config file example
    $ go run main.go --example
    ## >>>>> config file starts here >>>>>
    http:
        address: :8080
        usetls: false
    title: My App
    magicnums: []
## >>>>> config file ends here <<<<<<

#####  Load config from flags and environment
    $ APP_NAME="Best APP" go run main.go --http-addr=:8888 --http-use-tls --num=31 --num=61
    main.appCfg{ConfigBase:appconfig.ConfigBase{ShowHelp:false, PrintExample:false, ConfigFile:""}, HTTP:main.httpCfg{Address:":8888", UseTLS:true}, Title:"Best APP", MagicNums:[]int{31, 61}}


