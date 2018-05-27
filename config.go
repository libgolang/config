package config

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/magiconair/properties"
)

var (
	emptyString = ""
	flagMap     = make(map[string]*flagDef)
	loadedProps *properties.Properties
)

type flagDef struct {
	name     string
	usage    string
	defValue string

	valuePtr    *string
	stringValue string
	intValue    int64
	boolValue   bool
	floatValue  float64
}

// String defines a string flag
func String(name, defValue, usage string) *string {
	def := flagVar(name, defValue, usage)
	return &def.stringValue
}

// Int defines an integer flag
func Int(name string, defValue int, usage string) *int64 {
	def := flagVar(name, strconv.Itoa(defValue), usage)
	return &def.intValue
}

// Float defines an integer flag
func Float(name string, defValue float64, usage string) *float64 {
	def := flagVar(name, strconv.FormatFloat(defValue, 'f', -1, 64), usage)
	return &def.floatValue
}

// Bool defines an integer flag
func Bool(name string, defValue bool, usage string) *int64 {
	def := flagVar(name, strconv.FormatBool(defValue), usage)
	return &def.intValue
}

func flagVar(name, defValue, usage string) *flagDef {
	props := getConfig()
	def := &flagDef{
		name:     name,
		defValue: props.GetString(name, defValue),
		usage:    usage,
	}
	flagMap[name] = def
	return def
}

// Parse call parse on flags
func Parse() {
	for _, v := range flagMap {
		v.valuePtr = flag.String(v.name, v.defValue, v.usage)
	}
	flag.Parse()
	for _, v := range flagMap {
		v.stringValue = *v.valuePtr
		v.intValue, _ = strconv.ParseInt(v.stringValue, 10, 64)
		v.floatValue, _ = strconv.ParseFloat(v.stringValue, 64)
		v.boolValue, _ = strconv.ParseBool(v.stringValue)
	}
}

// PrintHelp prints flag helps
func PrintHelp() {
	flag.PrintDefaults()
}

//
type voidWriter struct {
}

func (w *voidWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func getConfig() *properties.Properties {

	if loadedProps != nil {
		return loadedProps
	}

	// get it from environment
	envFileName, ok := os.LookupEnv("CONFIG")
	if !ok {
		envFileName = "config.properties"
		if _, err := os.Stat(envFileName); os.IsNotExist(err) {
			cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
			envFileName = path.Join(cwd, envFileName)
		}
	}

	// get it from flags
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.SetOutput(&voidWriter{})
	configFilePtr := fs.String("config", envFileName, "Path to configuration file")
	_ = fs.Parse(os.Args[1:])
	p, err := properties.LoadFile(*configFilePtr, properties.UTF8)
	if err != nil {
		p = properties.NewProperties()
	}
	loadedProps = p // set global variable
	return p
}
