package instruments

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func MkChord(key string, cp [4]string) ([]byte, error) {
	clock := smf.MetricTicks(96)
	s := smf.New()
	s.TimeFormat = clock
	tr := smf.Track{}
	tr.Add(0, smf.MetaMeter(4, 4))
	tr.Add(0, smf.MetaTempo(140))

	// start
	for j := 0; j < 2; j++ {
		for i := 0; i < 4; i++ {
			c := chordNote(keyNoteChord(key), cp[i])

			for _, v := range c {
				tr.Add(0, midi.NoteOn(0, v, 100))
			}

			for j, v := range c {
				if j == 0 {
					tr.Add(clock.Ticks4th()*2, midi.NoteOff(0, v))
				} else {
					tr.Add(0, midi.NoteOff(0, v))
				}
			}
		}
	}
	// end

	tr.Close(0)
	s.Add(tr)

	buf := new(bytes.Buffer)
	_, err := s.WriteTo(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func chordNote(keyNoteChord uint8, degree_name string) []uint8 {
	chord_note := make([]uint8, 3)
	root := keyNoteChord

	switch true {
	case check_regexp(`bVII`, degree_name):
		root += 10
	case check_regexp(`VII`, degree_name):
		root += 11
	case check_regexp(`bVI`, degree_name):
		root += 8
	case check_regexp(`VI`, degree_name):
		root += 9
	case check_regexp(`#IV`, degree_name):
		root += 6
	case check_regexp(`IV`, degree_name):
		root += 5
	case check_regexp(`V`, degree_name):
		root += 7
	case check_regexp(`bIII`, degree_name):
		root += 3
	case check_regexp(`III`, degree_name):
		root += 4
	case check_regexp(`bII`, degree_name):
		root += 1
	case check_regexp(`II`, degree_name):
		root += 2
	default:
	}

	chord_note[0] = root

	chord_note[1] = root + 4
	if check_regexp(`m`, degree_name) {
		chord_note[1] = root + 3
	}
	if check_regexp(`sus4`, degree_name) {
		chord_note[1] = root + 5
	}

	chord_note[2] = root + 7
	if check_regexp(`b5`, degree_name) || check_regexp(`dim`, degree_name) {
		chord_note[2] = root + 6
	}

	if check_regexp(`7`, degree_name) && !check_regexp(`M7`, degree_name) {
		chord_note = append(chord_note, root+10)
	}
	if check_regexp(`M7`, degree_name) {
		chord_note = append(chord_note, root+11)
	}
	if check_regexp(`6`, degree_name) && !check_regexp(`\(6`, degree_name) {
		chord_note = append(chord_note, root+9)
	}

	if check_regexp(`\(6`, degree_name) || check_regexp(`13`, degree_name) {
		chord_note = append(chord_note, root+21)
	}
	if check_regexp(`9`, degree_name) {
		chord_note = append(chord_note, root+14)
	}
	if check_regexp(`11`, degree_name) && !check_regexp(`#11`, degree_name) {
		chord_note = append(chord_note, root+17)
	}
	if check_regexp(`#11`, degree_name) {
		chord_note = append(chord_note, root+18)
	}

	sort.Slice(chord_note, func(i, j int) bool {
		return chord_note[i] < chord_note[j]
	})

	return chord_note[:]
}

func check_regexp(reg, str string) bool {
	r := regexp.MustCompile(reg)
	return r.MatchString(str)
}

func keyNoteChord(key string) uint8 {
	noteNames := [12]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

	var note uint8

	for i, noteName := range noteNames {
		if key == noteName {
			note = uint8(i) + 48
			break
		}
	}

	if note == 0 {
		fmt.Println("入力された音階名が不正です")
		return 0
	}

	return note
}
