package audio

import (
	"github.com/gen2brain/malgo"
)

var (
	audioDevice *malgo.Device
	OnStop      func()
)

func initDevice(config malgo.DeviceConfig) {
	deviceCallbacks := malgo.DeviceCallbacks{
		Data: callbackData,
		Stop: callbackStop,
	}

	device, err := malgo.InitDevice(malgoContext.Context, config, deviceCallbacks)
	if err != nil {
		logger.Error("Device 초기화 실패 err=%w", err)
	}

	audioDevice = device
}

func callbackData(pOutputSample, pInputSamples []byte, frameCount uint32) {
	if frameCount == 0 {
		return
	}

	readN := readAudioStream(pOutputSample, int(frameCount))
	if readN <= 0 {
		return
	}
}

func callbackStop() {
	if OnStop != nil {
		OnStop()
	}
}

func disposeDevice() {
	if audioDevice != nil {
		defer audioDevice.Uninit()
		audioDevice = nil
	}
}
