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

// Tab is a group of measures to be displayed on a screen
type Tab struct {
	measures []Measure
}

// PrintAll writes the full tablature to stdout
func (t Tab) PrintAll() {
	sections := NewTabSection(t.measures[0].chords[0].instrument.strings)
	for _, m := range t.measures {
		if len(m.chords) == 0 {
			continue
		}
		sections.Add(m.chords)
	}
	t.PrintSection(sections)
}

// PrintSection writes a single TabSection to stdout
func (t Tab) PrintSection(s TabSection) {
	for _, str := range s.strings {
		var line string
		for _, char := range s.sequences[str] {
			line += char
		}
		fmt.Println(line)
	}
}

// TabSection is a running group of Measures to be printed together.
// Makes for easier book keeping of a Tab being printed to stdout.
type TabSection struct {
	strings   []InstrumentString
	sequences map[InstrumentString][]string
}

// NewTabSection returns an empty TabSection to start a new line of a Tab.
func NewTabSection(strings []InstrumentString) TabSection {
	sequences := make(map[InstrumentString][]string)
	return TabSection{strings, sequences}
}

// Add appends the tablature for a list of chords to the section.
func (t TabSection) Add(chords []Chord) {

	// Add a blank chord between each chord in the measure
	blank := BlankChord(chords[0].instrument)
	columns := []Chord{blank}
	for _, c := range chords {
		columns = append(columns, c)
		columns = append(columns, blank)
	}

	for _, chord := range columns {
		for _, str := range t.strings {
			// Convert Chord position into display value
			var value string
			fret, ok := chord.positions[str]
			if ok {
				value = fmt.Sprintf("%d", fret)
			} else {
				value = "-"
			}

			t.sequences[str] = append(t.sequences[str], value)
		}
	}
}

// main defines a Tab and prints it to the terminal
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

	// Declare final Tab object to write to the screen
	tab := Tab{[]Measure{
		measure, measure, measure,
	}}

	tab.PrintAll()
}
