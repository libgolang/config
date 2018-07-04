# Configuration Library

## Usage


```
package main

import (
	"fmt"

	"github.com/libgolang/config"
)

func main() {
	config.AppName = "my-application"
	strFlag := config.String("flag-1", "default!", "Flag 1 usage")
	boolFlag := config.Bool("flag.2", false, "Flag 2 usage")
	config.Parse()

	fmt.Printf("flag-1: %s\n", *strFlag)
	fmt.Printf("flag.2: %t\n", *boolFlag)
}
```

The code will read configuration variables like in the following order:

- Command Line Argument
-- flag-1: -flag-1 value1 | -flag-1=value1 | --flag-1 value1 | --flag-1=value1
-- flag.2: -flag.2 value2 | -flag.2=value2 | --flag.2 value2 | --flag.2=value2
- Environment Variable
-- flag-1: FLAG\_1=value1
-- flag.2: FLAG\_2=value2
- Configuration File (my-application.properties)
-- flag-1: flag-1: value1
-- flag.2: flag.2: value2
- Default value provided in code

