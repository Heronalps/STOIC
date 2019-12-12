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
For WTB photo repository, the average amount of photo per hour is 42.75, and stdev is 26.5
*/
func ImageCache() int {
	mean := 42.75
	stdev := 39.5
	return int(math.Ceil(math.Abs(RandNum(mean, stdev))))
}
