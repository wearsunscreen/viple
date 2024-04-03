// Player represents the current audio state.
package main

import (
	"bytes"
	_ "embed"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

//go:embed fail.ogg
var failOgg []byte

//go:embed triple.ogg
var tripleOgg []byte

type Player struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	seBytes      []byte
	seCh         chan []byte
}

func NewPlayer(audioContext *audio.Context, ogg []byte) (*Player, error) {
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
	player := &Player{
		audioContext: audioContext,
		audioPlayer:  p,
		seCh:         make(chan []byte),
	}

	player.audioPlayer.Play()

	return player, nil
}

func (p *Player) Close() error {
	return p.audioPlayer.Close()
}

func (p *Player) update() error {
	select {
	case p.seBytes = <-p.seCh:
		close(p.seCh)
		p.seCh = nil
	default:
	}
	return nil
}
