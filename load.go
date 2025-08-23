package appconfig

import "errors"

var (
	ErrStopExpected = errors.New(`a stop is expected`)
	ErrHelpShown    = errors.Join(ErrStopExpected, errors.New("help shown"))
	ErrExampleShown = errors.Join(ErrStopExpected, errors.New("example shown"))
)

// Load - loads field values from defaults, then from environment, when from flags, when from config, if specified
//   - config - a pointer to structure where the configuration is planned to be loaded
func Load[T any, PT interface{ *T }](receiver PT, envPrefix string) (errResult error) {
	ci, err := NewConfigInfo(receiver, envPrefix)
	if err != nil {
		return err
	}

	if err = ci.Load(receiver); err != nil {
		return err
	}

	if ci.HasHelpFlag() {
		ci.ShowHelp()
		errResult = errors.Join(errResult, ErrHelpShown)
	}

	if ci.HasExampleFlag() {
		if errLocal := ci.ShowExample(receiver); errLocal != nil {
			return errLocal
		}
		errResult = errors.Join(errResult, ErrExampleShown)
	}

	return errResult
}

// MustLoad - try to Load configuration, and panics if error!=nil
func MustLoad[T any, PT interface{ *T }](receiver PT, envPrefix string) {
	if err := Load(receiver, envPrefix); err != nil {
		panic(err)
	}
}
