package widget

import (
	"sync"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/Kor-SVS/cocoa/src/audio"
	"github.com/Kor-SVS/cocoa/src/audio/dsp"
	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/Kor-SVS/cocoa/src/ui/imguiw"
	"github.com/Kor-SVS/cocoa/src/util"
	"github.com/sasha-s/go-deadlock"
)

var (
	waveformPlotOnce     sync.Once
	waveformPlotInstance *WaveformPlot
)

type WaveformPlot struct {
	imguiw.PlotWidget

	logger *log.Logger

	title string

	sctx                *util.SimpleContext
	mtx                 *deadlock.RWMutex
	isShouldDataRefresh bool
	isCleard            bool

	sampleArray       *util.WaveformPlotData
	sampleArrayView   *util.WaveformPlotData
	maxSampleCount    int
	sampleCutIndex    *util.Index
	sampleCutIndexOld *util.Index
	sampleRate        int

	offset int32

	plotDrawEndEventArgs *util.PlotDrawEndEventArgs
}

// 싱글톤
func GetWaveformPlot() *WaveformPlot {
	var wp *WaveformPlot
	waveformPlotOnce.Do(func() {
		wp = &WaveformPlot{
			title:               "Wavefrom Data",
			sctx:                util.NewSimpleContext(),
			mtx:                 &deadlock.RWMutex{},
			isShouldDataRefresh: true,
			isCleard:            true,
			sampleArray:         util.NewWaveformPlotData(0, 0),
			sampleCutIndex:      util.NewIndex(0, 0),
			sampleCutIndexOld:   util.NewIndex(0, 0),
			maxSampleCount:      DefaultMaxSampleCount,
		}

		wp.sampleArrayView = util.NewWaveformPlotData(0, wp.maxSampleCount)
		if audio.IsAudioLoaded() {
			wp.sampleRate = int(audio.StreamFormat().SampleRate)
		}

		logOption := log.NewLoggerOption()
		logOption.Prefix = "[waveform]"
		wp.logger = plotLogger.NewSimpleLogger(logOption)

		wp.logger.Trace("Waveform Plot init...")

		wp.eventHandler_AudioStreamChanged()
	})

	waveformPlotInstance = wp
	return waveformPlotInstance
}

func (wp *WaveformPlot) PlotSetup(args imguiw.PlotSetupArgs) {
	var axisX1Flags imgui.PlotAxisFlags
	if args.IsLastSubPlot {
		axisX1Flags = imgui.PlotAxisFlagsNoLabel
	} else {
		axisX1Flags = imgui.PlotAxisFlagsNoLabel | imgui.PlotAxisFlagsNoTickLabels
	}

	imgui.PlotSetupAxisV(
		imgui.AxisX1,
		"PlotX",
		imgui.PlotAxisFlags(axisX1Flags),
	)
	imgui.PlotSetupAxisV(
		imgui.AxisY1,
		"PlotY",
		imgui.PlotAxisFlags(imgui.PlotAxisFlagsLock|
			imgui.PlotAxisFlagsNoTickLabels|
			imgui.PlotAxisFlagsNoGridLines|
			imgui.PlotAxisFlagsNoLabel|
			imgui.PlotAxisFlagsOpposite),
	)

	imgui.PlotSetupAxisLimitsV(imgui.AxisY1, -0.5, 0.5, imgui.PlotCondAlways)

	imgui.PlotSetupAxisLimitsConstraints(imgui.AxisX1, 0, args.AxisXLimitMax)

	if args.IsFitRequest {
		imgui.PlotSetupAxisLimitsV(imgui.AxisX1, 0, args.AxisXLimitMax, imgui.PlotCondAlways)
	}

	// if args.IsFitAudioStreamPos && audio.IsRunning() {
	// 	halfRange := (wp.plotDrawEndEventArgs.PlotPointEnd - wp.plotDrawEndEventArgs.PlotPointStart) / 2
	// 	pos := audio.Position().Seconds()
	// 	minLimit := pos - halfRange
	// 	maxLimit := pos + halfRange

	// 	imgui.PlotSetupAxisLimitsV(imgui.AxisX1, minLimit, maxLimit, imgui.PlotCondAlways)
	// } else {
	// 	imgui.PlotSetupAxisLimitsConstraints(imgui.AxisX1, 0, args.AxisXLimitMax)

	// 	if args.IsFitRequest {
	// 		imgui.PlotSetupAxisLimitsV(imgui.AxisX1, 0, args.AxisXLimitMax, imgui.PlotCondAlways)
	// 	}
	// }
}

