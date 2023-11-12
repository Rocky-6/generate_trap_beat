package service

import (
	"bytes"
	"context"
	"sort"

	"github.com/Rocky-6/trap/model"
	"github.com/Rocky-6/trap/repository"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

type chord struct {
	key             string
	chordInfomation []model.ChordInfomation
}

func NewChord(key string, chordInformation []model.ChordInfomation) repository.InstrumentsRepository {
	return &chord{
		key:             key,
		chordInfomation: chordInformation,
	}
}

func (chord *chord) MakeSMF(ctx context.Context) ([]byte, error) {
	clock := smf.MetricTicks(96)
	s := smf.New()
	s.TimeFormat = clock
	tr := smf.Track{}
	tr.Add(0, smf.MetaMeter(4, 4))
	tr.Add(0, smf.MetaTempo(140))

	// start
	for j := 0; j < 2; j++ {
		for i := 0; i < 4; i++ {
			c := chordNote(keyNoteChord(chord.key), chord.chordInfomation[i].DegreeName)

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

func chordNote(keyNoteChord uint8, degreeName string) []uint8 {
	chord_note := make([]uint8, 3)
	root := keyNoteChord

	switch true {
	case check_regexp(`bVII`, degreeName):
		root += 10
	case check_regexp(`VII`, degreeName):
		root += 11
	case check_regexp(`bVI`, degreeName):
		root += 8
	case check_regexp(`VI`, degreeName):
		root += 9
	case check_regexp(`#IV`, degreeName):
		root += 6
	case check_regexp(`IV`, degreeName):
		root += 5
	case check_regexp(`V`, degreeName):
		root += 7
	case check_regexp(`bIII`, degreeName):
		root += 3
	case check_regexp(`III`, degreeName):
		root += 4
	case check_regexp(`bII`, degreeName):
		root += 1
	case check_regexp(`II`, degreeName):
		root += 2
	default:
	}

	chord_note[0] = root

	chord_note[1] = root + 4
	if check_regexp(`m`, degreeName) {
		chord_note[1] = root + 3
	}
	if check_regexp(`sus4`, degreeName) {
		chord_note[1] = root + 5
	}

	chord_note[2] = root + 7
	if check_regexp(`b5`, degreeName) || check_regexp(`dim`, degreeName) {
		chord_note[2] = root + 6
	}

	if check_regexp(`7`, degreeName) && !check_regexp(`M7`, degreeName) {
		chord_note = append(chord_note, root+10)
	}
	if check_regexp(`M7`, degreeName) {
		chord_note = append(chord_note, root+11)
	}
	if check_regexp(`6`, degreeName) && !check_regexp(`\(6`, degreeName) {
		chord_note = append(chord_note, root+9)
	}

	if check_regexp(`\(6`, degreeName) || check_regexp(`13`, degreeName) {
		chord_note = append(chord_note, root+21)
	}
	if check_regexp(`9`, degreeName) {
		chord_note = append(chord_note, root+14)
	}
	if check_regexp(`11`, degreeName) && !check_regexp(`#11`, degreeName) {
		chord_note = append(chord_note, root+17)
	}
	if check_regexp(`#11`, degreeName) {
		chord_note = append(chord_note, root+18)
	}

	sort.Slice(chord_note, func(i, j int) bool {
		return chord_note[i] < chord_note[j]
	})

	return chord_note[:]
}

func keyNoteChord(key string) uint8 {
	noteNames := [12]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

	var note uint8

	for i, noteName := range noteNames {
		if key == noteName {
			note = uint8(i) + midi.C(5)
			break
		}
	}

	return note
}
