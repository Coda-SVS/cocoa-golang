package widget

import (
	"sync"
	"time"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/Kor-SVS/cocoa/src/audio"
	"github.com/Kor-SVS/cocoa/src/audio/dsp"
	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/Kor-SVS/cocoa/src/ui/imguiw"
)

var (
	waveformPlotOnce     sync.Once
	waveformPlotInstance *WaveformPlot
)

type WaveformPlot struct {
	logger *log.Logger

	isShouldDataRefresh bool

	sampleArrayX           []float64
	sampleArrayY           []float64
	sampleArrayViewX       []float64
	sampleArrayViewY       []float64
	maxSampleCount         int
	sampleCutStartIndex    int
	sampleCutEndIndex      int
	oldSampleCutStartIndex int
	oldSampleCutEndIndex   int

	axisXLimitMax float64
	offset        int

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
			sampleArrayX:        make([]float64, 0),
			sampleArrayY:        make([]float64, 0),
			sampleArrayViewX:    make([]float64, 0),
			sampleArrayViewY:    make([]float64, 0),
			maxSampleCount:      80000,
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
		imgui.PlotFlagsNoLegend|imgui.PlotFlagsNoTitle|imgui.PlotFlagsNoMenus,
	) {
		wp.updateData()
		dataLen := len(wp.sampleArrayY)

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

		imgui.PlotSetupLock()

		if dataLen > 0 {
			wp.updateViewData()

			imgui.PlotPlotLinedoublePtrdoublePtrV(
				"WavefromData",
				&wp.sampleArrayViewX,
				&wp.sampleArrayViewY,
				int32(len(wp.sampleArrayViewY)),
				imgui.PlotLineFlagsNone,
				int32(wp.offset),
				8,
			)

			wp.audioStreamPosDragLine()

			plotSize := imgui.PlotGetPlotSize()
			plotPos := imgui.PlotGetPlotPos()

			plotEndSize := plotSize.Add(plotPos)

			plotStartPoint := imgui.PlotPixelsToPlotFloatV(plotPos.X, plotPos.Y, imgui.AxisX1, imgui.AxisY1).X
			plotEndPoint := imgui.PlotPixelsToPlotFloatV(plotEndSize.X, plotEndSize.Y, imgui.AxisX1, imgui.AxisY1).X

			sampleRate := audio.StreamFormat().SampleRate
			wp.sampleCutStartIndex = max(0, int(plotStartPoint*float64(sampleRate)))
			wp.sampleCutEndIndex = min(dataLen, int(plotEndPoint*float64(sampleRate)))
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

		wp.sampleArrayY = sampleArray
		sampleArrayX := make([]float64, sampleCount)
		for i := 0; i < sampleCount; i++ {
			sampleArrayX[i] = float64(i) / float64(sampleRate)
		}
		wp.sampleArrayX = sampleArrayX

		wp.axisXLimitMax = audio.Duration().Seconds()

		wp.sampleCutStartIndex = 0
		wp.sampleCutEndIndex = sampleCount
	} else {
		wp.clear()
	}

	wp.isShouldDataRefresh = false
}

// Plot에 표시되는 데이터 처리 (화면 밖 데이터 Cut, 다운샘플링)
func (wp *WaveformPlot) updateViewData() {
	sampleCutStartIndex := wp.sampleCutStartIndex
	sampleCutEndIndex := wp.sampleCutEndIndex
	oldSampleCutStartIndex := wp.oldSampleCutStartIndex
	oldSampleCutEndIndex := wp.oldSampleCutEndIndex

	if oldSampleCutStartIndex != sampleCutStartIndex || oldSampleCutEndIndex != sampleCutEndIndex {
		maxSampleCount := wp.maxSampleCount
		realSampleCount := len(wp.sampleArrayY)
		viewSampleCount := sampleCutEndIndex - sampleCutStartIndex

		simpleCut := func() {
			wp.sampleArrayViewX = wp.sampleArrayX[sampleCutStartIndex:sampleCutEndIndex]
			wp.sampleArrayViewY = wp.sampleArrayY[sampleCutStartIndex:sampleCutEndIndex]
		}

		if realSampleCount > maxSampleCount && viewSampleCount > maxSampleCount {
			// 다운샘플링
			simpleCut()
			sampleArrayViewX, sampleArrayViewY, err := dsp.LTTB(
				wp.sampleArrayViewX,
				wp.sampleArrayViewY,
				maxSampleCount,
			)
			logger.Tracef("wp.sampleArrayViewY: %v, viewSampleCount: %v, maxSampleCount: %v", len(sampleArrayViewY), viewSampleCount, maxSampleCount)
			if err != nil {
				logger.Errorf("다운샘플링 오류 (err=%v)", err)
			} else {
				wp.sampleArrayViewX = sampleArrayViewX
				wp.sampleArrayViewY = sampleArrayViewY
			}
		} else {
			simpleCut()
		}

		wp.oldSampleCutStartIndex = sampleCutStartIndex
		wp.oldSampleCutEndIndex = sampleCutEndIndex
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
	clear(wp.sampleArrayX)
	clear(wp.sampleArrayY)
	clear(wp.sampleArrayViewX)
	clear(wp.sampleArrayViewY)
	wp.axisXLimitMax = 30
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
