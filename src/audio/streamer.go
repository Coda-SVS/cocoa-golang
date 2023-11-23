package audio

import (
	"path"
)

func Open(filepath string) {
	ext := path.Ext(filepath)

	switch ext {
	case ".wav":
	case ".mp3":
	}
}
