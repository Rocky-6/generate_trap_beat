package service

import (
	"bytes"
	"context"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

type clap struct{}

func (clap *clap) MakeSMF(ctx context.Context) ([]byte, error) {
	clock := smf.MetricTicks(96)
	s := smf.New()
	s.TimeFormat = clock
	tr := smf.Track{}
	tr.Add(0, smf.MetaMeter(4, 4))
	tr.Add(0, smf.MetaTempo(140))

	tr.Add(clock.Ticks4th()*2, midi.NoteOn(0, midi.C(5), 100))
	tr.Add(clock.Ticks64th(), midi.NoteOff(0, midi.C(5)))
	tr.Add(clock.Ticks4th()*4-clock.Ticks64th(), midi.NoteOn(0, midi.C(5), 100))
	tr.Add(clock.Ticks64th(), midi.NoteOff(0, midi.C(5)))
	tr.Add(clock.Ticks4th()*4-clock.Ticks64th(), midi.NoteOn(0, midi.C(5), 100))
	tr.Add(clock.Ticks64th(), midi.NoteOff(0, midi.C(5)))
	tr.Add(clock.Ticks4th()*4-clock.Ticks64th(), midi.NoteOn(0, midi.C(5), 100))
	tr.Add(clock.Ticks64th(), midi.NoteOff(0, midi.C(5)))

	tr.Close(0)
	s.Add(tr)

	buf := new(bytes.Buffer)
	_, err := s.WriteTo(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
