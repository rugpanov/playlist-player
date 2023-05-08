package player

import (
	"bytes"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"go-playlist-player/types"
	"os"
)

func createPlayer(songPath string) *types.PlayableImpl {
	// Read the mp3 file into memory
	fileBytes, err := os.ReadFile(songPath)

	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}

	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(fileBytes)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	initOtoContextIfRequired()

	// Create a new 'playable' that will handle our sound. Paused by default.
	internalPlayer := otoContext.NewPlayer(decodedMp3)

	var p = &types.PlayableImpl{}
	p.SetPlayable(internalPlayer)

	return p
}

var otoContext *oto.Context

func initOtoContextIfRequired() {
	if otoContext != nil {
		return
	}

	// Prepare an Oto context (this will use your default audio device) that will
	// Play all our sounds. Its configuration can't be changed later.

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	samplingRate := 44100

	// Number of channels (aka locations) to Play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	numOfChannels := 2

	// Bytes used by a channel to represent one sample. Either 1 or 2 (usually 2).
	audioBitDepth := 2

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	otoContext = otoCtx
}
