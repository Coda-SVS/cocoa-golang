package dsp

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
