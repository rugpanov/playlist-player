package state

import "sync"

type SafePlayState struct {
	state PlayState
	mu    sync.RWMutex
}

func (state *SafePlayState) GetState() PlayState {
	state.mu.RLock()
	defer state.mu.RUnlock()

	return state.state
}

func (state *SafePlayState) SetState(ps PlayState) {
	if state.GetState() == ps {
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	state.state = ps
}

type PlayState int

const (
	Playing PlayState = iota
	Paused
	Closed
)
