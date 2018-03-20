// cli defines the command-line interface for the asciitab app
package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	helptext string
)

func init() {
	helptext = `
usage: asciitab random [--key="key"]

Create ASCII tabs for guitar for use in practice or inspiration.
	`
}

func main() {
	// all tabs are for guitar currently
	guitar := NewGuitar()

	// parse "--abc=123" flags from command line
	flags := make(map[string]string)
	for _, arg := range os.Args {
		// all accepted flags currently start with "--"
		if arg[:2] != "--" {
			continue
		}
		// See if flag is form of "--flag" or "--flag=value"
		equalsIndex := strings.Index(arg, "=")
		if equalsIndex == -1 {
			flags[arg[2:]] = ""
			continue
		}
		// Flag is form of "--flag=value", set it and trim quotes
		value := strings.Trim(arg[equalsIndex+1:], `'"`)
		flags[arg[2:equalsIndex]] = value
	}

	// asciitab random --key=A#
	if len(os.Args) >= 2 && os.Args[1] == "random" {
		// key provided or random key?
		key, ok := flags["key"]
		if !ok {
			key = RandomNote()
		}
		// scale provided or random scale?
		scaleName, ok := flags["scale"]
		if !ok {
			scaleName = RandomScaleName()
		}
		scale, err := NewScale(scaleName, key)
		if err != nil {
			panic(err)
		}
		tab := RandomTab(guitar, scale)

		fmt.Printf("Random tab in %s\n", scale)
		tab.PrintAll()
		return
	}

	// default
	fmt.Println(strings.TrimSpace(helptext))
}
