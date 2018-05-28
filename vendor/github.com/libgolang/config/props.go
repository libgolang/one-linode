package config

import (
	"flag"
	"os"
	"path"
	"path/filepath"

	"github.com/magiconair/properties"
)

func getConfig() *properties.Properties {

	if loadedProps != nil {
		return loadedProps
	}

	// get it from environment
	envFileName, ok := os.LookupEnv(_env(ConfigFlagName))
	if !ok {
		envFileName = DefaultConfigFile
		if _, err := os.Stat(envFileName); os.IsNotExist(err) {
			cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
			envFileName = path.Join(cwd, envFileName)
		}
	}

	// get it from flags
	fs := flag.NewFlagSet(ConfigFlagName, flag.ContinueOnError)
	fs.SetOutput(&dummyWriter{})
	configFilePtr := fs.String(ConfigFlagName, envFileName, "Path to configuration file")
	_ = fs.Parse(os.Args[1:])
	p, err := properties.LoadFile(*configFilePtr, properties.UTF8)
	if err != nil {
		p = properties.NewProperties()
	}
	loadedProps = p // set global variable
	return p
}

//
type dummyWriter struct {
}

//
func (w *dummyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
