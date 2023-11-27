package audio

import (
	"github.com/Kor-SVS/cocoa/src/util"
	"github.com/gen2brain/malgo"
	"github.com/gopxl/beep"
)

type AudioStream struct {
	beep.StreamSeekCloser
	Format *beep.Format
}

var (
	AudioStreamBroker *util.Broker[EnumAudioStreamState] = util.NewBroker[EnumAudioStreamState]()
)

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

	deviceConfig := newDeviceConfig()

	deviceConfig.Playback.Format = malgo.FormatF32
	deviceConfig.Playback.Channels = 2
	deviceConfig.SampleRate = uint32(audioStream.Format.SampleRate)

	initDevice(deviceConfig)

	AudioStreamBroker.Publish(EnumAudioStreamOpen)
}

func GetAllSampleData() [][2]float64 {
	audioMutex.Lock()
	defer audioMutex.Unlock()

	pos := audioStream.Position()
	audioStream.Seek(0)

	buf := make([][2]float64, audioStream.Len())
	audioStream.Stream(buf)

	audioStream.Seek(pos)

	return buf
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

	util.FloatSampleToByteArray(audioBuffer, outBuffer)

	if readN == sampleLen && ok { // 스트림이 끝났을 경우
		return readN
	}
	return readN
}

func isAudioLoaded() bool {
	return audioStream != nil && audioDevice != nil
}

func Close() {
	audioMutex.Lock()
	defer func() {
		close()

		AudioStreamBroker.Publish(EnumAudioStreamClosed)
	}()

	defer audioMutex.Unlock()
}

func close() {
	if audioStream != nil {
		disposeDevice()
		audioStream.Close()
		audioStream = nil
		audioBuffer = nil
	}
}
