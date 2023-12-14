package window

import (
	"github.com/Kor-SVS/cocoa/src/log"
)

var logger *log.Logger

func init() {
	logOption := log.NewLoggerOption()
	logOption.Prefix = "[window]"
	logger = log.RootLogger().NewSimpleLogger(logOption)

	logger.Trace("Window init...")
}
