package audio

import (
	"github.com/Kor-SVS/cocoa/src/config"
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

func newDeviceConfig() malgo.DeviceConfig {
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)

	deviceName := audioDeviceConfig.GetString("DeviceName")
	var deviceInfo *malgo.DeviceInfo
	if deviceName != "" {
		deviceInfo = getDeviceInfoFromName(deviceName)
	}
	if deviceInfo == nil {
		deviceInfo = getDefaultDeviceInfo()
	}

	config.RootConfig.Set("Audio.Device.DeviceName", deviceInfo.Name())
	// audioDeviceConfig.Set("DeviceName", deviceInfo.Name())

	deviceConfig.Playback.DeviceID = deviceInfo.ID.Pointer()
	return deviceConfig
}

func getDevices() []malgo.DeviceInfo {
	devices, err := malgoContext.Devices(malgo.Playback)
	if err != nil {
		logger.Errorf("Playback Devices 로드 오류 (err=%v)", err)
		return nil
	}
	return devices
}

func getDeviceInfoFromName(name string) *malgo.DeviceInfo {
	devices := getDevices()
	if devices == nil {
		return nil
	}

	for _, device := range devices {
		if device.Name() == name {
			return &device
		}
	}

	logger.Errorf("해당하는 Playback DeviceInfo를 찾을 수 없음 (name=%v)", name)
	return nil
}

func getDefaultDeviceInfo() *malgo.DeviceInfo {
	devices := getDevices()
	if devices == nil {
		logger.Errorf("Playback DeviceInfo 로드 오류")
		return nil
	}

	for _, device := range devices {
		if device.IsDefault == 1 {
			return &device
		}
	}

	logger.Errorf("Default Playback Device를 찾을 수 없음")
	return nil
}

func disposeDevice() {
	if audioDevice != nil {
		defer audioDevice.Uninit()
		audioDevice = nil
	}
}
