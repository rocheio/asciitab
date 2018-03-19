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

	// asciitab random
	if len(os.Args) == 2 && os.Args[1] == "random" {
		scale := RandomScale()
		tab := RandomTab(guitar, scale)

		fmt.Printf("Random tab in %s\n", scale)
		tab.PrintAll()
		return
	}

	// asciitab random --key=A#
	if len(os.Args) == 3 && os.Args[1] == "random" {
		if os.Args[2][:6] == "--key=" {
			key := os.Args[2][6:]
			key = strings.Trim(key, `'"`)
			scaleName := RandomScaleName()
			scale, err := NewScale(scaleName, key)
			if err != nil {
				panic(err)
			}
			tab := RandomTab(guitar, scale)

			fmt.Printf("Random tab in %s\n", scale)
			tab.PrintAll()
			return
		}
	}

	// default
	fmt.Println(strings.TrimSpace(helptext))
}
