package log

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/Kor-SVS/cocoa/src/util"
)

type logBoxInternal struct {
	joinedPrefix string
}

type LogBox struct {
	logBoxInternal logBoxInternal
	callerLocation string
	message        string
}

func NewLogBox() *LogBox {
	return &LogBox{
		logBoxInternal: logBoxInternal{
			joinedPrefix: "",
		},
		message: "",
	}
}

func (lm *LogBox) AddCallStack(skip int) bool {
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

func (lm *LogBox) Message() string {
	return lm.message
}

func (lm *LogBox) BuildMessage() string {
	if lm.callerLocation == "" {
		return util.StringConcat(" ", lm.logBoxInternal.joinedPrefix, lm.message)
	} else {
		return util.StringConcat(" ", lm.callerLocation, lm.logBoxInternal.joinedPrefix, lm.message)
	}
}
