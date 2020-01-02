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
	runtime := client.SelectRunTime(56, app, version)
	fmt.Println("runtime : " + runtime)
}

func TestGetDeploymentTime(t *testing.T) {
	deploymentTime := client.GetDeploymentTime(runtime)
	fmt.Printf("runtime: %f \n", deploymentTime)
}

func TestCompareVersion(t *testing.T) {
	assert.Equal(t, client.CompareVersion("0", "1.0"), -1)
}

func TestRunOnNautilus(t *testing.T) {
	output, elapsed := client.RunOnNautilus(runtime, imageNum, app, version)
	fmt.Printf("Output : %v..\n", output)
	fmt.Printf("Elapsed : %v..\n", elapsed)
}
