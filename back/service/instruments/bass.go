package service

import (
	"bytes"
	"context"

	"github.com/Rocky-6/trap/model"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

type bass struct {
	key             string
	chordInfomation []model.ChordInfomation
}

func (bass *bass) MakeSMF(ctx context.Context) ([]byte, error) {
	clock := smf.MetricTicks(96)
	s := smf.New()
	s.TimeFormat = clock
	tr := smf.Track{}
	tr.Add(0, smf.MetaMeter(4, 4))
	tr.Add(0, smf.MetaTempo(140))

	c := bassNote(keyNoteBass(bass.key), bass.chordInfomation[0].DegreeName)
	tr.Add(0, midi.NoteOn(0, c, 100))
	tr.Add(clock.Ticks64th(), midi.NoteOff(0, c))

	c = bassNote(keyNoteBass(bass.key), bass.chordInfomation[2].DegreeName)
	tr.Add(clock.Ticks4th()*5-clock.Ticks64th(), midi.NoteOn(0, c, 100))
	tr.Add(clock.Ticks64th(), midi.NoteOff(0, c))

	c = bassNote(keyNoteBass(bass.key), bass.chordInfomation[0].DegreeName)
	tr.Add(clock.Ticks4th()*3-clock.Ticks64th(), midi.NoteOn(0, c, 100))
	tr.Add(clock.Ticks64th(), midi.NoteOff(0, c))

	c = bassNote(keyNoteBass(bass.key), bass.chordInfomation[1].DegreeName)
	tr.Add(clock.Ticks4th()*3-clock.Ticks64th(), midi.NoteOn(0, c, 100))
	tr.Add(clock.Ticks64th(), midi.NoteOff(0, c))

	c = bassNote(keyNoteBass(bass.key), bass.chordInfomation[2].DegreeName)
	tr.Add(clock.Ticks4th()*2+clock.Ticks8th()-clock.Ticks64th(), midi.NoteOn(0, c, 100))
	tr.Add(clock.Ticks64th(), midi.NoteOff(0, c))

	tr.Close(0)
	s.Add(tr)

	buf := new(bytes.Buffer)
	if _, err := s.WriteTo(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func bassNote(keyNoteChord uint8, degreeName string) uint8 {
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

	return root
}

func keyNoteBass(key string) uint8 {
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
