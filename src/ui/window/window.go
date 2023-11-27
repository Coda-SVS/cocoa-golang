package window

import (
	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/sasha-s/go-deadlock"
)

var logger *log.Logger

var windowMutex *deadlock.RWMutex = new(deadlock.RWMutex)

func init() {
	logOption := log.NewLoggerOption()
	logOption.Prefix = "[window]"
	logWriter := log.NewLogWriter(nil, nil, nil, nil)
	logger = log.NewLogger(logOption, logWriter)

	logger.Trace("Window init...")
}
