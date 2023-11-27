package audio

import (
	"github.com/Kor-SVS/cocoa/src/log"
)

var (
	logger *log.Logger
)

func init() {
	logOption := log.NewLoggerOption()
	logOption.Prefix = "[audio]"
	logWriter := log.NewLogWriter(nil, nil, nil, nil)
	logger = log.NewLogger(logOption, logWriter)

	logger.Trace("Audio init...")

	configInit()

	audioMutex.Lock()
	defer audioMutex.Unlock()

	initContext()
}

// 할당된 자원 해제
func Dispose() {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	defer disposeDevice()
	defer disposeContext()
	defer close()
}
