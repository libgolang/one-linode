package config

import "flag"

type int64FlagDef struct {
	name  string
	usage string

	valueDefault int64
	valuePtr     *int64
	value        int64
}

func (f *int64FlagDef) Flag() {
	f.valuePtr = flag.Int64(f.name, envInt64(f.name, f.valueDefault), f.usage)
}

func (f *int64FlagDef) Resolve() {
	f.value = *f.valuePtr
}
