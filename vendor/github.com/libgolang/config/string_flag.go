package config

import "flag"

type stringFlagDef struct {
	name  string
	usage string

	valueDefault string
	valuePtr     *string
	value        string
}

func (f *stringFlagDef) Flag() {
	f.valuePtr = flag.String(f.name, envString(f.name, f.valueDefault), f.usage)
}

func (f *stringFlagDef) Resolve() {
	f.value = *f.valuePtr
}
