// music defines the foundational structure of music used in tabs
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
	pitchProgression map[string]string
	scaleSteps       map[string][]int
)

func init() {
	rand.Seed(time.Now().Unix())
	pitchProgression = map[string]string{
		"Ab": "A",
		"A":  "A#",
		"A#": "B",
		"Bb": "B",
		"B":  "C",
		"C":  "C#",
		"C#": "D",
		"Db": "D",
		"D":  "D#",
		"D#": "E",
		"Eb": "E",
		"E":  "F",
		"F":  "F#",
		"F#": "G",
		"Gb": "G",
		"G":  "G#",
		"G#": "A",
	}
	scaleSteps = map[string][]int{
		"minor": []int{
			2, 1, 2, 2, 1, 2, 2,
		},
		"major": []int{
			2, 2, 1, 2, 2, 2, 1,
		},
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// InstrumentString is a physical string on an Instrument.
type InstrumentString struct {
	name string
}

// FretsInScale returns a list of frets in a given scale.
func (str InstrumentString) FretsInScale(scale Scale) []int {
	var frets []int

	fret := 0
	pitch := str.name
	for {
		if fret == 10 {
			break
		}
		if stringInSlice(pitch, scale.pitches) {
			frets = append(frets, fret)
		}
		pitch = pitchProgression[strings.ToUpper(pitch)]
		fret++
	}

	return frets
}

// indexOf returns the first fret where a pitch occurs on this string.
func (str InstrumentString) indexOf(targetPitch string) int {
	pitch := str.name
	for i := 0; i < 24; i++ {
		if pitch == targetPitch {
			return i
		}
		pitch = pitchProgression[strings.ToUpper(pitch)]
	}
	return -1
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
		"E", "A", "D", "G", "B", "e",
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
	positions  map[InstrumentString]int
}

// NewChord returns a Chord from a map of strings to positions
func NewChord(i Instrument, namesToPositions map[string]int) Chord {
	strPositions := make(map[InstrumentString]int)
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
	return NewChord(i, map[string]int{})
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
	section.AddStringNames()
	section.AddBar()

	for _, m := range *t.measures {
		if len(m.chords) == 0 {
			continue
		}
		section.AddChords(m.chords)
		section.AddBar()
	}

	t.PrintSection(section)
}

// PrintSection writes a single TabSection to stdout.
// Prints in reverse order of strings (top-to-bottom).
func (t Tab) PrintSection(s TabSection) {
	for i := len(s.strings) - 1; i >= 0; i-- {
		str := s.strings[i]
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

// AddBar appends a vertical line to this section to separate Measures.
func (t TabSection) AddBar() {
	for _, str := range t.strings {
		t.sequences[str] = append(t.sequences[str], "|")
	}
}

// AddStringNames appends the names of each string to start each section.
func (t TabSection) AddStringNames() {
	for _, str := range t.strings {
		t.sequences[str] = append(t.sequences[str], str.name)
	}
}

// Scale is a group of Pitches played together in harmony
// (e.g. Classical: A - G w/ flats and sharps)
type Scale struct {
	name    string
	root    string
	pitches []string
}

// String returns a string version of this Scale
func (s Scale) String() string {
	return fmt.Sprintf("<%s %s %s>", s.root, s.name, s.pitches)
}

// randomChord returns a Chord with random values (0-4)
func randomChord(inst Instrument) Chord {
	positions := make(map[string]int)
	for _, s := range inst.strings {
		// Positions on only ~half the strings
		if rand.Float64() < 0.5 {
			continue
		}
		positions[s.name] = rand.Intn(4)
	}
	return NewChord(inst, positions)
}

// randomChordInScale returns a Chord in the given Scale
func randomChordInScale(inst Instrument, scale Scale) Chord {
	positions := make(map[string]int)
	for _, str := range inst.strings {
		// find possible positions
		var possibles []int
		pitch := str.name
		for i := 0; i <= 4; i++ {
			if PitchMatches(pitch, scale.pitches) {
				possibles = append(possibles, i)
			}
			pitch = IncreasePitch(pitch)
		}
		// half chance to use each string
		if rand.Float64() > 0.7 {
			// equal chance to use each possible fret in scale
			fret := possibles[rand.Intn(len(possibles))]
			positions[str.name] = int(fret)
		}
	}
	return NewChord(inst, positions)
}

// singleNoteChord returns a 'chord' with just a single note
func singleNoteChord(inst Instrument, str string, fret int) Chord {
	positions := map[string]int{
		str: fret,
	}
	return NewChord(inst, positions)
}

// NewScale returns a Scale for a root note and minor or major step pattern.
func NewScale(name, root string) (Scale, error) {
	// start with the root note
	pitch := root
	pitches := []string{}
	steps, ok := scaleSteps[name]
	if !ok {
		return NilScale(), fmt.Errorf("scale not found: %s", name)
	}

	for _, step := range steps {
		pitches = append(pitches, pitch)
		// step up notes according to pattern
		for i := 0; i < step; i++ {
			oldpitch := pitch
			pitch, ok = pitchProgression[pitch]
			if !ok {
				return NilScale(), fmt.Errorf(
					"progression not found for pitch: %s", oldpitch,
				)
			}
		}
	}
	return Scale{name, root, pitches}, nil
}

// NilScale returns a Scale for returning with errors
func NilScale() Scale {
	return Scale{"", "", nil}
}

// RandomNote returns a random root note for scale progression
func RandomNote() string {
	keys := make([]string, len(pitchProgression))
	i := 0
	for k := range pitchProgression {
		keys[i] = k
		i++
	}
	return keys[rand.Intn(len(pitchProgression))]
}

// RandomScaleName returns the name of a random scale
func RandomScaleName() string {
	if rand.Float64() > 0.5 {
		return "major"
	}
	return "minor"
}

// RandomScale returns a major or minor scale with a random root note.
func RandomScale() Scale {
	name := RandomScaleName()
	note := RandomNote()
	scale, err := NewScale(name, note)
	if err != nil {
		panic(err)
	}
	return scale
}

// RandomTab returns a Tab with random notes from a Scale
func RandomTab(inst Instrument, scale Scale) Tab {
	// build a new tab using random chords in the scale
	tab := NewTab()
	for i := 0; i < 4; i++ {
		var chords []Chord
		for i := 0; i < 4; i++ {
			c := randomChordInScale(inst, scale)
			chords = append(chords, c)
		}
		measure := Measure{chords}
		tab.AddMeasure(measure)
	}

	return tab
}

// ScaleTab returns a Tab of all notes up and down the scale
func ScaleTab(inst Instrument, scale Scale) Tab {
	var scaleStrings []InstrumentString
	tab := NewTab()

	// Start on bottom string if <= 6 frets, second-to-bottom otherwise
	bassString := inst.strings[0]
	if bassString.indexOf(scale.root) == -1 || bassString.indexOf(scale.root) > 6 {
		scaleStrings = inst.strings[1:]
	} else {
		scaleStrings = inst.strings
	}

	// Find root note fret to start on
	rootFret := scaleStrings[0].indexOf(scale.root)

	// Add all frets on all strings in order
	for j := 0; j < len(scaleStrings); j++ {
		var chords []Chord

		str := scaleStrings[j]
		for _, fret := range str.FretsInScale(scale) {
			// Only use frets up to 4 after rootFret
			if fret < rootFret || fret > rootFret+4 {
				continue
			}
			c := singleNoteChord(inst, str.name, fret)
			chords = append(chords, c)
		}

		measure := Measure{chords}
		tab.AddMeasure(measure)
	}

	return tab
}
