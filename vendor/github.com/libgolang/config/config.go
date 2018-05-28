package config

import (
	"flag"

	"github.com/magiconair/properties"
)

var (
	// ConfigFlagName sets the flat to use when passing the config file name
	ConfigFlagName = "config"
	// DefaultConfigFile is the default name where configuration is pulled from
	DefaultConfigFile = "config.properties"
	flagMap           = make(map[string]FlagDef)
	loadedProps       *properties.Properties
)

// FlagDef interface
type FlagDef interface {
	Flag()
	Resolve()
}

type flagDef struct {
	name  string
	usage string

	stringDefault  string
	stringValuePtr *string
	stringValue    string

	intValue    int
	intValuePtr *int

	int64Value    int64
	int64ValuePtr *int64

	boolValue    bool
	boolValuePtr *bool

	floatValue    float64
	floatValuePtr *float64
}

// String defines a string flag
func String(name, def, usage string) *string {
	fd := &stringFlagDef{
		name:         name,
		usage:        usage,
		valueDefault: def,
	}
	flagMap[name] = fd
	return &fd.value
}

// Int defines an integer flag
func Int(name string, def int, usage string) *int {
	fd := &intFlagDef{
		name:         name,
		usage:        usage,
		valueDefault: def,
	}
	flagMap[name] = fd
	return &fd.value
}

// Int64 defines an integer 64 flag
func Int64(name string, def int64, usage string) *int64 {
	fd := &int64FlagDef{
		name:         name,
		usage:        usage,
		valueDefault: def,
	}
	flagMap[name] = fd
	return &fd.value
}

// Float defines a float flag
func Float(name string, def float64, usage string) *float64 {
	fd := &floatFlagDef{
		name:         name,
		usage:        usage,
		valueDefault: def,
	}
	flagMap[name] = fd
	return &fd.value
}

// Bool defines a boolean flag
func Bool(name string, def bool, usage string) *bool {
	fd := &boolFlagDef{
		name:         name,
		usage:        usage,
		valueDefault: def,
	}
	flagMap[name] = fd
	return &fd.value
}

// Var custom flag parsing.  see flag.Var
func Var(v FlagValue, name string, usage string) {
	fd := &varFlagDef{
		name:  name,
		usage: usage,
		value: v,
	}
	flagMap[name] = fd
}

// Parse call parse on flags
func Parse() {
	for _, v := range flagMap {
		v.Flag()
	}
	flag.Parse()
	for _, v := range flagMap {
		v.Resolve()
	}
}

// PrintHelp prints flag helps
func PrintHelp() {
	flag.PrintDefaults()
}
