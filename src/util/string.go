package util

import "strings"

// 문자열 Join
func StringConcat(sep string, strs ...string) string {
	return strings.Join(strs, sep)
}
