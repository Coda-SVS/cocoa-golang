package plot

import (
	g "github.com/AllenDang/giu"
	"github.com/Kor-SVS/cocoa/src/audio"
	"github.com/Kor-SVS/cocoa/src/util"
	"github.com/sasha-s/go-deadlock"
)

var mutex *deadlock.RWMutex = new(deadlock.RWMutex)

var (
	SampleArray    = make([]float64, 0)
	SamplePosArray = make([]float64, 0)
)

func init() {
	onStreamChanged()
}

func WaveformGUILoop() *g.PlotCanvasWidget {
	mutex.RLock()
	defer mutex.RUnlock()

	return g.Plot("Waveform Plot").
		AxisLimits(0, audio.Duration().Seconds(), -1, 1, g.ConditionOnce).
		XAxeFlags(g.PlotAxisFlagsNone).
		YAxeFlags(g.PlotAxisFlagsLock|g.PlotAxisFlagsNoTickLabels, g.PlotAxisFlagsNone, g.PlotAxisFlagsNone).
		Plots(g.LineXY("Waveform Line", SamplePosArray, SampleArray)).
		Size(-1, -1)
}

func onStreamChanged() {
	go func() {
		msgChan := audio.AudioStreamBroker.Subscribe()

		for msg := range msgChan {
			mutex.Lock()
			switch msg {
			case audio.EnumAudioStreamOpen:
				format := audio.StreamFormat()
				SampleArray = util.StereoToMono(audio.GetAllSampleData())
				SamplePosArray = make([]float64, len(SampleArray))
				for i := 0; i < len(SampleArray); i++ {
					SamplePosArray[i] = float64(i) / float64(format.SampleRate)
				}
			case audio.EnumAudioStreamClosed:
				SampleArray = make([]float64, 0)
				SamplePosArray = make([]float64, 0)
			}
			mutex.Unlock()
		}
	}()
}
