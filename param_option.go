package appconfig

import "strings"

type Option func(pi *ParamInfo)

var defaultOpts = []Option{
	OptEnvToUpper,
	OptFlagNameDoubleDash,
	OptFlagNameToLower,
}
var OptEnvToUpper Option = func(pi *ParamInfo) {
	pi.EnvName = strings.ToUpper(pi.EnvName)
}

var OptFlagNameDoubleDash Option = func(pi *ParamInfo) {
	pi.FlagName = "--" + pi.FlagName
}

var OptFlagNameToLower Option = func(pi *ParamInfo) {
	pi.FlagName = strings.ToLower(pi.FlagName)
}
