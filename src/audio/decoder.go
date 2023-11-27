package audio

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kor-SVS/cocoa/src/util"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/wav"
)

type Format string

var (
	UNKNOWN Format = "unknown"
	WAV     Format = "wav"
	FLAC    Format = "flac"
	MP3     Format = "mp3"
)

func GetDecoder(fpath string) (beep.StreamSeekCloser, *beep.Format, error) {
	ext_format, err := FileFormat(fpath)
	if err != nil {
		return nil, nil, err
	}

	f, err := os.Open(fpath)
	if err != nil {
		return nil, nil, err
	}

	var decoder beep.StreamSeekCloser
	var format beep.Format

	switch ext_format {
	case WAV:
		decoder, format, err = wav.Decode(f)
		if err != nil {
			return nil, nil, err
		}
	case FLAC:
		decoder, format, err = flac.Decode(f)
		if err != nil {
			return nil, nil, err
		}
	case MP3:
		decoder, format, err = mp3.Decode(f)
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, fmt.Errorf("파일 형식을 알 수 없습니다. fpath=%v", fpath)
	}

	return decoder, &format, nil
}

// FileFormat returns the known format of the passed path.
func FileFormat(fpath string) (Format, error) {
	if !util.FileExists(fpath) {
		return "", errors.New("invalid path")
	}

	ext := strings.ToLower(filepath.Ext(fpath))
	switch ext {
	case ".wav", ".wave":
		return WAV, nil
	case ".flac":
		return FLAC, nil
	case ".mp3":
		return MP3, nil
	}

	return UNKNOWN, nil
}
