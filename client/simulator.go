/*
Simulator module simulate the camera and Pi zero, and output images periodically
*/

package client

import (
	"archive/zip"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
func GenerateBatch(imageNum int, batchNo int) (string, int) {
	var (
		rootPath  string
		files     []string
		batchSize int
	)
	// Select volume
	switch volumnNo := rand.Intn(3); volumnNo {
	case 0:
		rootPath = "/opt"
	case 1:
		rootPath = "/opt2"
	case 2:
		rootPath = "/opt3"
	}
	rootPath = "/Users/michaelzhang/Downloads/WTB_samples"
	registryPath := rootPath + "/registryMap.gob"

	// Decode registry map
	decodeFile, err := os.Open(registryPath)
	if err != nil {
		panic(err)
	}
	defer decodeFile.Close()

	decoder := gob.NewDecoder(decodeFile)
	registryMap := make(map[int]string)

	decoder.Decode(&registryMap)
	fmt.Println("Done decoding registry map..")

	// Subtract the maximum batch size to ensure enough images following seqNo
	rand.Seed(time.Now().UnixNano())
	seqNo := rand.Intn(len(registryMap) - 200)

	// Override batchSize if image number is set
	if imageNum == 0 {
		batchSize = BatchSize()
	} else {
		batchSize = imageNum
	}

	// Copy to local buffer dir
	for idx := seqNo; idx < seqNo+batchSize; idx++ {
		files = append(files, registryMap[idx])
	}
	fmt.Println("Done appending file paths..")
	zipPath := rootPath + "/image_batch_" + strconv.Itoa(batchNo) + ".zip"
	// package the batch
	if err := ZipFiles(zipPath, files); err != nil {
		panic(err)
	}
	fmt.Println("Done zipping files...")
	return zipPath, batchSize
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

/*
ZipFiles compresses one or many files into one zip archive file
*/
func ZipFiles(zipName string, files []string) error {
	newZipFile, err := os.Create(zipName)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

/*
AddFileToZip adds seperate files to zip archive package
*/
func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	// header.Name = filename

	// Change to deflate to gain better compression
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
