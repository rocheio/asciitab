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

// BlankChord returns a Chord with no positions, representing a pause
func BlankChord(i Instrument) Chord {
	return NewChord(i, map[string]int32{})
}

// Measure is a series of Chords that should be displayed together
type Measure struct {
	chords []Chord
}

// printMeasure writes a multi-column measure to stdout
func printMeasure(m Measure) {
	if len(m.chords) == 0 {
		return
	}

	// Add a blank chord between each chord in the measure
	blank := BlankChord(m.chords[0].instrument)
	chords := []Chord{blank}
	for _, c := range m.chords {
		chords = append(chords, c)
		chords = append(chords, blank)
	}

	// Transform list of Chords [E, Am, C#] into a map of
	// each string to its fret position sequence
	fretSequences := make(map[InstrumentString][]string)
	strings := m.chords[0].instrument.strings

	for _, chord := range chords {
		for _, str := range strings {
			// Convert Chord position into display value
			var value string
			fret, ok := chord.positions[str]
			if ok {
				value = fmt.Sprintf("%d", fret)
			} else {
				value = "-"
			}

			fretSequences[str] = append(fretSequences[str], value)
		}
	}

	for _, str := range strings {
		var line string
		for _, char := range fretSequences[str] {
			line += char
		}
		fmt.Println(line)
	}
}

func main() {
	// Declare instrument and strings
	guitar := NewInstrument([]string{
		"e", "B", "G", "D", "A", "E",
	})

	// Declare chords used in the tab
	c := NewChord(guitar, map[string]int32{
		"B": 1,
		"D": 2,
		"A": 3,
	})

	// Declare blocks of chords to display together
	measure := Measure{[]Chord{
		c, c, c,
	}}

	printMeasure(measure)
}
