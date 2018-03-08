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

// GetString returns an instruments string by name (e.g. "A")
func (i Instrument) GetString(name string) (InstrumentString, error) {
	for _, s := range i.strings {
		if s.name == name {
			return s, nil
		}
	}
	var nullString InstrumentString
	return nullString, fmt.Errorf("string '%s' not found on instrument '%s'", name, i)
}

// Chord is a fixed position on the frets of an Instrument
type Chord struct {
	instrument Instrument
	positions  map[InstrumentString]int32
}

// NewChord returns a Chord from a map of strings to positions
func NewChord(i Instrument, namesToPositions map[string]int32) Chord {
	strPositions := make(map[InstrumentString]int32)
	for name, position := range namesToPositions {
		str, err := i.GetString(name)
		if err != nil {
			panic(err)
		}
		strPositions[str] = position
	}

	c := Chord{
		instrument: i,
		positions:  strPositions,
	}

	return c
}

func printChordAsTab(c Chord) {
	for _, s := range c.instrument.strings {
		var displayValue string
		fret, ok := c.positions[s]
		if ok {
			displayValue = fmt.Sprintf("%d", fret)
		} else {
			displayValue = "-"
		}
		fmt.Println(displayValue)
	}
}

func main() {
	guitar := NewInstrument([]string{
		"e", "B", "G", "D", "A", "E",
	})

	c := NewChord(guitar, map[string]int32{
		"B": 1,
		"D": 2,
		"A": 3,
	})

	fmt.Println(guitar.strings)

	printChordAsTab(c)
}
