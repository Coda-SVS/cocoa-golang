package util

import (
	"strings"

	"golang.org/x/text/language"
)

// 문자열 Join
func StringConcat(sep string, strs ...string) string {
	return strings.Join(strs, sep)
}

func NormalizeLanguageCode(s string) string {
	matcher := language.NewMatcher([]language.Tag{
		language.English,
		language.Arabic,
		language.French,
		language.German,
		language.Hindi,
		language.Indonesian,
		language.Italian,
		language.Japanese,
		language.Korean,
		language.Malay,
		language.Portuguese,
		language.Russian,
		language.SimplifiedChinese,
		language.TraditionalChinese,
		language.Chinese,
		language.Spanish,
		language.Thai,
		language.Turkish,
		language.Vietnamese,
	})

	preferred, _ := language.Parse(s)
	code, _, _ := matcher.Match(preferred)
	base, _ := code.Base()
	langauge, _ := code.Region()

	if base.String() == "zh" { //중국어 예외처리
		return "zh-" + strings.ToLower(langauge.String())
	}

	return base.String()
}
