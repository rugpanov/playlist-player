package main

import (
	"fmt"
	"go-playlist-player/player"
	"go-playlist-player/types"
	"go-playlist-player/types/commands"
	"go-playlist-player/types/state"
	"log"
	"net/http"
	"os"
	"time"
)

var commandQueue = make(chan commands.PlayCommand, 10000)

func handleRequests() {
	uploadPage := http.FileServer(http.Dir("./static"))

	http.Handle("/", uploadPage)
	http.HandleFunc("/player/play", handlePlay)
	http.HandleFunc("/player/pause", handlePause)
	http.HandleFunc("/player/next", handleNext)
	http.HandleFunc("/player/prev", handlePrev)
	http.HandleFunc("/playlist/upload", upload)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handlePlay(w http.ResponseWriter, _ *http.Request) {
	commandQueue <- commands.Play
	respondOK(w)
}

func handlePause(w http.ResponseWriter, _ *http.Request) {
	commandQueue <- commands.Pause
	respondOK(w)
}

func handleNext(w http.ResponseWriter, _ *http.Request) {
	commandQueue <- commands.Next
	respondOK(w)
}

func handlePrev(w http.ResponseWriter, _ *http.Request) {
	commandQueue <- commands.Prev
	respondOK(w)
}

func upload(w http.ResponseWriter, r *http.Request) {
	path, err := UploadFile(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("500 - " + err.Error()))
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}
	myPlayer.AddSongFromDisc(&types.Song{
		Path: path,
	})
	respondOK(w)
}

func respondOK(w http.ResponseWriter) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}

var myPlayer *player.MyPlayer = nil

func main() {
	go handleRequests()
	playerImpl := initializePlayer()
	go dispatchCommands(playerImpl)

	addAllFilesFromDir("music")

	fmt.Println("Player started. Visit http://localhost:8080")
	// we should never stop music
	for true {
		time.Sleep(10000)
	}
}

func initializePlayer() *player.MyPlayer {
	var pausedState = state.SafePlayState{}
	pausedState.SetState(state.Paused)
	var closedState = state.SafePlayState{}
	closedState.SetState(state.Closed)
	myPlayer = &player.MyPlayer{
		WantedState:     &pausedState,
		CurrentState:    &closedState,
		Commands:        &commandQueue,
		CurrentPlayable: &types.PlayableImpl{},
	}

	return myPlayer
}

func addAllFilesFromDir(path string) {
	if !dirExist(path) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			log.Fatal(err)
			return
		}
	}

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(files) == 0 {
		err := DownloadFile("music/Cease Fire.mp3", "https://legismusic.s3.amazonaws.com/epic/Cease+Fire.mp3")
		if err != nil {
			log.Fatal(err)
			return
		}
		err = DownloadFile("music/Bloody.mp3", "https://legismusic.s3.amazonaws.com/epic/Bloody.mp3")
		if err != nil {
			log.Fatal(err)
			return
		}

		files, err = os.ReadDir(path)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		myPlayer.AddSongFromDisc(&types.Song{
			Title: file.Name(),
			Path:  path + "/" + file.Name(),
		})
	}
}

func dirExist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func dispatchCommands(player *player.MyPlayer) {
	for command := range *player.Commands {
		switch command {
		case commands.Play:
			player.Play()
		case commands.Pause:
			player.Pause()
		case commands.Prev:
			player.Prev()
		case commands.Next:
			player.Next()
		default:
			fmt.Println("Unknown command" + string(rune(command)))
		}
	}
}
