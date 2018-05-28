package config

import "flag"

type intFlagDef struct {
	name  string
	usage string

	valueDefault int
	valuePtr     *int
	value        int
}

func (f *intFlagDef) Flag() {
	f.valuePtr = flag.Int(f.name, envInt(f.name, f.valueDefault), f.usage)
}

func (f *intFlagDef) Resolve() {
	f.value = *f.valuePtr
}
