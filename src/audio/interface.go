package audio

import (
	"time"

	"github.com/gopxl/beep"
)

func IsAudioLoaded() bool {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	return isAudioLoaded()
}

func IsRunning() bool {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	return audioDevice.IsStarted()
}

func Play() {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if !isAudioLoaded() {
		return
	}

	if !audioDevice.IsStarted() {
		err := audioDevice.Start()
		if err != nil {
			logger.Errorf("오디오 재생 실패 (err=%v)", err)
		}
	}
}

func Pause() {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if !isAudioLoaded() {
		return
	}

	if audioDevice.IsStarted() {
		err := audioDevice.Stop()
		if err != nil {
			logger.Errorf("오디오 중지 실패 (err=%v)", err)
		}
	}
}

func Stop() {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if !isAudioLoaded() {
		return
	}

	if audioDevice.IsStarted() {
		err := audioDevice.Stop()
		if err != nil {
			logger.Errorf("오디오 중지 실패 (err=%v)", err)
			return
		}
	}
	err := audioStream.Seek(0)
	if err != nil {
		logger.Errorf("오디오 위치 이동 실패 (err=%v)", err)
		return
	}
}

func SetPosition(t time.Duration) {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if !isAudioLoaded() {
		return
	}

	isRunning := audioDevice.IsStarted()
	if isRunning {
		err := audioDevice.Stop()
		if err != nil {
			logger.Errorf("오디오 중지 실패 (err=%v)", err)
			return
		}
	}

	err := audioStream.Seek(int(t.Seconds() * float64(audioStream.Format.SampleRate)))
	if err != nil {
		logger.Errorf("오디오 위치 이동 실패 (err=%v)", err)
		return
	}

	if isRunning {
		err := audioDevice.Start()
		if err != nil {
			logger.Errorf("오디오 재생 실패 (err=%v)", err)
		}
	}
}

func Position() time.Duration {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if !isAudioLoaded() {
		return 0
	}

	rawPos := float64(audioStream.Position()) / float64(audioStream.Format.SampleRate)
	pos := time.Duration(rawPos * float64(time.Second))
	return pos
}

func Duration() time.Duration {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if !isAudioLoaded() {
		return 0
	}

	rawDur := float64(audioStream.Len()) / float64(audioStream.Format.SampleRate)
	dur := time.Duration(rawDur * float64(time.Second))
	return dur
}

func StreamFormat() *beep.Format {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if !isAudioLoaded() {
		return nil
	}

	format := beep.Format{}
	format.NumChannels = audioStream.Format.NumChannels
	format.Precision = audioStream.Format.Precision
	format.SampleRate = audioStream.Format.SampleRate

	return &format
}
