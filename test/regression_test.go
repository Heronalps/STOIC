package test

import (
	"fmt"
	"testing"

	"github.com/heronalps/STOIC/client"
)

func TestRegress2(t *testing.T) {
	client.Regress2()
}

func TestRegress(t *testing.T) {
	coef, intercept := client.Regress(runtime, app, version, numDP)
	fmt.Printf("coef : %v \n", coef)
	fmt.Printf("intercept : %v \n", intercept)
}
