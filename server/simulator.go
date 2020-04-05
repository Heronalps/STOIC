/*
Simulator module simulate the camera and Pi zero, and output images periodically
*/

package server

import (
	"encoding/gob"
	"io"
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
BatchSize outputs a random number of images in a certain period.
For WTB photo repository, the average amount of photo per hour is 42.75, and stdev is 26.5
*/
func BatchSize() int {
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
		arr[i] = BatchSize()
	}
	return arr
}

/*
RegisterImages persist a map to file system to register the mapping
from picture sequence number to file name
*/
func RegisterImages(rootPath string) {
	var (
		seqNo int = 1
	)
	registryMap := make(map[int]string)
	// Exclude thumbnails (xxx_t.jpg)
	re := regexp.MustCompile(`.*[0-9][^_t]\.jpg`)
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

/*
GenerateBatch generates a batch of images randomly selected from 3 volumes /opt /opt2 /opt3
1. generate 3 random numbers, volumnNo, seqNo, batch size
2. select corresponding pictures, cache at a local dir and package them
3. Return the path string to the server socket that sends the package to edge server
*/
func GenerateBatch() string {
	var (
		registryPath string
	)
	// Select volume
	switch volumnNo := rand.Intn(3); volumnNo {
	case 0:
		registryPath = "/opt/registryMap.gob"
	case 1:
		registryPath = "/opt2/registryMap.gob"
	case 2:
		registryPath = "/opt3/registryMap.gob"
	}
	registryPath = "/Users/michaelzhang/Downloads/WTB_samples/registryMap.gob"

	// Decode registry map
	decodeFile, err := os.Open(registryPath)
	if err != nil {
		panic(err)
	}
	defer decodeFile.Close()

	decoder := gob.NewDecoder(decodeFile)
	registryMap := make(map[int]string)

	decoder.Decode(&registryMap)
	// Subtract the maximum batch size to ensure enough images following seqNo
	seqNo := rand.Intn(len(registryMap) - 200)
	batchSize := BatchSize()
	for idx := seqNo; idx < seqNo+batchSize; idx++ {

	}

	return ""
}

/*
CopyFile copies the image from path to local buffer folder
*/
func CopyFile(source string, target string) bool {
	sourceFile, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer sourceFile.Close()

	sourceCopied, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer sourceCopied.Close()

	_, err = io.Copy(sourceCopied, sourceFile)
	if err != nil {
		return false
	}
	return true
}
