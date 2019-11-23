/*
Simulator module simulate the camera and Pi zero, and output images periodically
*/

package server

import (
	"math"
	"math/rand"
	"time"
)

/*
RandNum outputs a random sample subject to Guassian Distribution
*/
func RandNum(mean float64, stdev float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.NormFloat64()*stdev + mean
}

/*
ImageCache outputs a random number of images in a certain period.
*/
func ImageCache() int64 {
	mean := 42.75
	stdev := 26.5
	return int64(math.Ceil(math.Abs(RandNum(mean, stdev))))
}
