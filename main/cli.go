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

Create ASCII tabs for guitar and ukulele for use in practice or inspiration.
	`
}

// parseFlags returns a map of "--key=value" args from the input.
func parseFlags() map[string]string {
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
	return flags
}

func main() {
	// parse "--abc=123" flags from command line
	flags := parseFlags()

	// key provided or random key?
	key, ok := flags["key"]
	if !ok {
		key = RandomNote()
	}

	// scale name provided or random scale?
	scaleName, ok := flags["scale"]
	if !ok {
		scaleName = RandomScaleName()
	}

	scale, err := NewScale(scaleName, key)
	if err != nil {
		panic(err)
	}

	// guitar or ukulele?
	var instrument Instrument
	instName, ok := flags["instrument"]
	if ok && instName == "ukulele" {
		instrument = NewUkulele()
	} else {
		instrument = NewGuitar()
	}

	// asciitab random --key=A#
	if len(os.Args) >= 2 && os.Args[1] == "random" {
		tab := RandomTab(instrument, scale)
		fmt.Printf("Random tab in %s\n", scale)
		tab.PrintAll()
		return
	}

	// asciitab scale --key=B
	if len(os.Args) >= 2 && os.Args[1] == "scale" {
		tab := ScaleTab(instrument, scale)
		fmt.Printf("Basic scale in %s\n", scale)
		tab.PrintAll()
		return
	}

	// default
	fmt.Println(strings.TrimSpace(helptext))
}
