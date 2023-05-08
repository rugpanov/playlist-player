# Go Playlist Player

This module provides a way to manage a playlist with the following endpoints:

- `/player/play`: starts the playback
- `/player/pause`: pauses the playback
- `/playlist/upload`: adds a new song to the end of the playlist
- `/player/next`: plays the next song in the playlist
- `/player/prev`: plays the previous song in the playlist

## Technical Requirements

- The module have a well-defined interface for interacting with the playlist
- Playing a song does not block the control methods
- The `/player/play` GET endpoint starts playback
- The next song starts playing automatically after the current song ends
- The `/player/pause` GET endpoint pauses the current playback, and when playback is resumed, it continues from the
  pause point
- The `/playlist/upload` POST endpoint adds a new song to the end of the playlist. Implementation should be done
  considering concurrent access
- The `/player/next` GET endpoint should stop the current playback and start playing the next song. Therefore, the
  current playback should be stopped, and the next song should start playing
- The `/player/prev` GET endpoint should stop the current playback and start playing the previous song
- The implementation of all endpoints should be done considering concurrent access
- All implementations must be thoroughly tested and optimized for performance.

## How to run locally

```shell
go run .
```

## Music
Music is provided by [Legis Music](https://legismusic.com/).