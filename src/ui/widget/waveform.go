package widget

import (
	"sync"
	"time"

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

func NewWaveformPlotData() *WaveformPlotData {
	return &WaveformPlotData{
		X: make([]float64, 0),
		Y: make([]float64, 0),
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
	logger *log.Logger

	isShouldDataRefresh bool

	sampleArray       *WaveformPlotData
	sampleArrayView   *WaveformPlotData
	maxSampleCount    int
	sampleCutIndex    *util.ArrayIndex
	sampleCutIndexOld *util.ArrayIndex

	axisXLimitMax float64
	offset        int32
	isFitRequest  bool

	audioStreamPosID           int32
	audioStreamPos             float64
	audioStreamPos_out_clicked bool
	audioStreamPos_out_hovered bool
	audioStreamPos_held        bool
	audioStreamPos_IsPaused    bool
}

// 싱글톤
func GetWaveformPlot() *WaveformPlot {
	waveformPlotOnce.Do(func() {
		waveformPlotInstance = &WaveformPlot{
			isShouldDataRefresh: true,
			sampleArray:         NewWaveformPlotData(),
			sampleArrayView:     NewWaveformPlotData(),
			sampleCutIndex:      util.NewArrayIndex(0, 0),
			sampleCutIndexOld:   util.NewArrayIndex(0, 0),
			maxSampleCount:      50000,
			axisXLimitMax:       30,
		}

		logOption := log.NewLoggerOption()
		logOption.Prefix = "[waveform]"
		waveformPlotInstance.logger = logger.NewSimpleLogger(logOption)

		waveformPlotInstance.logger.Trace("waveform init...")

		waveformPlotInstance.eventHandler_AudioStreamChanged()
	})

	return waveformPlotInstance
}

func (wp *WaveformPlot) View() {
	if imgui.PlotBeginPlotV(
		imguiw.T("Waveform"),
		imgui.Vec2{X: -1, Y: -1},
		imgui.PlotFlagsNoLegend|imgui.PlotFlagsNoTitle|imgui.PlotFlagsNoMenus|imgui.PlotFlagsNoMouseText,
	) {
		wp.updateData()
		dataLen := wp.sampleArray.LengthY()

		imgui.PlotSetupAxisV(
			imgui.AxisX1,
			"WavefromPlotX",
			imgui.PlotAxisFlags(imgui.PlotAxisFlagsNoLabel),
		)
		imgui.PlotSetupAxisV(
			imgui.AxisY1,
			"WavefromPlotY",
			imgui.PlotAxisFlags(imgui.PlotAxisFlagsLock|imgui.PlotAxisFlagsNoTickLabels|imgui.PlotAxisFlagsNoGridLines|imgui.PlotAxisFlagsNoLabel),
		)

		imgui.PlotSetupAxisLimitsConstraints(imgui.AxisX1, 0, wp.axisXLimitMax)
		imgui.PlotSetupAxisLimitsV(imgui.AxisY1, -0.5, 0.5, imgui.PlotCondAlways)

		if wp.isFitRequest {
			imgui.PlotSetupAxisLimitsV(imgui.AxisX1, 0, wp.axisXLimitMax, imgui.PlotCondAlways)
			wp.isFitRequest = false
		}

		imgui.PlotSetupLock()

		if dataLen > 0 {
			mouseWheelDelta := imguiw.Context.IO().MouseWheel()
			if mouseWheelDelta == 0 {
				wp.updateViewData()
			}

			imgui.PlotPlotLinedoublePtrdoublePtrV(
				"WavefromData",
				&wp.sampleArrayView.X,
				&wp.sampleArrayView.Y,
				int32(wp.sampleArrayView.LengthY()),
				imgui.PlotLineFlagsNone,
				wp.offset,
				8,
			)

			wp.audioStreamPosDragLine()

			plotSize := imgui.PlotGetPlotSize()
			plotPos := imgui.PlotGetPlotPos()

			plotEndSize := plotSize.Add(plotPos)

			plotStartPoint := imgui.PlotPixelsToPlotFloatV(plotPos.X, plotPos.Y, imgui.AxisX1, imgui.AxisY1).X
			plotEndPoint := imgui.PlotPixelsToPlotFloatV(plotEndSize.X, plotEndSize.Y, imgui.AxisX1, imgui.AxisY1).X

			sampleRate := audio.StreamFormat().SampleRate
			wp.sampleCutIndex.Start = max(0, int(plotStartPoint*float64(sampleRate)))
			wp.sampleCutIndex.End = min(dataLen, int(plotEndPoint*float64(sampleRate)))
		}

		if imgui.PlotIsPlotHovered() && imgui.IsMouseDoubleClicked(imgui.MouseButtonLeft) {
			wp.isFitRequest = true
		}

		imgui.PlotEndPlot()
	}
}

// 현재 오디오 스트림에서 데이터 불러오기
// TODO: 메모리 최적화 (현재 약 4분 44100Hz 오디오의 경우 500~800MB의 메모리 소모)
func (wp *WaveformPlot) updateData() {
	if !wp.isShouldDataRefresh {
		return
	}

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

	wp.isFitRequest = true
	wp.isShouldDataRefresh = false
}

// Plot에 표시되는 데이터 처리 (화면 밖 데이터 Cut, 다운샘플링)
func (wp *WaveformPlot) updateViewData() {
	sampleCutIndex := *wp.sampleCutIndex
	sampleCutIndexOld := *wp.sampleCutIndexOld

	if !sampleCutIndex.Equal(sampleCutIndexOld) {
		maxSampleCount := wp.maxSampleCount
		realSampleCount := wp.sampleArray.LengthY()
		viewSampleCount := sampleCutIndex.Size()

		simpleCut := func() {
			wp.sampleArrayView.X = wp.sampleArray.X[sampleCutIndex.Start:sampleCutIndex.End]
			wp.sampleArrayView.Y = wp.sampleArray.Y[sampleCutIndex.Start:sampleCutIndex.End]
		}

		if realSampleCount > maxSampleCount && viewSampleCount > maxSampleCount {
			// 다운샘플링
			simpleCut()
			sampleArrayViewX, sampleArrayViewY, err := dsp.LTTB(
				wp.sampleArrayView.X,
				wp.sampleArrayView.Y,
				maxSampleCount,
			)
			if err != nil {
				logger.Errorf("다운샘플링 오류 (err=%v)", err)
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

// 재생 위치 표시 및 조작
func (wp *WaveformPlot) audioStreamPosDragLine() {
	if !audio.IsAudioLoaded() {
		return
	}

	wp.audioStreamPos = audio.Position().Seconds()
	_audioStreamPos := wp.audioStreamPos
	imgui.PlotDragLineXV(
		wp.audioStreamPosID,
		&wp.audioStreamPos,
		imgui.NewVec4(255, 0, 0, 255),
		1,
		imgui.PlotDragToolFlagsDelayed,
		&wp.audioStreamPos_out_clicked,
		&wp.audioStreamPos_out_hovered,
		&wp.audioStreamPos_held,
	)

	if wp.audioStreamPos_held {
		if !wp.audioStreamPos_IsPaused && audio.IsRunning() {
			audio.Pause()
			wp.audioStreamPos_IsPaused = true
		}
		if _audioStreamPos != wp.audioStreamPos {
			audio.SetPosition(time.Duration(wp.audioStreamPos * float64(time.Second)))
		}
	}

	if wp.audioStreamPos_out_clicked && wp.audioStreamPos_IsPaused {
		audio.Play()
		wp.audioStreamPos_IsPaused = false
	}
}

func (wp *WaveformPlot) clear() {
	wp.sampleArray.Clear()
	wp.sampleArrayView.Clear()
	wp.axisXLimitMax = 30
	wp.sampleCutIndex = util.NewArrayIndex(0, 0)
	wp.sampleCutIndexOld = util.NewArrayIndex(0, 0)
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
