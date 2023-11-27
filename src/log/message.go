package log

import "github.com/Kor-SVS/cocoa/src/util"

type logBoxInternal struct {
	joinedPrefix string
}

type LogBox struct {
	logBoxInternal logBoxInternal
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

func (lm *LogBox) Message() string {
	return lm.message
}

func (lm *LogBox) BuildMessage() string {
	return util.StringConcat(lm.logBoxInternal.joinedPrefix, " ", lm.message)
}
