package audio

import (
	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/gen2brain/malgo"
)

var logger *log.Logger

func init() {
	logOption := log.NewLoggerOption()
	logOption.Prefix = "[audio]"
	logWriter := log.NewLogWriter(nil, nil, nil, nil)
	logger = log.NewLogger(logOption, logWriter)

	logger.Trace("Audio init...")

	audioMutex.Lock()
	defer audioMutex.Unlock()

	initContext(nil, malgo.ContextConfig{})
}

// 할당된 자원 해제
func Dispose() {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	defer disposeDevice()
	defer disposeContext()
	defer close()
}
