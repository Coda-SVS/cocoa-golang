package util

import (
	"encoding/binary"
	"math"
)

func FloatSampleToByteArray(inB [][2]float64, outB []byte) {
	sampleLen := len(inB)

	inBufferIdx := 0
	for inBufferIdx < sampleLen {
		byteArrayIdx := (inBufferIdx * 2) * 4
		binary.NativeEndian.PutUint32(outB[byteArrayIdx:byteArrayIdx+4], math.Float32bits(float32(inB[inBufferIdx][0])))
		binary.NativeEndian.PutUint32(outB[byteArrayIdx+4:byteArrayIdx+8], math.Float32bits(float32(inB[inBufferIdx][1])))
		inBufferIdx++
	}
}

func MonoToStereo(in []float64) [][2]float64 {
	out := make([][2]float64, len(in))
	for i, s := range in {
		if s < -1 {
			s = -1
		} else if s > 1 {
			s = 1
		}
		out[i][0] = s
		out[i][1] = s
	}
	return out
}

func StereoToMono(in [][2]float64) []float64 {
	out := make([]float64, len(in))
	var temp float64
	for i, s := range in {
		temp = (s[0] + s[1]) / float64(2)
		if temp < -1 {
			temp = -1
		} else if temp > 1 {
			temp = 1
		}
		out[i] = temp
	}
	return out
}
