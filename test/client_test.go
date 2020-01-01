package test

import (
	"fmt"
	"testing"

	"github.com/heronalps/STOIC/client"
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
