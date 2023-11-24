package audio

import "github.com/Kor-SVS/cocoa/src/log"

var logger *log.Logger

func init() {
	logOption := log.NewLoggerOption()
	logOption.Prefix = "[audio]"
	logWriter := log.NewLogWriter(nil, nil, nil, nil)
	logger = log.NewLogger(logOption, logWriter)

	logger.Trace("Audio init...")
}
