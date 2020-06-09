package test

import (
	"encoding/gob"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/heronalps/STOIC/client"
	"github.com/stretchr/testify/assert"
)

func TestQueryGPUNum(t *testing.T) {
	var (
		numGPU interface{}
	)
	numGPU = client.QueryGPUNum(namespace, deployment)
	_, ok := numGPU.(int64)
	if !ok {
		t.Errorf("The query of GPU number was not successful ...\n")
	}
}

func TestSelectRunTime(t *testing.T) {
	runtime, _ := client.SelectRunTime(56, zipPath, app, version, "")
	fmt.Println("runtime : " + runtime)
}

func TestGetDeploymentTime(t *testing.T) {
	deploymentTimes := client.GetDeploymentTime(runtime)
	fmt.Printf("runtime: %v \n", deploymentTimes)
}

func TestCompareVersion(t *testing.T) {
	assert.Equal(t, client.CompareVersion("0", "1.0"), -1)
}

func TestRunOnNautilus(t *testing.T) {
	output, _, timeLog := client.RunOnNautilus(runtime, zipPath, imageNum, app, version, transferTime)
	fmt.Printf("Output : %v..\n", string(output))
	assert.NotNil(t, timeLog)
}

func TestRunOnEdge(t *testing.T) {
	output, _, timeLog := client.RunOnEdge(zipPath, imageNum, app, version)
	fmt.Printf("Output : %v..\n", string(output))
	assert.NotNil(t, timeLog)
}

func TestLogTimes(t *testing.T) {
	predTimeLog := &client.TimeLog{
		Total:      3.0,
		Transfer:   1.0,
		Deployment: 1.0,
		Processing: 1.0,
	}
	actTimeLog := &client.TimeLog{
		Total:      6.0,
		Transfer:   2.0,
		Deployment: 2.0,
		Processing: 2.0,
	}
	client.LogTimes(imageNum, app, version, runtime, predTimeLog, actTimeLog)
}

func TestMedian(t *testing.T) {
	assert.Equal(t, 0.2, client.Median([]float64{0.1, 0.2, 0.3}))
	assert.Equal(t, 0.25, client.Median([]float64{0.1, 0.2, 0.3, 0.4}))
}

func TestGetWindowSize(t *testing.T) {
	optWinSize := client.GetWindowSize(runtime)
	fmt.Println(optWinSize)
}

func TestUpdateWindowSizes(t *testing.T) {
	client.UpdateWindowSizes()
}

func TestBatchSize(t *testing.T) {
	num := client.BatchSize()
	fmt.Println(time.Now().UnixNano())
	fmt.Println(time.Now())
	fmt.Println(num)
}

func TestGenerateWorkLoad(t *testing.T) {
	workload := client.GenerateWorkLoad(6700)
	fmt.Println(workload)
}

func TestServerWorkload(t *testing.T) {
	fmt.Println(client.Workload[223])
}

func TestRegisterImages(t *testing.T) {
	// var paths [3]string
	// paths[0] = "/opt"
	// paths[1] = "/opt2"
	// paths[2] = "/opt3"
	var paths [1]string
	paths[0] = "/home/ubuntu/WTB_samples/2015-11-01"
	for _, path := range paths {
		client.RegisterImages(path)
		decodeFile, err := os.Open(path + "/registryMap.gob")
		if err != nil {
			panic(err)
		}
		defer decodeFile.Close()

		decoder := gob.NewDecoder(decodeFile)
		registryMap := make(map[int]string)

		decoder.Decode(&registryMap)
		fmt.Println(registryMap)
	}
}

func TestGenerateBatch(t *testing.T) {
	path, _ := client.GenerateBatch(0, 1)
	fmt.Println(path)
}

func TestCopyFile(t *testing.T) {
	source := "/Users/michaelzhang/Downloads/WTB_samples/time_lapse/Main_2013-07-30_10:43:56_17411_12-398.jpg"
	target := "/Users/michaelzhang/Downloads/test.jpg"
	client.CopyFile(source, target)
}

func TestGetTransferTime(t *testing.T) {
	zipPath := "/Users/michaelzhang/go/src/github.com/heronalps/STOIC/image_batch_2.zip"
	client.GetTransferTime(zipPath)
}

func TestListBucket(t *testing.T) {
	client.ListBucket()
}

func TestTransferBatch(t *testing.T) {
	path := "/Users/michaelzhang/go/src/github.com/heronalps/STOIC/image_batch_2.zip"
	client.TransferBatch(path)
}

func TestListObject(t *testing.T) {
	bucket := "test-bucket"
	client.ListObject(bucket)
}

func TestHomeDir(t *testing.T) {
	fmt.Println(client.HomeDir())
}

func TestParseImageFileName(t *testing.T) {
	fmt.Println(client.ParseImageFileName("/Users/michaelzhang/Downloads/WTB_samples/time_lapse/Main_2013-08-03_13:20:10_5192_16-407.jpg"))
}

func TestDecodeImageRegistry(t *testing.T) {
	client.DecodeImageRegistry("/Users/michaelzhang/Downloads/WTB_samples/registryMap.gob")
}

func TestBatchSizeECDF(t *testing.T) {
	fmt.Println(client.BatchSizeECDF())
}
