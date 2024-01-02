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

	if sp.plotDrawEndEventArgs != nil {
		if len(sp.freqArray) != 0 {
			imgui.PlotPushColormapPlotColormap(imgui.PlotColormapViridis)
			imgui.PlotPlotHeatmapdoublePtrV(
				sp.title,
				&sp.freqArray,
				int32(sp.n_bin),
				int32(sp.maxWidthSize),
				sp.ferq_min,
				sp.ferq_max,
				"",
				imgui.NewPlotPoint(sp.plotDrawEndEventArgs.PlotPointStart, 0.5),
				imgui.NewPlotPoint(sp.plotDrawEndEventArgs.PlotPointEnd, -0.5),
				imgui.PlotHeatmapFlagsNone,
			)
			imgui.PlotPopColormap()
		}

		plotDrawEndEventArgs := *sp.plotDrawEndEventArgs
		sp.sampleCutIndex.Start = max(0, int(plotDrawEndEventArgs.PlotPointStart*float64(sp.sampleRate)))
		sp.sampleCutIndex.End = min(len(sp.sampleArray), int(plotDrawEndEventArgs.PlotPointEnd*float64(sp.sampleRate)))
	}
}

// 현재 오디오 스트림에서 데이터 불러오기
func (sp *SpectrogramPlot) UpdateData() {
	if sp.isShouldDataRefresh {
		// dataContext, cancel := context.WithCancel(context.Background())

		if audio.IsAudioLoaded() {
			sp.sampleRate = int(audio.StreamFormat().SampleRate)
			sp.sampleArray = audio.GetMonoAllSampleData()
			sp.spectrogram.SetSampleData(sp.sampleArray, sp.sampleRate)
			sp.spectrogram.Coefficients(context.TODO()) // caching

			sp.n_bin = sp.spectrogram.NumBin()

			freqArraySize := sp.maxWidthSize * sp.spectrogram.NumBin()
			if len(sp.freqArray) != freqArraySize {
				sp.freqArray = make([]float64, freqArraySize)
			}

			sp.isCleard = false
			sp.isShouldDataRefresh = false
		} else {
			sp.clear()
		}
	}

	sp.updateViewData()
}

func (sp *SpectrogramPlot) updateViewData() {
	if sp.isShouldDataRefresh {
		return
	}

	sampleCutIndex := *sp.sampleCutIndex
	sampleCutIndexOld := *sp.sampleCutIndexOld

	if !sampleCutIndex.Equal(sampleCutIndexOld) {
		if len(sp.sampleArray) < sampleCutIndex.End {
			sampleCutIndex.End = len(sp.sampleArray)
		}
		if sampleCutIndex.End <= sampleCutIndex.Start {
			sampleCutIndex.Start = sampleCutIndex.End - 1
		}

		freqs, width, height := sp.spectrogram.Coefficients(context.TODO())
		viewSampleSize := min(width, sampleCutIndex.End) - max(0, sampleCutIndex.Start)
		for x := 0; x < sp.maxWidthSize; x++ {
			mx := sampleCutIndex.Start + int(util.MapRange(float64(x), 0, float64(sp.maxWidthSize), 0, float64(viewSampleSize)))
			for y := 0; y < height; y++ {
				sp.freqArray[y*sp.maxWidthSize+x] = freqs[y*width+mx]
			}
		}

		// TEST CODE
		// freqs, width, height := sp.spectrogram.Coefficients()
		// sp.freqArray = freqs
		// sp.maxWidthSize = width
		// sp.n_bin = height

		sp.sampleCutIndexOld = &sampleCutIndex
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
