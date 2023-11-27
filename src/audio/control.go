package audio

import "time"

func Play() {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if audioStream == nil || audioDevice == nil {
		return
	}

	if !audioDevice.IsStarted() {
		audioDevice.Start()
	}
}

func Pause() {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if audioStream == nil || audioDevice == nil {
		return
	}

	if audioDevice.IsStarted() {
		audioDevice.Stop()
	}
}

func Stop() {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if audioStream == nil || audioDevice == nil {
		return
	}

	if audioDevice.IsStarted() {
		audioDevice.Stop()
	}
}

func SetPosition(t time.Duration) {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if audioStream == nil || audioDevice == nil {
		return
	}

	audioStream.Seek(int(t.Seconds() * float64(audioStream.Format.SampleRate)))
}

func Position() *time.Duration {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if audioStream == nil {
		return nil
	}

	rawPos := float64(audioStream.Position()) / float64(audioStream.Format.SampleRate)
	pos := time.Duration(rawPos * float64(time.Second))
	return &pos
}

func Duration() *time.Duration {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if audioStream == nil {
		return nil
	}

	rawDur := float64(audioStream.Len()) / float64(audioStream.Format.SampleRate)
	dur := time.Duration(rawDur * float64(time.Second))
	return &dur
}
