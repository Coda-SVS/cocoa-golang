package log

import "github.com/Kor-SVS/cocoa/src/util"

type logBoxEx struct {
	logger *Logger
	logBox *LogBox
}

var (
	logBoxHandleChannel chan *logBoxEx
)

func init() {
	logBoxHandleChannel = make(chan *logBoxEx, 32)

	wg := util.GetWaitGroup()
	wg.Add(1)

	go func() {
		wg := util.GetWaitGroup()
		defer wg.Done()

		for logBoxEx := range logBoxHandleChannel {
			logger := logBoxEx.logger
			logBox := logBoxEx.logBox
			logger.logBoxHandler(logBox, true)
		}

		box := NewLogBox(TRACE)
		box.SetIsAsync(false)
		box.message = "비동기 로그 처리 스레드 종료됨"
		box.AddCallStack(0)
		RootLogger().Direct(box)
	}()
}

func Dispose() {
	if logBoxHandleChannel != nil {
		close(logBoxHandleChannel)
		logBoxHandleChannel = nil
	}
}
