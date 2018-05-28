package config

import "flag"

type boolFlagDef struct {
	name  string
	usage string

	valueDefault bool
	valuePtr     *bool
	value        bool
}

func (f *boolFlagDef) Flag() {
	f.valuePtr = flag.Bool(f.name, envBool(f.name, f.valueDefault), f.usage)
}

func (f *boolFlagDef) Resolve() {
	f.value = *f.valuePtr
}
