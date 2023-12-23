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
	spectrogramPlotOnce     sync.Once
	spectrogramPlotInstance *SpectrogramPlot
)

type SpectrogramPlot struct {
	imguiw.PlotWidget

	logger *log.Logger

	title string

	isShouldDataRefresh bool
	isCleard            bool

	spectrogram  *dsp.Spectrogram
	freqArray    []float64
	sampleRate   int // sampling rate
	n_bin        int // spectrogram bin count
	ferq_min     float64
	ferq_max     float64
	maxWidthSize int

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
			isShouldDataRefresh: true,
			spectrogram:         dsp.NewSpectrogram(),
			sampleCutIndex:      util.NewIndex(0, 0),
			sampleCutIndexOld:   util.NewIndex(0, 0),
			ferq_min:            0,
			ferq_max:            1,
			maxWidthSize:        1024,
		}

		logOption := log.NewLoggerOption()
		logOption.Prefix = "[spectrogram]"
		sp.logger = plotLogger.NewSimpleLogger(logOption)

		sp.logger.Trace("Spectrogram Plot init...")

		sp.eventHandler_AudioStreamChanged()
	})

	spectrogramPlotInstance = sp
	return spectrogramPlotInstance
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
		if audio.IsAudioLoaded() {
			sp.sampleRate = int(audio.StreamFormat().SampleRate)
			sp.sampleArray = audio.GetMonoAllSampleData()
			sp.spectrogram.SetSampleData(sp.sampleArray, sp.sampleRate)
			sp.spectrogram.Coefficients() // caching

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

		freqs, width, height := sp.spectrogram.Coefficients()
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
