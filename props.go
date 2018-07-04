package config

import (
	"fmt"
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
	configFlagName := fmt.Sprintf(ConfigFlagName, AppName) // ConfigFlagName: %s_config -> APPNAME_CONFIG
	envFileName, ok := os.LookupEnv(_env(configFlagName))
	if !ok { // if app_config env variable not set:
		envFileName = fmt.Sprintf(ConfigFileNameFormat, AppName) // ConfigFileNameFormat %s.properties
		if _, err := os.Stat(envFileName); os.IsNotExist(err) {  // if not exists:
			cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
			envFileName = path.Join(cwd, envFileName) // $PWD/config.properties
		}
	}

	p, err := properties.LoadFile(envFileName, properties.UTF8)
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
