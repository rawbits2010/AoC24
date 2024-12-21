package main

import "math"

type ZScore struct {
	count     int
	mean      float64
	variance  float64
	threshold float64
}

func NewZScore(threshold float64) ZScore {
	return ZScore{
		count:     0,
		mean:      0,
		variance:  0,
		threshold: threshold,
	}
}

func (zs *ZScore) Add(value float64) bool {

	zs.count++

	diff := value - zs.mean
	zs.mean += diff / float64(zs.count)
	zs.variance += diff * (value - zs.mean)

	if zs.count == 0 {
		return false
	}

	stdDev := math.Sqrt(zs.variance / float64(zs.count))
	zScore := (value - zs.mean) / stdDev

	return zScore > zs.threshold
}
