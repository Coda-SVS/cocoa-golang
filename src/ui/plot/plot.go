package plot

import (
	"github.com/Kor-SVS/cocoa/src/log"
)

var logger *log.Logger

func init() {
	logOption := log.NewLoggerOption()
	logOption.Prefix = "[plot]"
	logger = log.RootLogger().NewSimpleLogger(logOption)

	logger.Trace("Plot init...")
}
