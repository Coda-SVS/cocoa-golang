package dsp

// Source: https://github.com/xigh/spectrogram

import (
	"context"
	"math"
	"math/cmplx"
	"sync"

	"github.com/Kor-SVS/cocoa/src/util"
)

// FFT N Bin -> Freq
func NBin2Freq(i, sampleRate, n_bin int) float64 {
	return float64(i) * (float64(sampleRate) / float64(n_bin) / 2.0)
}

// Freq -> FFT N Bin
func Freq2NBin(freq float64, sampleRate, n_bin int) int {
	return int(freq / (float64(sampleRate) / float64(n_bin) / 2.0))
}

func Mel(freq float64) float64 {
	return 2595 * math.Log10(1+freq/700)
}

func dft(samples []float64, freqs []complex128) {
	// freqs := make([]complex128, len(samples))

	arg := -2.0 * math.Pi / float64(len(samples))
	for k := 0; k < len(samples); k++ {
		r, i := 0.0, 0.0
		for n := 0; n < len(samples); n++ {
			r += samples[n] * math.Cos(arg*float64(n)*float64(k))
			i += samples[n] * math.Sin(arg*float64(n)*float64(k))
		}
		freqs[k] = complex(r, i)
	}
}

func hfft(samples []float64, freqs []complex128, n, step int) {
	if n == 1 {
		freqs[0] = complex(samples[0], 0)
		return
	}

	half := n / 2

	hfft(samples, freqs, half, 2*step)
	hfft(samples[step:], freqs[half:], half, 2*step)

	for k := 0; k < half; k++ {
		a := -2 * math.Pi * float64(k) / float64(n)
		e := cmplx.Rect(1, a) * freqs[k+half]

		freqs[k], freqs[k+half] = freqs[k]+e, freqs[k]-e
	}
}

func FFT(
	samples []float64,
	freqArrayWidth int,
	n_bin int,
	windowFn func([]float64) []float64,
	isDFT bool,
) (freqArray []float64) {
	freqArraySize := freqArrayWidth * n_bin
	freqArray = make([]float64, freqArraySize)

	sampleArrayLength := len(samples)

	for x := 1; x < freqArrayWidth; x++ {
		n0 := int64(util.MapRange(float64(x-1), 0, float64(freqArrayWidth), 0, float64(sampleArrayLength)))
		n1 := int64(util.MapRange(float64(x), 0, float64(freqArrayWidth), 0, float64(sampleArrayLength)))
		c := n0 + (n1-n0)/2

		subSampleArray := make([]float64, n_bin*2)
		for i := 0; i < len(subSampleArray); i++ {
			s := 0.0
			n := int(c) - n_bin + int(i)
			if n >= 0 && n < sampleArrayLength {
				s = samples[n]
			}
			subSampleArray[i] = s
		}

		subSampleArray = windowFn(subSampleArray)

		freqs := make([]complex128, n_bin*2)
		if isDFT {
			dft(subSampleArray, freqs)
		} else {
			hfft(subSampleArray, freqs, n_bin*2, 1)
		}

		for y := 0; y < n_bin; y++ {
			freqArray[y*freqArrayWidth+(x-1)] = cmplx.Abs(freqs[y])
		}
	}

	return freqArray
}

func ParallelFFT(
	ctx context.Context,
	samples []float64,
	sampleRate int,
	freqArrayWidth int,
	n_bin int,
	windowFn func([]float64) []float64,
	isDFT bool,
) (freqArray []float64) {
	freqArraySize := freqArrayWidth * n_bin
	freqArray = make([]float64, freqArraySize)

	isCancelled := func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
			return false
		}
	}

	if isCancelled() {
		return nil
	}

	subSampleArrayPool := &sync.Pool{
		New: func() interface{} {
			array := make([]float64, n_bin*2)
			return &array
		},
	}

	freqArrayPool := &sync.Pool{
		New: func() interface{} {
			array := make([]complex128, n_bin*2)
			return &array
		},
	}

	wg := &sync.WaitGroup{}
	sampleLength := len(samples)

	for x := 1; x < freqArrayWidth; x++ {
		if isCancelled() {
			return nil
		}

		wg.Add(1)

		xIdx := x

		goPool.Submit(func() {
			defer wg.Done()

			if isCancelled() {
				return
			}

			n0 := int64(util.MapRange(float64(xIdx-1), 0, float64(freqArrayWidth), 0, float64(sampleLength)))
			n1 := int64(util.MapRange(float64(xIdx), 0, float64(freqArrayWidth), 0, float64(sampleLength)))
			c := int(n0 + (n1-n0)/2)

			subSampleArray := *subSampleArrayPool.Get().(*[]float64)
			for i := 0; i < len(subSampleArray); i++ {
				s := 0.0
				n := c - n_bin + int(i)
				if n >= 0 && n < sampleLength {
					s = samples[n]
				}
				subSampleArray[i] = s
			}

			subSampleArray = windowFn(subSampleArray)

			if isCancelled() {
				subSampleArrayPool.Put(&subSampleArray)
				return
			}

			freqs := *freqArrayPool.Get().(*[]complex128)
			if isDFT {
				dft(subSampleArray, freqs)
			} else {
				hfft(subSampleArray, freqs, n_bin*2, 1)
			}

			for y := 0; y < n_bin; y++ {
				srcIdx := y
				dstIdx := y*freqArrayWidth + (xIdx - 1)
				freqArray[dstIdx] = cmplx.Abs(freqs[srcIdx])
			}

			subSampleArrayPool.Put(&subSampleArray)
			freqArrayPool.Put(&freqs)
		})
	}

	if isCancelled() {
		return nil
	}

	wg.Wait()

	return freqArray
}
