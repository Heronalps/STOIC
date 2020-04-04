/*
Simulator module simulate the camera and Pi zero, and output images periodically
*/

package server

import (
	"encoding/gob"
	"math"
	"math/rand"
	"os"
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

/*
GenerateWorkLoad generates constant work load
*/
func GenerateWorkLoad(length int) []int {
	arr := make([]int, length)
	for i := 0; i < length; i++ {
		arr[i] = ImageCache()
	}
	return arr
}

/*
RegisterImages persist a map to file system to register the mapping
from picture sequence number to file name
*/
func RegisterImages(path string) {
	registryMap := make(map[string]string)
	registryMap["abc"] = "def"

	registryFile, err := os.Create(path + "/registryMap.gob")
	if err != nil {
		panic(err)
	}

	encoder := gob.NewEncoder(registryFile)
	if err := encoder.Encode(registryMap); err != nil {
		panic(err)
	}
	registryFile.Close()
}
