package audio

import "github.com/sasha-s/go-deadlock"

var (
	audioMutex *deadlock.Mutex = new(deadlock.Mutex)
)
