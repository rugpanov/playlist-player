package commands

type PlayCommand int

const (
	Play PlayCommand = iota
	Pause
	Next
	Prev
)
