package test

import (
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
