package widget

import (
	"math/cmplx"
	"sync"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/Kor-SVS/cocoa/src/audio"
	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/Kor-SVS/cocoa/src/ui/imguiw"
	"github.com/Kor-SVS/cocoa/src/util"
	"gonum.org/v1/gonum/dsp/fourier"
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

	isShouldDataRefresh bool
	isCleard            bool

	fft *fourier.FFT

	freqArray  []float64
	sampleRate int // sampling rate
	n_freq     int // FFT frequency count
	n_bin      int // spectrogram bin count
	ferq_min   float64
	ferq_max   float64

	sampleArray       []float64
	sampleCutIndex    *util.Index
	sampleCutIndexOld *util.Index

	plotDrawEndEventArgs *util.PlotDrawEndEventArgs
}

// 싱글톤
func GetSpectrogramPlot() *SpectrogramPlot {
	spectrogramPlotOnce.Do(func() {
		spectrogramPlotInstance = &SpectrogramPlot{
			title:               "Spectrogram Data",
			isShouldDataRefresh: true,
			freqArray:           make([]float64, 0),
			sampleArray:         make([]float64, 0),
			sampleCutIndex:      util.NewIndex(0, 0),
			sampleCutIndexOld:   util.NewIndex(0, 0),
			ferq_min:            0,
			ferq_max:            1,
		}

		logOption := log.NewLoggerOption()
		logOption.Prefix = "[spectrogram]"
		spectrogramPlotInstance.logger = plotLogger.NewSimpleLogger(logOption)

		spectrogramPlotInstance.logger.Trace("Spectrogram Plot init...")

		spectrogramPlotInstance.eventHandler_AudioStreamChanged()
	})

	return spectrogramPlotInstance
}

func (sp *SpectrogramPlot) Plot() {
	if sp.isShouldDataRefresh {
		return
	}

	if len(sp.freqArray) != 0 {
		imgui.PlotPushColormapPlotColormap(imgui.PlotColormapViridis)
		imgui.PlotPlotHeatmapdoublePtrV(
			sp.title,
			&sp.freqArray,
			int32(sp.n_bin),
			int32(sp.n_freq),
			sp.ferq_min,
			sp.ferq_max,
			"",
			imgui.NewPlotPoint(sp.plotDrawEndEventArgs.PlotPointStart, 0.5),
			imgui.NewPlotPoint(sp.plotDrawEndEventArgs.PlotPointEnd, -0.5),
			imgui.PlotHeatmapFlagsNone,
		)
		imgui.PlotPopColormap()
	}

	if sp.plotDrawEndEventArgs != nil {
		plotDrawEndEventArgs := *sp.plotDrawEndEventArgs
		sp.sampleCutIndex.Start = max(0, int(plotDrawEndEventArgs.PlotPointStart*float64(sp.sampleRate)))
		sp.sampleCutIndex.End = min(len(sp.sampleArray), int(plotDrawEndEventArgs.PlotPointEnd*float64(sp.sampleRate)))
	}
}

func (sp *SpectrogramPlot) fftInit(n int) {
	if sp.fft == nil {
		sp.fft = fourier.NewFFT(n * 2)
	} else {
		sp.fft.Reset(n * 2)
	}

	sp.sampleRate = int(audio.StreamFormat().SampleRate)
	sp.n_bin = n

	spectrogramPlotInstance.logger.Infof("FFT init (n_fft: %v)", sp.fft.Len())
}

// 현재 오디오 스트림에서 데이터 불러오기
func (sp *SpectrogramPlot) UpdateData() {
	if sp.isShouldDataRefresh {
		if audio.IsAudioLoaded() {
			sp.fftInit(512) // n_fft param
			sp.sampleArray = audio.GetMonoAllSampleData()

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

		samples := sp.sampleArray[sampleCutIndex.Start:sampleCutIndex.End]

		heatmapWidth := 768 // resolution param
		freqArraySize := heatmapWidth * (sp.n_bin * 2)
		if cap(sp.freqArray) < freqArraySize {
			sp.freqArray = make([]float64, freqArraySize)
		} else {
			sp.freqArray = sp.freqArray[:freqArraySize]
		}
		sp.n_freq = heatmapWidth

		sub := make([]float64, sp.n_bin*2)
		freqs := make([]complex128, len(sub)/2+1)
		for x := 1; x < heatmapWidth; x++ {
			n0 := int64(util.MapRange(float64(x-1), 0, float64(heatmapWidth), 0, float64(len(samples))))
			n1 := int64(util.MapRange(float64(x-0), 0, float64(heatmapWidth), 0, float64(len(samples))))
			c := n0 + (n1-n0)/2

			for i := 0; i < len(sub); i += 1 {
				s := 0.0
				n := int(c) - sp.n_bin + int(i)
				if n >= 0 && n < len(samples) {
					s = samples[n]
				}
				sub[i] = s
			}

			sub = window.Blackman(sub)

			freqs := sp.fft.Coefficients(freqs, sub)

			for y := 0; y < len(freqs); y++ {
				sp.freqArray[y*heatmapWidth+x] = cmplx.Abs(freqs[y])
			}
		}

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
