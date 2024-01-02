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
	"gonum.org/v1/gonum/dsp/window"
)

var (
	spectrogramPlotOnce     sync.Once
	spectrogramPlotInstance *SpectrogramPlot
)

type SpectrogramPlot struct {
	imguiw.PlotWidget

	logger *log.Logger

	title string

	sctx                *util.SimpleContext
	mtx                 *deadlock.RWMutex
	isShouldDataRefresh bool
	isCleard            bool

	freqArray     []float32
	sampleRate    int // sampling rate
	n_bin         int // spectrogram bin count
	ferq_min      float64
	ferq_max      float64
	maxHeightSize int
	maxWidthSize  int

	sampleArray       []float64
	sampleCutIndex    *util.Index
	sampleCutIndexOld *util.Index

	plotDrawEndEventArgs *util.PlotDrawEndEventArgs

	// temp state
	wheelEventBufferIdx int
	wheelEventBuffer    []bool
}

// 싱글톤
func GetSpectrogramPlot() *SpectrogramPlot {
	var sp *SpectrogramPlot
	spectrogramPlotOnce.Do(func() {
		sp = &SpectrogramPlot{
			title:               "Spectrogram Data",
			sctx:                util.NewSimpleContext(),
			mtx:                 &deadlock.RWMutex{},
			isShouldDataRefresh: true,
			sampleCutIndex:      util.NewIndex(0, 0),
			sampleCutIndexOld:   util.NewIndex(0, 0),
			n_bin:               512,
			ferq_min:            0,
			ferq_max:            1,
			maxWidthSize:        512,

			wheelEventBuffer: make([]bool, 3),
		}

		sp.maxHeightSize = sp.n_bin

		logOption := log.NewLoggerOption()
		logOption.Prefix = "[spectrogram]"
		sp.logger = plotLogger.NewSimpleLogger(logOption)

		sp.logger.Trace("Spectrogram Plot init...")

		sp.eventHandler_AudioStreamChanged()
	})

	spectrogramPlotInstance = sp
	return spectrogramPlotInstance
}

func (sp *SpectrogramPlot) PlotSetup(args imguiw.PlotSetupArgs) {
	sp.mtx.RLock()
	sampleRate := sp.sampleRate
	sp.mtx.RUnlock()

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
		imgui.PlotAxisFlags(imgui.PlotAxisFlagsNoTickLabels|
			imgui.PlotAxisFlagsNoGridLines|
			imgui.PlotAxisFlagsNoLabel|
			imgui.PlotAxisFlagsOpposite),
	)

	imgui.PlotSetupAxisScalePlotScale(imgui.AxisY1, imgui.PlotScaleMel)

	imgui.PlotSetupAxisLimitsConstraints(imgui.AxisY1, math.SmallestNonzeroFloat64, float64(sampleRate)/2)
	imgui.PlotSetupAxisLimitsConstraints(imgui.AxisX1, 0, args.AxisXLimitMax)

	if args.IsFitRequest {
		imgui.PlotSetupAxisLimitsV(imgui.AxisX1, 0, args.AxisXLimitMax, imgui.PlotCondAlways)
		imgui.PlotSetupAxisLimitsV(imgui.AxisY1, math.SmallestNonzeroFloat64, float64(sampleRate)/2, imgui.PlotCondAlways)
	}

	// if args.IsFitAudioStreamPos && audio.IsRunning() {
	// 	halfRange := (sp.plotDrawEndEventArgs.PlotPointEnd - sp.plotDrawEndEventArgs.PlotPointStart) / 2
	// 	pos := audio.Position().Seconds()
	// 	minLimit := pos - halfRange
	// 	maxLimit := pos + halfRange

	// 	imgui.PlotSetupAxisLimitsV(imgui.AxisX1, minLimit, maxLimit, imgui.PlotCondAlways)
	// } else {
	// 	imgui.PlotSetupAxisLimitsConstraints(imgui.AxisX1, 0, args.AxisXLimitMax)

	// 	if args.IsFitRequest {
	// 		imgui.PlotSetupAxisLimitsV(imgui.AxisX1, 0, args.AxisXLimitMax, imgui.PlotCondAlways)
	// 		imgui.PlotSetupAxisLimitsV(imgui.AxisY1, math.SmallestNonzeroFloat64, float64(sampleRate)/2, imgui.PlotCondAlways)
	// 	}
	// }
}

