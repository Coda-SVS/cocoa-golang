package imguiw

import (
	"github.com/Kor-SVS/cocoa/src/i18n"
	"github.com/eduardolat/goeasyi18n"
)

func T(translateKey string, options ...goeasyi18n.Options) string {
	tString := i18n.T(translateKey, options...)
	if tString == "" {
		tString = translateKey
	}
	return Context.FontAtlas.RegisterString(tString)
}

func RS(text string) string {
	return Context.FontAtlas.RegisterString(text)
}
