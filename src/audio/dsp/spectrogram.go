package dsp

import (
	"github.com/sasha-s/go-deadlock"
	"gonum.org/v1/gonum/dsp/window"
)

type Spectrogram struct {
	mtx *deadlock.Mutex

	windowFn func([]float64) []float64

	freqArray      []float64
	freqArrayWidth int
	sampleRate     int // sampling rate
	n_bin          int // spectrogram bin count

	sampleArray []float64
	isCached    bool

	isDFT bool
}

func NewSpectrogram() *Spectrogram {
	st := &Spectrogram{
		mtx:         &deadlock.Mutex{},
		windowFn:    window.Blackman,
		freqArray:   nil,
		sampleArray: nil,
		n_bin:       1024,
		isDFT:       false,
	}
	return st
}

func (st *Spectrogram) SetWindowFunc(windowFn func([]float64) []float64) {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	st.windowFn = windowFn
	st.burstCache()
}

func (st *Spectrogram) SetSampleData(sampleArray []float64, sampleRate int) {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	st.sampleArray = sampleArray
	st.sampleRate = sampleRate
	st.burstCache()
}

func (st *Spectrogram) NumBin() int {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	return st.n_bin
}

func (st *Spectrogram) SetNumBin(n_bin int) {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	if st.n_bin == n_bin {
		return
	}

	st.n_bin = n_bin
	st.burstCache()
}

func (st *Spectrogram) ResetFFTData() {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	st.burstCache()
}

func (st *Spectrogram) ResetSampleData() {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	st.isCached = false
	st.sampleArray = nil
}

func (st *Spectrogram) burstCache() {
	st.isCached = false
	st.freqArray = nil
}

func (st *Spectrogram) Coefficients() ([]float64, int, int) {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	if st.sampleArray == nil {
		return nil, 0, 0
	}

	if st.isCached {
		return st.freqArray, st.freqArrayWidth, st.n_bin
	}

	// 조정 가능
	st.freqArrayWidth = len(st.sampleArray)

	st.freqArray = ParallelFFT(
		st.sampleArray,
		st.freqArrayWidth,
		st.n_bin,
		st.windowFn,
		st.isDFT,
	)

	st.isCached = true

	return st.freqArray, st.freqArrayWidth, st.n_bin
}