func (wp *WaveformPlot) Plot() {
	if wp.isShouldDataRefresh {
		return
	}

	wp.mtx.RLock()
	sampleArrayView := *wp.sampleArrayView
	title := wp.title
	offset := wp.offset
	plotDrawEndEventArgs := wp.plotDrawEndEventArgs
	wp.mtx.RUnlock()

	dataLength := sampleArrayView.LengthY()

	if dataLength > 0 {
		imgui.PlotPlotLinedoublePtrdoublePtrV(
			title,
			&sampleArrayView.X,
			&sampleArrayView.Y,
			int32(dataLength),
			imgui.PlotLineFlagsNone,
			offset,
			8,
		)
	}

	if plotDrawEndEventArgs != nil {
		wp.mtx.Lock()
		sampleRate := wp.sampleRate
		realSampleCount := wp.sampleArray.LengthY()
		wp.sampleCutIndex.Start = max(0, int(plotDrawEndEventArgs.PlotPointStart*float64(sampleRate)))
		wp.sampleCutIndex.End = min(realSampleCount, int(plotDrawEndEventArgs.PlotPointEnd*float64(sampleRate)))
		wp.mtx.Unlock()
	}
}

// 현재 오디오 스트림에서 데이터 업데이트
func (wp *WaveformPlot) UpdateData() {
	if wp.isShouldDataRefresh {
		if audio.IsAudioLoaded() {
			wp.mtx.Lock()

			format := audio.StreamFormat()
			sampleArray := audio.GetMonoAllSampleData()
			wp.sampleRate = int(format.SampleRate)
			sampleCount := len(sampleArray)

			wp.sampleArray.Y = sampleArray
			sampleArrayX := make([]float64, sampleCount)
			for i := 0; i < sampleCount; i++ {
				sampleArrayX[i] = float64(i) / float64(wp.sampleRate)
			}
			wp.sampleArray.X = sampleArrayX

			wp.sampleCutIndex.Start = 0
			wp.sampleCutIndex.End = sampleCount

			wp.isCleard = false
			wp.isShouldDataRefresh = false

			wp.mtx.Unlock()
		} else {
			wp.clear()
		}
	} else {
		wp.updateViewData()
	}
}

// Plot에 표시되는 데이터 처리 (화면 밖 데이터 Cut, 다운샘플링)
func (wp *WaveformPlot) updateViewData() {
	if wp.isShouldDataRefresh {
		return
	}

	wp.mtx.Lock()
	sampleCutIndex := *wp.sampleCutIndex
	sampleCutIndexOld := *wp.sampleCutIndexOld
	wp.sampleCutIndexOld = &sampleCutIndex
	wp.mtx.Unlock()

	if !sampleCutIndex.Equal(sampleCutIndexOld) {
		wp.sctx.Cancel()
		sctx := util.NewSimpleContext()
		wp.sctx = sctx

		wp.mtx.Lock()

		if sctx.IsCancelled() {
			wp.mtx.Unlock()
			return
		}

		maxSampleCount := wp.maxSampleCount
		realSampleCount := wp.sampleArray.LengthY()
		viewSampleCount := sampleCutIndex.Size()

		if realSampleCount < sampleCutIndex.End {
			sampleCutIndex.End = realSampleCount
		}
		if sampleCutIndex.End <= sampleCutIndex.Start {
			sampleCutIndex.Start = sampleCutIndex.End - 1
		}

		sIdx := max(0, sampleCutIndex.Start-wp.sampleRate/3)
		eIdx := min(realSampleCount, sampleCutIndex.End+wp.sampleRate/3)

		simpleCut := func(sIdx, eIdx int) {
			sampleX := wp.sampleArray.X[sIdx:eIdx]
			sampleY := wp.sampleArray.Y[sIdx:eIdx]
			// copy(wp.sampleArrayView.X, sampleX)
			// copy(wp.sampleArrayView.Y, sampleY)
			wp.sampleArrayView.X = sampleX
			wp.sampleArrayView.Y = sampleY
		}

		if realSampleCount > maxSampleCount && viewSampleCount > maxSampleCount {
			// 다운샘플링
			sampleArrayViewX, sampleArrayViewY, err := dsp.LTTB(
				sctx.Context(),
				wp.sampleArray.X[sIdx:eIdx],
				wp.sampleArray.Y[sIdx:eIdx],
				maxSampleCount,
			)
			if err != nil {
				wp.logger.Errorf("다운샘플링 오류 (err=%v)", err)
				simpleCut(sIdx, eIdx)
			} else {
				// LOG
				// wp.logger.Tracef("다운샘플링 적용 (srcLen=(%v,%v), dstLen=(%v,%v), sIdx:eIdx=%v:%v, realSampleCount=%v, viewSampleCount=%v, maxSampleCount=%v)",
				// 	len(wp.sampleArray.X),
				// 	len(wp.sampleArray.Y),
				// 	len(sampleArrayViewX),
				// 	len(sampleArrayViewY),
				// 	sIdx, eIdx,
				// 	realSampleCount,
				// 	viewSampleCount,
				// 	maxSampleCount,
				// )
				wp.sampleArrayView.X = sampleArrayViewX
				wp.sampleArrayView.Y = sampleArrayViewY
			}
		} else {
			simpleCut(sIdx, eIdx)
		}

		wp.mtx.Unlock()
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
	if wp.isCleard {
		return
	}

	wp.sampleArray.Clear()
	wp.sampleArrayView.Clear()
	wp.sampleCutIndex = util.NewIndex(0, 0)
	wp.sampleCutIndexOld = util.NewIndex(0, 0)
	wp.isCleard = true
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
