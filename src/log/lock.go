package log

import "github.com/sasha-s/go-deadlock"

var logMutex *deadlock.Mutex = new(deadlock.Mutex)
