package imguiw

import (
	"github.com/Kor-SVS/cocoa/src/i18n"
	"github.com/eduardolat/goeasyi18n"
)

func T(translateKey string, options ...goeasyi18n.Options) string {
	return Context.FontAtlas.RegisterString(i18n.T(translateKey, options...))
}

func RS(text string) string {
	return Context.FontAtlas.RegisterString(text)
}
