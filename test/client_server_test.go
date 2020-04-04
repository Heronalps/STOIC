package test

import (
	"encoding/gob"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/heronalps/STOIC/client"
	"github.com/heronalps/STOIC/server"
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
	runtime, _ := client.SelectRunTime(56, app, version, "")
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
	output, _, timeLog := client.RunOnNautilus(runtime, imageNum, app, version, transferTime)
	fmt.Printf("Output : %v..\n", string(output))
	assert.NotNil(t, timeLog)
}

func TestRunOnEdge(t *testing.T) {
	output, _, timeLog := client.RunOnEdge(imageNum, app, version)
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

func TestImageCache(t *testing.T) {
	num := server.ImageCache()
	fmt.Println(time.Now().UnixNano())
	fmt.Println(time.Now())
	fmt.Println(num)
}

func TestGenerateWorkLoad(t *testing.T) {
	workload := server.GenerateWorkLoad(6700)
	fmt.Println(workload)
}

func TestServerWorkload(t *testing.T) {
	fmt.Println(server.Workload[223])
}

func TestRegisterImages(t *testing.T) {
	path := "/Users/michaelzhang/Downloads/WTB_samples"
	//server.RegisterImages(path)
	decodeFile, err := os.Open(path + "/registryMap.gob")
	if err != nil {
		panic(err)
	}
	defer decodeFile.Close()

	decoder := gob.NewDecoder(decodeFile)
	registryMap := make(map[int32]string)

	decoder.Decode(&registryMap)
	//fmt.Println(registryMap)
}
