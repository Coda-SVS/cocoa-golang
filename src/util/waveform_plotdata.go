package util

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
