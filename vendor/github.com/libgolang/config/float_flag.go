package config

import "flag"

type floatFlagDef struct {
	name  string
	usage string

	valueDefault float64
	valuePtr     *float64
	value        float64
}

func (f *floatFlagDef) Flag() {
	f.valuePtr = flag.Float64(f.name, envFloat(f.name, f.valueDefault), f.usage)
}

func (f *floatFlagDef) Resolve() {
	f.value = *f.valuePtr
}
