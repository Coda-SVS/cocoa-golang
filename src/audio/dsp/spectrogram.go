package dsp

import (
	"context"
	"sync/atomic"

	"github.com/sasha-s/go-deadlock"
	"gonum.org/v1/gonum/dsp/window"
)

type Spectrogram struct {
	mtx       *deadlock.Mutex
	isWorking *atomic.Bool

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

func (st *Spectrogram) Coefficients(ctx context.Context, freqArrayWidth int) ([]float64, int, int) {
	if st.sampleArray == nil {
		return nil, 0, 0
	}

	if st.isCached {
		return st.freqArray, st.freqArrayWidth, st.n_bin
	}

	st.mtx.Lock()
	defer st.mtx.Unlock()
	st.isWorking.Store(true)

	// 조정 가능
	st.freqArrayWidth = min(freqArrayWidth, len(st.sampleArray))

	st.freqArray = ParallelFFT(
		ctx,
		st.sampleArray,
		st.sampleRate,
		st.freqArrayWidth,
		st.n_bin,
		st.windowFn,
		st.isDFT,
	)

	st.isCached = true
	st.isWorking.Store(false)

	return st.freqArray, st.freqArrayWidth, st.n_bin
}

func (st *Spectrogram) Coefficients_NonBlock(ctx context.Context, freqArrayWidth int) ([]float64, int, int) {
	if st.isWorking.Load() {
		return nil, 0, 0
	}

	if !st.isCached {
		go st.Coefficients(ctx, freqArrayWidth)
		return nil, 0, 0
	}

	return st.freqArray, st.freqArrayWidth, st.n_bin
}
