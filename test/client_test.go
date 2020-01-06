package test

import (
	"fmt"
	"testing"

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
	output, elapsed, _ := client.RunOnNautilus(runtime, imageNum, app, version)
	fmt.Printf("Output : %v..\n", output)
	fmt.Printf("Elapsed : %v..\n", elapsed)
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
	client.LogTimes(predTimeLog, actTimeLog)
}
