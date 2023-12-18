package util

import (
	"runtime"
	"strconv"
	"strings"
)

func GetCallFileFromCallStack(skip int) (string, bool) {
	i := skip + 1
	for {
		_, fileName, lineNum, ok := runtime.Caller(i)
		if !ok {
			return "", false
		}

		if strings.Contains(fileName, "cocoa") && strings.Contains(fileName, "src") {
			fileName = strings.Split(fileName, "src")[1][1:]
			return StringConcat(":", fileName, strconv.Itoa(lineNum)), true
		}
		i++
	}

}
