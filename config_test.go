package config

import (
	flag "github.com/spf13/pflag"
	"os"
	"testing"
)

func TestEnv(t *testing.T) {

	// given
	appName := "xyz"
	key := "key"
	envKey := "XYZ_KEY"
	val := "val-a"
	_ = os.Setenv(envKey, val)
	defer func() {
		_ = os.Unsetenv(envKey)
	}()

	// when
	conf := NewConf()
	conf.SetName(appName)
	conf.DefineString(key, "", "")
	res := conf.GetString(key)

	// then
	if conf.name != appName {
		t.Error("Expected app name", appName)
	}

	if res != val {
		t.Errorf("Exptected \"%s\", but got \"%s\"\n", val, res)
	}
}

func TestParams(t *testing.T) {
	// given
	appName := "xyz"
	key := "xyz"
	expected := "xyz-val"

	// when
	// backup flags
	origOsArgs := os.Args
	origCommandLine := flag.CommandLine

	// override flags
	os.Args = []string{"./test", "--xyz", expected}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// restore flags
	defer func() {
		flag.CommandLine = origCommandLine
		os.Args = origOsArgs
	}()

	//
	conf := NewConf()
	conf.SetName(appName)
	conf.DefineString(key, "defalt value", "Description")
	res := conf.GetString(key)

	// then
	if res != expected {
		t.Errorf("Exptected \"%s\", but got \"%s\"\n", expected, res)
	}
}

func TestPropertiesFile(t *testing.T) {
	// given
	appName := "testyml"
	key := "test.a"
	expected := "abc"

	// backup flags
	origOsArgs := os.Args
	origCommandLine := flag.CommandLine

	// override flags
	os.Args = []string{"./test"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// restore flags
	defer func() {
		flag.CommandLine = origCommandLine
		os.Args = origOsArgs
	}()

	// when
	conf := NewConf()
	conf.SetName(appName)
	conf.DefineString(key, "", "")
	res := conf.GetString(key)

	// then
	if res != expected {
		t.Errorf("Exptected \"%s\", but got \"%s\"\n", expected, res)
	}
}

func TestConfigOverrideParamFile(t *testing.T) {
	// given
	appName := "testyml"
	key := "test.override.abc"
	expected := "abc"

	// backup flags
	origOsArgs := os.Args
	origCommandLine := flag.CommandLine

	// override flags
	os.Args = []string{"./test", "--config", "testoverride.properties"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// restore flags
	defer func() {
		flag.CommandLine = origCommandLine
		os.Args = origOsArgs
	}()

	// when
	conf := NewConf()
	conf.SetName(appName)
	conf.DefineString(key, "", "")
	res := conf.GetString(key)

	// then
	if res != expected {
		t.Errorf("Exptected \"%s\", but got \"%s\"\n", expected, res)
	}
}

func TestConfigOverrideEnvFile(t *testing.T) {
	// given
	appName := "ENVCONFIG"
	key := "test.override.abc"
	expected := "abc"
	envKey := "ENVCONFIG_CONFIG"

	// backup flags
	origOsArgs := os.Args
	origCommandLine := flag.CommandLine
	_ = os.Setenv(envKey, "testoverride.properties")

	// override flags
	os.Args = []string{"./test"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// restore flags
	defer func() {
		flag.CommandLine = origCommandLine
		os.Args = origOsArgs
		_ = os.Unsetenv(envKey)
	}()

	// when
	conf := NewConf()
	conf.SetName(appName)
	conf.DefineString(key, "", "")
	res := conf.GetString(key)

	// then
	if res != expected {
		t.Errorf("Exptected \"%s\", but got \"%s\"\n", expected, res)
	}
}
