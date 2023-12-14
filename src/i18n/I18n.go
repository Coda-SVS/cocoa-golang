package i18n

import (
	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/Kor-SVS/cocoa/src/util"
	"golang.org/x/text/language"

	"github.com/eduardolat/goeasyi18n"
)

var (
	logger           *log.Logger
	currentLangCode  string
	i18nInstance     *goeasyi18n.I18n
	languageNameList []string
	languageDataMap  map[string]*LanguageData
)

type LanguageData struct {
	Code       string
	Label      string
	DataGetter func() goeasyi18n.TranslateStrings
}

func init() {
	logOption := log.NewLoggerOption()
	logOption.Prefix = "[i18n]"
	logger = log.RootLogger().NewSimpleLogger(logOption)

	logger.Trace("I18n init...")

	currentLangCode = util.NormalizeLanguageCode(language.Korean.String())

	i18nInstance = goeasyi18n.NewI18n(goeasyi18n.Config{
		FallbackLanguageName: currentLangCode,
	})

	i18nInstance.AddLanguage(currentLangCode, getKoTranslateMap())

	languageNameList = make([]string, 1)
	languageDataMap = make(map[string]*LanguageData)
}

func GetLanguageList() []string {
	return languageNameList
}

func GetLanguageMap() map[string]*LanguageData {
	return languageDataMap
}

func SetLanguage(lang string) {
	normalizedLang := util.NormalizeLanguageCode(lang)
	languageData, ok := languageDataMap[normalizedLang]

	if ok {
		if !i18nInstance.HasLanguage(normalizedLang) {
			i18nInstance.AddLanguage(normalizedLang, languageData.DataGetter())
		}
		currentLangCode = lang
		logger.Tracef("언어 적용됨 (lang=%v, normalizedLang=%v)", lang, normalizedLang)
	} else {
		logger.Errorf("언어 데이터가 존재하지 않음 (lang=%v, normalizedLang=%v)", lang, normalizedLang)
	}
}

func T(translateKey string, options ...goeasyi18n.Options) string {
	return i18nInstance.Translate(currentLangCode, translateKey, options...)
}
