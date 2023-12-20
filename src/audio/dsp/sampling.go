package dsp

import (
	"fmt"
	"math"
)

// LTTB down-samples the data to contain only threshold number of points that
// have the same visual shape as the original data
func LTTB(x []float64, y []float64, threshold int) ([]float64, []float64, error) {
	if threshold < 3 {
		threshold = 3
	}

	if len(x) != len(y) {
		return nil, nil, fmt.Errorf("LTTB data mismatch Error! len(x) != len(y): x length is %v, y length is %v", len(x), len(y))
	}

	if threshold >= len(y) {
		return x, y, nil // Nothing to do
	}

	sampledX, sampledY, err := LTTB_Buffer(x, y, nil, nil, threshold)

	return sampledX, sampledY, err
}

// LTTB down-samples the data to contain only threshold number of points that
// have the same visual shape as the original data
// (With Output Buffer)
func LTTB_Buffer(x, y, outx, outy []float64, threshold int) ([]float64, []float64, error) {
	if threshold < 3 {
		threshold = 3
	}

	if len(x) != len(y) {
		return nil, nil, fmt.Errorf("LTTB data mismatch Error! len(x) != len(y): x length is %v, y length is %v", len(x), len(y))
	}

	if threshold >= len(y) {
		return x, y, nil // Nothing to do
	}

	var sampledX []float64
	if cap(outx) < threshold {
		sampledX = make([]float64, 0, threshold)
	} else {
		sampledX = outx[0:0]
	}

	var sampledY []float64
	if cap(outy) < threshold {
		sampledY = make([]float64, 0, threshold)
	} else {
		sampledY = outy[0:0]
	}

	// Bucket size. Leave room for start and end data points
	every := float64(len(y)-2) / float64(threshold-2)

	sampledX = append(sampledX, x[0]) // Always add the first point
	sampledY = append(sampledY, y[0]) // Always add the first point

	bucketStart := 1
	bucketCenter := int(math.Floor(every)) + 1

	var a int

	for i := 0; i < threshold-2; i++ {

		bucketEnd := int(math.Floor(float64(i+2)*every)) + 1

		// Calculate point average for next bucket (containing c)
		avgRangeStart := bucketCenter
		avgRangeEnd := bucketEnd

		if avgRangeEnd >= len(y) {
			avgRangeEnd = len(y)
		}

		avgRangeLength := float64(avgRangeEnd - avgRangeStart)

		var avgX, avgY float64
		for ; avgRangeStart < avgRangeEnd; avgRangeStart++ {
			avgX += x[avgRangeStart]
			avgY += y[avgRangeStart]
		}
		avgX /= avgRangeLength
		avgY /= avgRangeLength

		// Get the range for this bucket
		rangeOffs := bucketStart
		rangeTo := bucketCenter

		// Point a
		pointAX := x[a]
		pointAY := y[a]

		maxArea := float64(-1.0)

		var nextA int
		for ; rangeOffs < rangeTo; rangeOffs++ {
			// Calculate triangle area over three buckets
			area := (pointAX-avgX)*(y[rangeOffs]-pointAY) - (pointAX-x[rangeOffs])*(avgY-pointAY)
			// We only care about the relative area here.
			// Calling math.Abs() is slower than squaring
			area *= area
			if area > maxArea {
				maxArea = area
				nextA = rangeOffs // Next a is this b
			}
		}

		// Pick this point from the bucket
		sampledX = append(sampledX, x[nextA])
		sampledY = append(sampledY, y[nextA])

		a = nextA // This a is the next a (chosen b)

		bucketStart = bucketCenter
		bucketCenter = bucketEnd
	}

	// Always add last
	sampledX = append(sampledX, x[len(x)-1])
	sampledY = append(sampledY, y[len(y)-1])

	return sampledX, sampledY, nil
}
