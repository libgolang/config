package config

//
import (
	"github.com/magiconair/properties"
	flag "github.com/spf13/pflag"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode"
)

//
var (
	conf = NewConf()
)

// Conf config structure
type Conf struct {
	isInitialized bool
	name          string
	props         *properties.Properties
	flagMap       map[string]flagMapType
}

// NewConf *Conf constructor
func NewConf() *Conf {
	return &Conf{
		isInitialized: false,
		name:          "app",
		flagMap:       make(map[string]flagMapType),
	}
}

// DefineString defines a flag
func DefineString(name string, def string, usage string) { conf.DefineString(name, def, usage) }

// DefineBool defines a bool flag
func DefineBool(name string, def bool, usage string) { conf.DefineBool(name, def, usage) }

// DefineInt defines an int flag
func DefineInt(name string, def int, usage string) { conf.DefineInt(name, def, usage) }

// SetName sets the configuration name used as
// prefix and configruation file
func SetName(name string) { conf.SetName(name) }

// SetName sets the configuration name used as
// prefix and configruation file
func (c *Conf) SetName(name string) {
	c.name = name
}

// GetString get configuration string
func GetString(name string) string { return conf.GetString(name) }

// GetInt get configuration int
func GetInt(name string) int { return conf.GetInt(name) }

// GetBool get configuration bool
func GetBool(name string) bool { return conf.GetBool(name) }

// GetString get configuration string
func (c *Conf) GetString(name string) string {
	//
	c.init()
	return c.findConfig(name).String()
}

// GetInt get configuration int
func (c *Conf) GetInt(name string) int {
	c.init()
	return c.findConfig(name).Int()
}

// GetBool get configuration bool
func (c *Conf) GetBool(name string) bool {
	c.init()
	return c.findConfig(name).Bool()
}

func (c *Conf) findConfig(name string) flagMapType {
	ptr, found := c.flagMap[strings.ToLower(name)]
	if !found {
		panic("Configuration parameter \"" + name + "\" does not exist.  It must be defined first.")
	}
	return ptr
}

func (c *Conf) init() {
	//
	if c.isInitialized {
		return
	}

	// Flags
	configFilePtr := flag.String("config", "", "Configuration File")
	flag.Parse()

	//Env
	if *configFilePtr == "" {
		if configFile, found := os.LookupEnv(c.toEnvKey("config")); found {
			*configFilePtr = configFile
		}
	}

	//Config Name
	if *configFilePtr == "" {
		var cwd string
		if dir, err := os.Getwd(); err != nil {
			cwd = dir
		} else {
			cwd = "./"
		}
		files := []string{
			path.Join(cwd, c.name+".properties"),
			path.Join("${HOME}", c.name+".properties"),
			path.Join("/etc", c.name+".properties"),
		}
		properties.LogPrintf = func(fmt string, args ...interface{}) {} // no logging please
		if props, err := properties.LoadAll(files, properties.UTF8, true); err != nil {
			c.props = properties.NewProperties()
		} else {
			c.props = props
		}
	} else {
		c.props = properties.MustLoadAll([]string{*configFilePtr}, properties.UTF8, true)
	}

	//
	c.isInitialized = true
}

// DefineString defines a flag on the config struct
func (c *Conf) DefineString(name string, def string, usage string) {
	f := &flagMapTypeString{
		parent: c,
		name:   name,
		usage:  usage,
		defVal: def,
		isSet:  false,
	}

	c.flagMap[strings.ToLower(name)] = f
	flag.Var(f, name, usage)
}

// DefineInt defines an int flag on the config struct
func (c *Conf) DefineInt(name string, def int, usage string) {
	c.DefineString(name, strconv.Itoa(def), usage)
}

// DefineBool defines a bool flag on the config struct
func (c *Conf) DefineBool(name string, def bool, usage string) {
	c.DefineString(name, strconv.FormatBool(def), usage)
}

func (c *Conf) toEnvKey(s string) string {
	newStr := make([]rune, 0, len(s))
	for i, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if unicode.IsLower(r) {
				r = unicode.ToUpper(r)
			}
			newStr = append(newStr, r)
		} else {
			// do not add _ if the previous one was _
			if i > 0 && newStr[i-1] == '_' {
				continue
			}
			newStr = append(newStr, '_')
		}
	}
	return strings.ToUpper(c.name) + "_" + string(newStr)
}

type flagMapType interface {
	String() string
	Bool() bool
	Int() int
}

type flagMapTypeString struct {
	parent *Conf
	name   string
	usage  string
	defVal string
	val    *string
	isSet  bool
}

func (f *flagMapTypeString) Set(v string) error {
	f.val = &v
	f.isSet = true
	return nil
}

func (f *flagMapTypeString) Type() string {
	return "string"
}

func (f *flagMapTypeString) String() string {

	// Return default value if not initialized.
	// this is to make it compatible with pflag
	// otherwise we get an NPE
	if !f.parent.isInitialized {
		return f.defVal
	}

	// Flags
	if f.isSet {
		return *f.val
	}

	// Env
	envKey := f.parent.toEnvKey(f.name)
	if val, found := os.LookupEnv(envKey); found {
		return val
	}

	// Get value from file
	return f.parent.props.GetString(f.name, f.defVal)

}

func (f *flagMapTypeString) Int() int {
	i, err := strconv.Atoi(*f.val)
	if err != nil {
		return 0
	}
	return i
}

func (f *flagMapTypeString) Bool() bool {
	v, err := strconv.ParseBool(*f.val)
	if err != nil {
		return false
	}
	return v
}
