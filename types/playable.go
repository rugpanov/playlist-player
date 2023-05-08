package types

import (
	"io"
	"sync"
)

type Playable interface {
	Pause()
	Play()
	IsPlaying() bool
	io.Closer
}

type PlayableImpl struct {
	mu       sync.RWMutex
	playable Playable
}

func (p *PlayableImpl) GetPlayable() Playable {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.playable
}

func (p *PlayableImpl) SetPlayable(playable Playable) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.playable = playable
}

func (p *PlayableImpl) IsPlayingAsync() bool {
	return p.GetPlayable().IsPlaying()
}

func (p *PlayableImpl) PlayAsync() {
	p.GetPlayable().Play()
}

func (p *PlayableImpl) PauseAsync() {
	p.GetPlayable().Pause()
}

func (p *PlayableImpl) Close() error {
	var err = p.GetPlayable().Close()
	return err
}
