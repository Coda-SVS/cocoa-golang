package resource

import (
	"embed"
	"io/fs"

	"github.com/Kor-SVS/cocoa/src/core"
	"github.com/Kor-SVS/cocoa/src/log"
)

var logger *log.Logger

//go:embed *
var resourceFS embed.FS

func init() {
	logOption := log.NewLoggerOption()
	logOption.Prefix = "[resource]"
	logger = log.RootLogger().NewSimpleLogger(logOption)

	logger.Trace("resource init...")
}

func GetResourceFS(path string) fs.FS {
	if path == "" {
		return resourceFS
	} else {
		subFS, err := fs.Sub(resourceFS, path)
		if err != nil {
			err = core.NewErrorW(err, true)
			logger.Errorf("resourceFS 로드 오류 (err=%v, path=%v)", err, path)
			panic(err)
		}
		return subFS
	}
}
