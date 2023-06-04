package instruments_test

import (
	"os"
	"testing"

	"github.com/Rocky-6/trap/database"
	"github.com/Rocky-6/trap/instruments"
)

func TestMkKick(t *testing.T) {
	midiData, err := instruments.MkKick()
	if err != nil {
		t.Errorf("MkKick returned an error: %v", err)
	}

	f, err := os.Create("kick.mid")
	if err != nil {
		t.Errorf("cannot create kick.mid")
	}
	defer f.Close()

	_, err = f.Write(midiData)
	if err != nil {
		t.Errorf("cannot write kick.mid")
	}
}

func TestMkHihat(t *testing.T) {
	midiData, err := instruments.MkHihat()
	if err != nil {
		t.Errorf("MkHihat returned an error: %v", err)
	}

	f, err := os.Create("hihat.mid")
	if err != nil {
		t.Errorf("cannot create hihat.mid")
	}
	defer f.Close()

	_, err = f.Write(midiData)
	if err != nil {
		t.Errorf("cannot write hihat.mid")
	}
}

func TestMkClap(t *testing.T) {
	midiData, err := instruments.MkClap()
	if err != nil {
		t.Errorf("MkClap returned an error: %v", err)
	}

	f, err := os.Create("clap.mid")
	if err != nil {
		t.Errorf("cannot create clap.mid")
	}
	defer f.Close()

	_, err = f.Write(midiData)
	if err != nil {
		t.Errorf("cannot write clap.mid")
	}
}

func TestMkMelody(t *testing.T) {
	midiData, err := instruments.MkMelody("C")
	if err != nil {
		t.Errorf("MkMelody returned an error: %v", err)
	}

	f, err := os.Create("melody.mid")
	if err != nil {
		t.Errorf("cannot create melody.mid")
	}
	defer f.Close()

	_, err = f.Write(midiData)
	if err != nil {
		t.Errorf("cannot write melody.mid")
	}
}

func TestMkChord(t *testing.T) {
	cp := database.DispChord()
	midiData, err := instruments.MkChord("C", cp)
	if err != nil {
		t.Errorf("MkChord returned an error: %v", err)
	}

	f, err := os.Create("chord.mid")
	if err != nil {
		t.Errorf("cannot create chord.mid")
	}
	defer f.Close()

	_, err = f.Write(midiData)
	if err != nil {
		t.Errorf("cannot write chord.mid")
	}
}

func TestMkBass(t *testing.T) {
	cp := database.DispChord()
	midiData, err := instruments.MkBass("C", cp)
	if err != nil {
		t.Errorf("MkBass returned an error: %v", err)
	}

	f, err := os.Create("bass.mid")
	if err != nil {
		t.Errorf("cannot create bass.mid")
	}
	defer f.Close()

	_, err = f.Write(midiData)
	if err != nil {
		t.Errorf("cannot write bass.mid")
	}
}
