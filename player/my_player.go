package player

import (
	"fmt"
	"go-playlist-player/types"
	"go-playlist-player/types/commands"
	"go-playlist-player/types/state"
	"sync/atomic"
	"time"
)

type MyPlayer struct {
	playlist         types.SyncPlaylist
	currentSongIndex atomic.Int64
	CurrentPlayable  *types.PlayableImpl
	WantedState      *state.SafePlayState
	CurrentState     *state.SafePlayState
	Commands         *chan commands.PlayCommand
}

func (player *MyPlayer) startPlayer() {
	for player.WantedState.GetState() == state.Paused {
		player.CurrentPlayable.PauseAsync()
		player.CurrentState.SetState(state.Paused)
		time.Sleep(time.Millisecond)
	}

	if player.WantedState.GetState() == state.Playing {
		player.CurrentPlayable.PlayAsync()
		player.CurrentState.SetState(state.Playing)
	}

	for player.CurrentPlayable.IsPlayingAsync() && player.WantedState.GetState() != state.Closed {
		for player.WantedState.GetState() == state.Paused {
			player.CurrentPlayable.PauseAsync()
			player.CurrentState.SetState(state.Paused)
			time.Sleep(time.Millisecond)
		}

		if player.WantedState.GetState() == state.Playing {
			player.CurrentPlayable.PlayAsync()
			player.CurrentState.SetState(state.Playing)
			time.Sleep(time.Millisecond)
		}
	}

	if player.WantedState.GetState() != state.Closed {
		go player.Next()
		for player.WantedState.GetState() != state.Closed {
			time.Sleep(time.Millisecond)
		}
	}

	var err = player.CurrentPlayable.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
	player.CurrentState.SetState(state.Closed)
}

func (player *MyPlayer) Play() {
	fmt.Println("start playing")
	player.WantedState.SetState(state.Playing)

	for player.CurrentState.GetState() != state.Playing {
		time.Sleep(time.Millisecond)
	}
}

func (player *MyPlayer) Pause() {
	fmt.Println("pause playing")
	player.WantedState.SetState(state.Paused)

	for player.CurrentState.GetState() != state.Paused {
		time.Sleep(time.Millisecond)
	}
}

func (player *MyPlayer) AddSongFromDisc(song *types.Song) {
	player.playlist.AddSong(song)

	if player.playlist.Length() == 1 {
		go player.initSingleSongPlayer(false)
	}
}

func (player *MyPlayer) Next() {
	fmt.Println("next song")

	if int(player.currentSongIndex.Load()+1) < player.playlist.Length() {
		player.currentSongIndex.Add(1)
	} else {
		player.currentSongIndex.Swap(0)
	}

	var shouldPlay = player.CurrentState.GetState() == state.Playing
	player.WantedState.SetState(state.Closed)
	player.initSingleSongPlayer(shouldPlay)

	for player.CurrentState.GetState() != state.Paused && player.CurrentState.GetState() != state.Playing {
		time.Sleep(time.Millisecond)
	}
}

func (player *MyPlayer) Prev() {
	fmt.Println("prev song")

	if int(player.currentSongIndex.Load()) > 0 {
		player.currentSongIndex.Add(-1)
	} else {
		player.currentSongIndex.Swap(int64(player.playlist.Length() - 1))
	}

	var shouldPlay = player.CurrentState.GetState() == state.Playing
	player.WantedState.SetState(state.Closed)
	player.initSingleSongPlayer(shouldPlay)
	for player.CurrentState.GetState() != state.Paused && player.CurrentState.GetState() != state.Playing {
		time.Sleep(time.Millisecond)
	}
}

func (player *MyPlayer) initSingleSongPlayer(play bool) {
	song := player.playlist.GetSong(int(player.currentSongIndex.Load()))

	for player.CurrentState.GetState() != state.Closed {
		time.Sleep(time.Millisecond)
	}
	player.CurrentPlayable.SetPlayable(createPlayer(song.Path).GetPlayable())

	if play {
		player.WantedState.SetState(state.Playing)
	} else {
		player.WantedState.SetState(state.Paused)
	}

	go player.startPlayer()
}
