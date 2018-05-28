package config

import (
	"flag"
	"os"
)

type varFlagDef struct {
	name      string
	usage     string
	setCalled bool
	value     FlagValue
}

func (f *varFlagDef) Flag() {
	flag.Var(f, f.name, f.usage)
	f.setCalled = false
}

func (f *varFlagDef) Resolve() {
	if !f.setCalled {
		if v, ok := os.LookupEnv(_env(f.name)); ok {
			f.value.FromString(v)
		} else if v, ok := getConfig().Get(f.name); ok {
			f.value.FromString(v)
		}
	}
}

func (f *varFlagDef) String() string {
	return f.value.ToString()
}
func (f *varFlagDef) Set(value string) error {
	f.setCalled = true
	return f.value.FromString(value)
}

// FlagValue interface
type FlagValue interface {
	ToString() string
	FromString(value string) error
}
