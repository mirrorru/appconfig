# appconfig

##Application configuration

This library is designed to simplify loading the application configuration into structures.

Allows initializing simple fields via default values, command line flags, or environment variables. Complex structures (slices and maps) can be initialized via a configuration file.

It is possible to organize the output of a hint and an example configuration file.

General usage example:
```GO
    //...
	type httpServiceCfg struct {
        Address string `default:":8080" help:"Address to listen HTTP requests"`
        UseTLS  bool   `help:"Use TLS (HTTPS)"`
    }

    type appConfig struct {
        // Include appconfig.ConfigBase or specify "magic" tags, if you need to process --help, --example or --config=file_name
        Title string `default:"My App" env:"name" flag:"name" help:"Name of application"`
        HTTP  httpServiceCfg
    }

    cfg := appConfig{}
    err := appconfig.Load(&cfg, "APP")

    if errors.Is(err, appconfig.ErrStopExpected) {
        os.Exit(0)
    }
	
    if err != nil {
    panic(err)
    }
    //...
```