func (sp *SpectrogramPlot) Plot() {
	if sp.isShouldDataRefresh {
		return
	}

	sp.mtx.RLock()
	plotDrawEndEventArgs := sp.plotDrawEndEventArgs
	title := sp.title
	freqArray := sp.freqArray
	maxHeightSize := sp.maxHeightSize
	maxWidthSize := sp.maxWidthSize
	ferq_min := sp.ferq_min
	ferq_max := sp.ferq_max
	sampleRate := sp.sampleRate
	sampleArray := sp.sampleArray

	wheelEventBufferIdx := sp.wheelEventBufferIdx
	wheelEventBufferIdx++
	if wheelEventBufferIdx >= len(sp.wheelEventBuffer) {
		wheelEventBufferIdx = 0
	}
	sp.wheelEventBuffer[wheelEventBufferIdx] = imguiw.Context.IO().MouseWheel() != 0
	isWheelEvent := false
	for _, c := range sp.wheelEventBuffer {
		if c {
			isWheelEvent = true
			break
		}
	}
	sp.wheelEventBufferIdx = wheelEventBufferIdx
	sp.mtx.RUnlock()

	if plotDrawEndEventArgs != nil {
		if len(freqArray) != 0 && !isWheelEvent {
			imgui.PlotPushColormapPlotColormap(imgui.PlotColormapViridis)
			// t := time.Now()
			imgui.PlotPlotHeatmapFloatPtrV(
				title,
				freqArray,
				int32(maxHeightSize),
				int32(maxWidthSize),
				ferq_min,
				ferq_max,
				"",
				imgui.NewPlotPoint(plotDrawEndEventArgs.PlotPointStart, float64(sampleRate)/2),
				imgui.NewPlotPoint(plotDrawEndEventArgs.PlotPointEnd, 0),
				imgui.PlotHeatmapFlagsNone,
			)
			// sp.logger.Tracef("diff_t: %v", time.Since(t))
			imgui.PlotPopColormap()
		}

		sp.mtx.Lock()
		sp.sampleCutIndex.Start = max(0, int(plotDrawEndEventArgs.PlotPointStart*float64(sampleRate)))
		sp.sampleCutIndex.End = min(len(sampleArray), int(plotDrawEndEventArgs.PlotPointEnd*float64(sampleRate)))
		sp.mtx.Unlock()
	}
}

// 현재 오디오 스트림에서 데이터 업데이트
func (sp *SpectrogramPlot) UpdateData() {
	if sp.isShouldDataRefresh {
		if audio.IsAudioLoaded() {
			sp.mtx.Lock()

			sp.sampleRate = int(audio.StreamFormat().SampleRate)
			sp.sampleArray = audio.GetMonoAllSampleData()

			freqArraySize := sp.maxWidthSize * sp.maxHeightSize
			if len(sp.freqArray) != freqArraySize {
				sp.freqArray = make([]float32, freqArraySize)
			}

			sp.isCleard = false
			sp.isShouldDataRefresh = false

			sp.mtx.Unlock()
		} else {
			sp.clear()
		}
	} else {
		sp.updateViewData()
	}
}

func (sp *SpectrogramPlot) updateViewData() {
	if sp.isShouldDataRefresh {
		return
	}

	sp.mtx.Lock()
	sampleCutIndex := *sp.sampleCutIndex
	sampleCutIndexOld := *sp.sampleCutIndexOld
	sp.sampleCutIndexOld = &sampleCutIndex
	sp.mtx.Unlock()

	if !sampleCutIndex.Equal(sampleCutIndexOld) {
		sp.sctx.Cancel()
		sctx := util.NewSimpleContext()
		sp.sctx = sctx

		if sctx.IsCancelled() {
			return
		}

		sp.mtx.RLock()
		realSampleCount := len(sp.sampleArray)
		sampleArray := sp.sampleArray
		sampleRate := sp.sampleRate
		maxHeightSize := sp.maxHeightSize
		maxWidthSize := sp.maxWidthSize
		n_bin := sp.n_bin
		sp.mtx.RUnlock()

		if realSampleCount < sampleCutIndex.End {
			sampleCutIndex.End = realSampleCount
		}
		if sampleCutIndex.End <= sampleCutIndex.Start {
			sampleCutIndex.Start = sampleCutIndex.End - 1
		}

		sIdx := max(0, sampleCutIndex.Start-n_bin)
		eIdx := min(realSampleCount, sampleCutIndex.End+n_bin)

		freqArray := dsp.ParallelFFT(
			sctx.Context(),
			sampleArray[sIdx:eIdx],
			sampleRate,
			maxWidthSize,
			n_bin,
			window.BartlettHann,
			false,
		)

		if sctx.IsCancelled() {
			return
		}

		sp.mtx.Lock()

		for y := 0; y < maxHeightSize; y++ {
			for x := 0; x < maxWidthSize; x++ {
				sp.freqArray[y*maxWidthSize+x] = float32(freqArray[y*maxWidthSize+x])
			}
		}

		sp.mtx.Unlock()
	}
}

func (sp *SpectrogramPlot) Title() string {
	return sp.title
}

func (sp *SpectrogramPlot) EventHandler(eventArgs any) {
	switch castEventArgs := eventArgs.(type) {
	case util.PlotDrawEndEventArgs:
		sp.plotDrawEndEventArgs = &castEventArgs
	}
}

func (sp *SpectrogramPlot) clear() {
	sp.mtx.Lock()
	defer sp.mtx.Unlock()

	if sp.isCleard {
		return
	}

	clear(sp.sampleArray)
	sp.sampleCutIndex = util.NewIndex(0, 0)
	sp.sampleCutIndexOld = util.NewIndex(0, 0)
	sp.isCleard = true
}

func (sp *SpectrogramPlot) IsDisposed() bool {
	return false
}

// 오디오 스트림의 이벤트 수신
func (sp *SpectrogramPlot) eventHandler_AudioStreamChanged() {
	go func() {
		msgChan := audio.AudioStreamBroker().Subscribe()

		for msg := range msgChan {
			sp.logger.Tracef("[Callback] AudioStreamChanged (msg=%v)", msg)
			switch msg {
			case audio.EnumAudioStreamOpen:
				sp.isShouldDataRefresh = true
			case audio.EnumAudioStreamClosed:
				sp.isShouldDataRefresh = true
			}
		}
	}()
}
