/*
Simulator module simulate the camera and Pi zero, and output images periodically
*/

package server

import (
	"encoding/gob"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
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
func RegisterImages(rootPath string) {
	var (
		seqNo int32 = 1
	)
	registryMap := make(map[int32]string)
	re := regexp.MustCompile(`.*\.jpg`)
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// fmt.Printf("path : %s \n", path)
		// fmt.Printf("seqNo : %v \n", seqNo)
		if match := re.FindString(path); len(match) > 0 {
			registryMap[seqNo] = path
			seqNo++
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	registryFile, err := os.Create(rootPath + "/registryMap.gob")
	if err != nil {
		panic(err)
	}

	encoder := gob.NewEncoder(registryFile)
	if err := encoder.Encode(registryMap); err != nil {
		panic(err)
	}
	registryFile.Close()
}
