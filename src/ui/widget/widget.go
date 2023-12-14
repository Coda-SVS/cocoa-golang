package widget

import (
	"github.com/Kor-SVS/cocoa/src/log"
)

var (
	logger *log.Logger
)

func init() {
	logOption := log.NewLoggerOption()
	logOption.Prefix = "[widget]"
	logger = log.RootLogger().NewSimpleLogger(logOption)

	logger.Trace("Widget init...")
}
