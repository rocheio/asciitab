package main

import "fmt"

// InstrumentString is a physical string on an Instrument.
type InstrumentString struct {
	name string
}

// Instrument has strings that are played to make music.
type Instrument struct {
	strings []InstrumentString
}

// NewInstrument returns a new instrument with InstrumentStrings
// named from a list of character strings.
func NewInstrument(stringNames []string) Instrument {
	var strings []InstrumentString
	for _, str := range stringNames {
		strings = append(strings, InstrumentString{str})
	}
	return Instrument{
		strings: strings,
	}
}

func main() {
	guitar := NewInstrument([]string{
		"E", "A", "D", "G", "B", "e",
	})
	fmt.Println(guitar.strings)
}
