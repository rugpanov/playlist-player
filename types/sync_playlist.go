package types

import (
	"sync"
)

type SyncPlaylist struct {
	playlist []*Song
	mu       sync.RWMutex
}

func (syncPlaylist *SyncPlaylist) Length() int {
	syncPlaylist.mu.RLock()
	defer syncPlaylist.mu.RUnlock()

	return len(syncPlaylist.playlist)
}

func (syncPlaylist *SyncPlaylist) GetSong(index int) *Song {
	syncPlaylist.mu.RLock()
	defer syncPlaylist.mu.RUnlock()

	return syncPlaylist.playlist[index]
}

func (syncPlaylist *SyncPlaylist) AddSong(song *Song) {
	syncPlaylist.mu.Lock()
	defer syncPlaylist.mu.Unlock()

	syncPlaylist.playlist = append(syncPlaylist.playlist, song)
}
