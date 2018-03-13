package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var pitchProgression map[string]string

func init() {
	rand.Seed(time.Now().Unix())
	pitchProgression = map[string]string{
		"A":  "A#",
		"A#": "B",
		"B":  "C",
		"C":  "C#",
		"C#": "D",
		"D":  "D#",
		"D#": "E",
		"E":  "F",
		"F":  "F#",
		"F#": "G",
		"G":  "G#",
		"G#": "A",
	}
}

// InstrumentString is a physical string on an Instrument.
type InstrumentString struct {
	name string
}

// PitchMatches returns True if this string is in a list of pitches
func PitchMatches(possible string, pitches []string) bool {
	for _, p := range pitches {
		if strings.ToUpper(possible) == strings.ToUpper(p) {
			return true
		}
	}
	return false
}

// IncreasePitch returns one pitch higher than passed in
func IncreasePitch(pitch string) string {
	p, ok := pitchProgression[strings.ToUpper(pitch)]
	if !ok {
		return pitch
	}
	return p
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

// NewGuitar returns an Instrument configured as a standard guitar
func NewGuitar() Instrument {
	return NewInstrument([]string{
		"e", "B", "G", "D", "A", "E",
	})
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

// AddChord adds a Chord to the current Measure
func (m Measure) AddChord(c Chord) {
	m.chords = append(m.chords, c)
}

// Tab is a group of measures to be displayed on a screen
type Tab struct {
	measures *[]Measure
}

// NewTab returns a new, empty Tab for building on
func NewTab() Tab {
	m := []Measure{}
	return Tab{&m}
}

// AddMeasure adds a Measure to the current Tab
func (t Tab) AddMeasure(m Measure) {
	*t.measures = append(*t.measures, m)
}

// PrintAll writes the full tablature to stdout
func (t Tab) PrintAll() {
	measures := *t.measures
	section := NewTabSection(measures[0].chords[0].instrument.strings)
	section.AddBarLine()

	for _, m := range *t.measures {
		if len(m.chords) == 0 {
			continue
		}
		section.AddChords(m.chords)
		section.AddBarLine()
	}

	t.PrintSection(section)
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
// Makes for easier bookkeeping of a Tab being printed to stdout.
type TabSection struct {
	strings   []InstrumentString
	sequences map[InstrumentString][]string
}

// NewTabSection returns an empty TabSection to start a new line of a Tab.
func NewTabSection(strings []InstrumentString) TabSection {
	sequences := make(map[InstrumentString][]string)
	return TabSection{strings, sequences}
}

// AddChords converts chords to tablature and appends to this section.
func (t TabSection) AddChords(chords []Chord) {

	// Add a blank chord between each chord in the measure
	blank := BlankChord(chords[0].instrument)
	columns := []Chord{blank, blank}
	for _, c := range chords {
		columns = append(columns, c)
		columns = append(columns, blank)
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

// AddBarLine appends a vertical line to this section to separate Measures.
func (t TabSection) AddBarLine() {
	for _, str := range t.strings {
		t.sequences[str] = append(t.sequences[str], "|")
	}
}

// Scale is a group of Pitches played together in harmony
// (e.g. Classical: A - G w/ flats and sharps)
type Scale struct {
	pitches []string
}

// randomChord returns a Chord with random values (0-4)
func randomChord(inst Instrument) Chord {
	positions := make(map[string]int32)
	for _, s := range inst.strings {
		// Positions on only ~half the strings
		if rand.Float64() < 0.5 {
			continue
		}
		positions[s.name] = rand.Int31n(4)
	}
	return NewChord(inst, positions)
}

// randomChordInScale returns a Chord in the given Scale
func randomChordInScale(inst Instrument, scale Scale) Chord {
	positions := make(map[string]int32)
	for _, str := range inst.strings {
		// find possible positions
		var possibles []int
		pitch := str.name
		for i := 0; i <= 3; i++ {
			if PitchMatches(pitch, scale.pitches) {
				possibles = append(possibles, i)
			}
			pitch = IncreasePitch(pitch)
		}
		// half chance to use each string
		if rand.Float64() > 0.5 {
			// equal chance to use each possible fret in scale
			fret := possibles[rand.Intn(len(possibles))]
			positions[str.name] = int32(fret)
		}
	}
	return NewChord(inst, positions)
}

// main defines a Tab and prints it to the terminal
func main() {
	guitar := NewGuitar()
	aMinor := Scale{[]string{
		"A", "B", "C", "D", "E", "F", "G",
	}}
	tab := NewTab()
	for i := 0; i < 4; i++ {
		var chords []Chord
		for i := 0; i < 4; i++ {
			c := randomChordInScale(guitar, aMinor)
			chords = append(chords, c)
		}
		measure := Measure{chords}
		tab.AddMeasure(measure)
	}
	tab.PrintAll()
}
