// AudioPlayer represents the current audio state.
package main

import (
	"bytes"
	_ "embed"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

var (
	sampleRate   = 48000
	audioContext *audio.Context
)

//go:embed fail.ogg
var failOgg []byte

//go:embed triple.ogg
var tripleOgg []byte

type AudioPlayer struct {
	audioPlayer *audio.Player
}

func init() {
	audioContext = audio.NewContext(sampleRate)
}

func PlaySound(ogg []byte) (*AudioPlayer, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}
	var s audioStream
	var err error
	s, err = vorbis.DecodeWithoutResampling(bytes.NewReader(ogg))
	if err != nil {
		return nil, err
	}
	p, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}
	player := &AudioPlayer{
		audioPlayer: p,
	}

	player.audioPlayer.Play()

	return player, nil
}

func (p *AudioPlayer) Close() error {
	return p.audioPlayer.Close()
}

func (p *AudioPlayer) update() error {
	return nil
}
