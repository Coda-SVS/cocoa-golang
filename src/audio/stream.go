package audio

import (
	"github.com/gen2brain/malgo"
	"github.com/gopxl/beep"
)

type AudioStream struct {
	beep.StreamSeekCloser
	Format *beep.Format
}

var (
	audioStream             *AudioStream                // 오디오 데이터 스트림
	bufferCallbackFuncArray []func(buffer [][2]float64) // 오디오 버퍼 접근 함수 콜백
	audioBuffer             [][2]float64                // 오디오 버퍼
)

func Open(fpath string) {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	logger.Trace("오디오 파일을 여는 중...")

	decoder, format, err := GetDecoder(fpath)
	if err != nil {
		logger.Errorf("오디오 파일을 열지 못했습니다. err=%v", err)
		decoder.Close()
	}

	if decodeErr := decoder.Err(); decodeErr != nil {
		logger.Errorf("디코드 오류 발생 decodeErr=%v", decodeErr)
		decoder.Close()
	}

	audioStream = &AudioStream{}
	audioStream.StreamSeekCloser = decoder
	audioStream.Format = format
	audioBuffer = nil

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatF32
	deviceConfig.Playback.Channels = 2
	deviceConfig.SampleRate = uint32(audioStream.Format.SampleRate)
	deviceConfig.Alsa.NoMMap = 1

	initDevice(deviceConfig)
}

func readAudioStream(outBuffer []byte, frameCount int) int {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	if audioStream == nil {
		return 0
	}

	if audioBuffer == nil || len(audioBuffer) != frameCount {
		audioBuffer = make([][2]float64, frameCount)
	}

	sampleLen := audioStream.Len()
	readN, ok := audioStream.Stream(audioBuffer)

	for _, callback := range bufferCallbackFuncArray {
		callback(audioBuffer)
	}

	floatSampleToByteArray(audioBuffer, outBuffer)

	if readN == sampleLen && ok { // 스트림이 끝났을 경우
		return readN
	}
	return readN
}

func Close() {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	close()
}

func close() {
	if audioStream != nil {
		disposeDevice()
		audioStream.Close()
		audioStream = nil
		audioBuffer = nil
	}
}
