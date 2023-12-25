package log

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/Kor-SVS/cocoa/src/util"
)

type logBoxInternal struct {
	joinedPrefix string
}

type LogBox struct {
	isLocked       bool // 로그 정보 동결
	isAsync        bool // 비동기 로그처리 여부
	logLevel       LogLevel
	goid           int
	logBoxInternal logBoxInternal
	callerLocation string
	message        string
}

func NewLogBox(logLevel LogLevel) *LogBox {
	lm := &LogBox{
		isLocked: false,
		isAsync:  true,
		logLevel: logLevel,
		logBoxInternal: logBoxInternal{
			joinedPrefix: "",
		},
		message: "",
	}
	lm.AddGoID()
	return lm
}

func (lm *LogBox) AddCallStack(skip int) bool {
	if lm.isLocked {
		return false
	}

	_, fileName, lineNum, ok := runtime.Caller(skip + 1)
	if strings.Contains(fileName, "src") {
		fileName = strings.Split(fileName, "src")[1][1:]
	}

	if !ok {
		return false
	} else {
		lm.callerLocation = util.StringConcat(":", fileName, strconv.Itoa(lineNum))
		return true
	}
}

func (lm *LogBox) AddCallStackFromError(err *util.ErrorW, callStackSkip int) bool {
	if lm.isLocked {
		return false
	}

	ok := err.AddCallStack(callStackSkip+1, false)

	if !ok {
		return false
	} else {
		lm.callerLocation = err.CallFile()
		return true
	}
}

func (lm *LogBox) AddGoID() {
	if lm.isLocked {
		return
	}

	lm.goid = int(util.GetCurrentGoID())
}

func (lm *LogBox) SetLock() {
	if lm.isLocked {
		return
	}

	lm.isLocked = true
}

func (lm *LogBox) IsAsync() bool {
	return lm.isAsync
}

func (lm *LogBox) SetIsAsync(value bool) {
	if lm.isLocked {
		return
	}

	lm.isAsync = value
}

func (lm *LogBox) LogLevel() LogLevel {
	return lm.logLevel
}

func (lm *LogBox) Message() string {
	return lm.message
}

func (lm *LogBox) BuildMessage(isConsole bool) (rMsg string) {
	msgBuf := make([]string, 0, 3)

	if lm.callerLocation != "" {
		msgBuf = append(msgBuf, lm.callerLocation)
	}

	if lm.goid > 0 {
		if isConsole {
			msgBuf = append(msgBuf, util.GetColorForID(lm.goid).Sprintf("[grtn %v]", lm.goid))
		} else {
			msgBuf = append(msgBuf, fmt.Sprintf("[grtn %v]", lm.goid))
		}
	}

	msgBuf = append(msgBuf, lm.logBoxInternal.joinedPrefix)
	msgBuf = append(msgBuf, lm.message)

	return util.StringConcat(" ", msgBuf...)
}
