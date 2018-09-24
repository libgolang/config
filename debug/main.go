package main

import (
	"fmt"
	"github.com/libgolang/config"
)

func main() {
	config.SetName("myapp")
	config.DefineString("flag1", "default-flag1", "Usage")
	flag1 := config.GetString("flag1")

	/*
		config.AppName = "debug"
		strFlag := config.String("flag1", "default!", "Flag 1 usage")
		boolFlag := config.Bool("flag2", false, "Flag 2 usage")
		config.Parse()
	*/

	fmt.Printf("flag1: %s\n", flag1)
}
