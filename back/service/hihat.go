package service

import (
	"bytes"
	"context"

	"github.com/Rocky-6/trap/repository"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

type hihat struct{}

func NewHihat() repository.InstrumentsRepository {
	return &hihat{}
}

func (hihat *hihat) MakeSMF(ctx context.Context) ([]byte, error) {
	clock := smf.MetricTicks(96)
	s := smf.New()
	s.TimeFormat = clock
	tr := smf.Track{}
	tr.Add(0, smf.MetaMeter(4, 4))
	tr.Add(0, smf.MetaTempo(140))

	tr.Add(0, midi.NoteOn(0, midi.C(5), 100))
	tr.Add(clock.Ticks64th(), midi.NoteOff(0, midi.C(5)))
	for i := 0; i < 31; i++ {
		tr.Add(clock.Ticks8th()-clock.Ticks64th(), midi.NoteOn(0, midi.C(5), 100))
		tr.Add(clock.Ticks64th(), midi.NoteOff(0, midi.C(5)))
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
