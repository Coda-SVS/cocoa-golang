package log

import (
	"os"
	"path"
	"strings"
)

// 프로그램 작업 디렉토리 가져오기
func getExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return path.Dir(ex)
}

// 문자열 Join
func stringConcat(strs ...string) string {
	return strings.Join(strs, "")
}
