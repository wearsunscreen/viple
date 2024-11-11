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

//go:embed resources/fail.ogg
var failOgg []byte

//go:embed resources/win.ogg
var winOgg []byte

//go:embed resources/triple.ogg
var tripleOgg []byte

//go:embed resources/brick.ogg
var brickOgg []byte

//go:embed resources/paddle.ogg
var paddleOgg []byte

//go:embed resources/music.ogg
var musicOgg []byte

type AudioPlayer struct {
	audioPlayer *audio.Player
}

func init() {
	audioContext = audio.NewContext(sampleRate)
}

func PlaySound(ogg []byte) error {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}
	var s audioStream
	var err error
	s, err = vorbis.DecodeWithoutResampling(bytes.NewReader(ogg))
	if err != nil {
		return err
	}
	p, err := audioContext.NewPlayer(s)
	if err != nil {
		return err
	}
	player := &AudioPlayer{
		audioPlayer: p,
	}

	player.audioPlayer.Play()

	return nil
}

func PlaySoundForever(ogg []byte) error {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}
	var s audioStream
	var err error
	s, err = vorbis.DecodeWithoutResampling(bytes.NewReader(ogg))
	if err != nil {
		return err
	}
	loop := audio.NewInfiniteLoop(s, s.Length())
	p, err := audioContext.NewPlayer(loop)
	if err != nil {
		return err
	}
	player := &AudioPlayer{
		audioPlayer: p,
	}

	player.audioPlayer.Play()

	return nil
}

func (p *AudioPlayer) Close() error {
	return p.audioPlayer.Close()
}
