package widget

import (
	"sync"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/Kor-SVS/cocoa/src/audio"
	"github.com/Kor-SVS/cocoa/src/audio/dsp"
	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/Kor-SVS/cocoa/src/ui/imguiw"
	"github.com/Kor-SVS/cocoa/src/util"
)

var (
	waveformPlotOnce     sync.Once
	waveformPlotInstance *WaveformPlot
)

type WaveformPlotData struct {
	X []float64
	Y []float64
}

func NewWaveformPlotData(length, capacity int) *WaveformPlotData {
	return &WaveformPlotData{
		X: make([]float64, length, capacity),
		Y: make([]float64, length, capacity),
	}
}

func (wd *WaveformPlotData) LengthX() int {
	return len(wd.X)
}

func (wd *WaveformPlotData) LengthY() int {
	return len(wd.Y)
}

func (wd *WaveformPlotData) Clear() {
	clear(wd.X)
	clear(wd.Y)
}

type WaveformPlot struct {
	imguiw.PlotWidget

	logger *log.Logger

	title string

	isShouldDataRefresh bool

	sampleArray       *WaveformPlotData
	sampleArrayView   *WaveformPlotData
	maxSampleCount    int
	sampleCutIndex    *util.Index
	sampleCutIndexOld *util.Index

	axisXLimitMax float64
	offset        int32

	plotDrawEndEventArgs *util.PlotDrawEndEventArgs
}

// 싱글톤
func GetWaveformPlot() *WaveformPlot {
	waveformPlotOnce.Do(func() {
		waveformPlotInstance = &WaveformPlot{
			title:               "Wavefrom Data",
			isShouldDataRefresh: true,
			sampleArray:         NewWaveformPlotData(0, 0),
			sampleCutIndex:      util.NewIndex(0, 0),
			sampleCutIndexOld:   util.NewIndex(0, 0),
			maxSampleCount:      50000,
			axisXLimitMax:       DefaultAxisXLimitMax,
		}

		waveformPlotInstance.sampleArrayView = NewWaveformPlotData(0, waveformPlotInstance.maxSampleCount)

		logOption := log.NewLoggerOption()
		logOption.Prefix = "[waveform]"
		waveformPlotInstance.logger = plotLogger.NewSimpleLogger(logOption)

		waveformPlotInstance.logger.Trace("Waveform Plot init...")

		waveformPlotInstance.eventHandler_AudioStreamChanged()
	})

	return waveformPlotInstance
}

func (wp *WaveformPlot) Plot() {
	dataLen := wp.sampleArray.LengthY()

	if dataLen > 0 {
		imgui.PlotPlotLinedoublePtrdoublePtrV(
			wp.title,
			&wp.sampleArrayView.X,
			&wp.sampleArrayView.Y,
			int32(wp.sampleArrayView.LengthY()),
			imgui.PlotLineFlagsNone,
			wp.offset,
			8,
		)

		plotDrawEndEventArgs := *wp.plotDrawEndEventArgs

		sampleRate := audio.StreamFormat().SampleRate
		wp.sampleCutIndex.Start = max(0, int(plotDrawEndEventArgs.PlotPointStart*float64(sampleRate)))
		wp.sampleCutIndex.End = min(dataLen, int(plotDrawEndEventArgs.PlotPointEnd*float64(sampleRate)))
	}
}

// 현재 오디오 스트림에서 데이터 불러오기
// TODO: 메모리 최적화 (현재 약 4분 44100Hz 오디오의 경우 500~800MB의 메모리 소모)
func (wp *WaveformPlot) UpdateData() {
	if wp.isShouldDataRefresh {
		if audio.IsAudioLoaded() {
			format := audio.StreamFormat()
			sampleArray := dsp.StereoToMono(audio.GetAllSampleData())
			sampleRate := int(format.SampleRate)
			sampleCount := len(sampleArray)

			wp.sampleArray.Y = sampleArray
			sampleArrayX := make([]float64, sampleCount)
			for i := 0; i < sampleCount; i++ {
				sampleArrayX[i] = float64(i) / float64(sampleRate)
			}
			wp.sampleArray.X = sampleArrayX

			wp.axisXLimitMax = audio.Duration().Seconds()

			wp.sampleCutIndex.Start = 0
			wp.sampleCutIndex.End = sampleCount
		} else {
			wp.clear()
		}

		wp.isShouldDataRefresh = false
	}

	wp.updateViewData()
}

// Plot에 표시되는 데이터 처리 (화면 밖 데이터 Cut, 다운샘플링)
func (wp *WaveformPlot) updateViewData() {
	if !audio.IsAudioLoaded() {
		return
	}

	sampleCutIndex := *wp.sampleCutIndex
	sampleCutIndexOld := *wp.sampleCutIndexOld

	if !sampleCutIndex.Equal(sampleCutIndexOld) {
		maxSampleCount := wp.maxSampleCount
		realSampleCount := wp.sampleArray.LengthY()
		viewSampleCount := sampleCutIndex.Size()

		simpleCut := func() {
			sampleX := wp.sampleArray.X[sampleCutIndex.Start:sampleCutIndex.End]
			sampleY := wp.sampleArray.Y[sampleCutIndex.Start:sampleCutIndex.End]
			copy(wp.sampleArrayView.X, sampleX)
			copy(wp.sampleArrayView.Y, sampleY)
		}

		if realSampleCount > maxSampleCount && viewSampleCount > maxSampleCount {
			// 다운샘플링
			sampleArrayViewX, sampleArrayViewY, err := dsp.LTTB_Buffer(
				wp.sampleArray.X[sampleCutIndex.Start:sampleCutIndex.End],
				wp.sampleArray.Y[sampleCutIndex.Start:sampleCutIndex.End],
				wp.sampleArrayView.X,
				wp.sampleArrayView.Y,
				maxSampleCount,
			)
			if err != nil {
				logger.Errorf("다운샘플링 오류 (err=%v)", err)
				simpleCut()
			} else {
				wp.sampleArrayView.X = sampleArrayViewX
				wp.sampleArrayView.Y = sampleArrayViewY
			}
		} else {
			simpleCut()
		}

		wp.sampleCutIndexOld = &sampleCutIndex
	}
}

func (wp *WaveformPlot) Title() string {
	return wp.title
}

func (wp *WaveformPlot) EventHandler(eventArgs any) {
	switch castEventArgs := eventArgs.(type) {
	case util.PlotDrawEndEventArgs:
		wp.plotDrawEndEventArgs = &castEventArgs
	}
}

func (wp *WaveformPlot) clear() {
	wp.sampleArray.Clear()
	wp.sampleArrayView.Clear()
	wp.axisXLimitMax = 30
	wp.sampleCutIndex = util.NewIndex(0, 0)
	wp.sampleCutIndexOld = util.NewIndex(0, 0)
}

func (wp *WaveformPlot) IsDisposed() bool {
	return false
}

// 오디오 스트림의 이벤트 수신
func (wp *WaveformPlot) eventHandler_AudioStreamChanged() {
	go func() {
		msgChan := audio.AudioStreamBroker().Subscribe()

		for msg := range msgChan {
			wp.logger.Tracef("[Callback] AudioStreamChanged (msg=%v)", msg)
			switch msg {
			case audio.EnumAudioStreamOpen:
				wp.isShouldDataRefresh = true
			case audio.EnumAudioStreamClosed:
				wp.isShouldDataRefresh = true
			}
		}
	}()
}
