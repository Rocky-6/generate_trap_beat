package instruments

import (
	"bytes"
	"math/rand"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func MkMelody(key string) ([]byte, error) {
	clock := smf.MetricTicks(96)
	s := smf.New()
	s.TimeFormat = clock
	tr := smf.Track{}
	tr.Add(0, smf.MetaMeter(4, 4))
	tr.Add(0, smf.MetaTempo(140))

	key_note := keyNoteMelody(key)
	note := uint8(0)
	tick := uint32(0)
	n := byte(69)
	length := 64

	for {
	L:
		for {
			switch n {
			case 60, 62, 64, 67, 69, 72, 74, 76, 81, 84:
				if note != n {
					note = n
					break L
				}
				n = n + byte(rand.Intn(5)) - 2
			default:
				n = n + byte(rand.Intn(5)) - 2
			}
		}

		if length == 2 {
			tick = clock.Ticks8th()
			tr.Add(0, midi.NoteOn(0, note+key_note, 100))
			tr.Add(tick, midi.NoteOff(0, note+key_note))
			break
		}

		switch rand.Intn(10) {
		case 0, 1, 2, 3:
			tick = clock.Ticks8th()
			length -= 2
		case 4, 5, 6, 7, 8, 9:
			tick = clock.Ticks4th()
			length -= 4
		default:
		}

		tr.Add(0, midi.NoteOn(0, note+key_note, 100))
		tr.Add(tick, midi.NoteOff(0, note+key_note))
		n = n + byte(rand.Intn(5)) - 2
		if length == 0 {
			break
		}
	}
	tr.Close(0)
	s.Add(tr)

	buf := new(bytes.Buffer)
	_, err := s.WriteTo(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func keyNoteMelody(key string) uint8 {
	noteNames := [12]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

	var note uint8

	for i, noteName := range noteNames {
		if key == noteName {
			note = uint8(i)
			break
		}
	}

	return note
}
