package dsp

import (
	"runtime"

	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/panjf2000/ants/v2"
)

var (
	logger *log.Logger
	goPool *ants.Pool
)

func init() {
	logOption := log.NewLoggerOption()
	logOption.Prefix = "[dsp]"
	logger = log.RootLogger().NewSimpleLogger(logOption)

	logger.Trace("DSP init...")

	pool, err := ants.NewPool(runtime.NumCPU())
	if err != nil {
		logger.Errorf("DSP goPool 초기화 실패 (err=%v)", err)
	}
	goPool = pool
}

func Dispose() {
	defer goPool.Release()
}
