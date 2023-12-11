package i18n

import "github.com/eduardolat/goeasyi18n"

func init() {
	langLabel := "한국어 (Korean)"
	languageData := LanguageData{
		Code:       "ko",
		Label:      langLabel,
		DataGetter: getKoTranslateMap,
	}

	languageNameList = append(languageNameList, langLabel)
	languageDataMap["ko"] = &languageData
}

func getKoTranslateMap() goeasyi18n.TranslateStrings {
	result := []goeasyi18n.TranslateString{{
		Key:     "File",
		Default: "파일",
	}, {
		Key:     "OpenFile",
		Default: "파일열기",
	}, {
		Key:     "Audio",
		Default: "오디오",
	}, {
		Key:     "AudioPlay",
		Default: "재생",
	}, {
		Key:     "AudioPause",
		Default: "일시정지",
	}, {
		Key:     "AudioStop",
		Default: "중지",
	},
	}

	return result
}
