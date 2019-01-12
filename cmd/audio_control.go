package cmd

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type AudioControl struct {
	Sounds map[string]Audio

	engine *Engine
}

func NewAudioControl() AudioControl {
	return AudioControl{
		Sounds: make(map[string]Audio),
	}
}

func (ac *AudioControl) Initialize(e *Engine) {
	ac.engine = e
}

func (ac *AudioControl) Load(path string, name string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	s, format, err := wav.Decode(f)
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})

	ac.Sounds[name] = Audio{
		S:      &s,
		Format: &format,
		Done:   done,
	}
}

func (ac *AudioControl) Play(name string) {
	audio := ac.Sounds[name]
	speaker.Init(audio.Format.SampleRate, audio.Format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(*audio.S, beep.Callback(func() {
		close(audio.Done)
	})))
}

type Audio struct {
	S      *beep.StreamSeekCloser
	Format *beep.Format
	Done   chan struct{}
}
