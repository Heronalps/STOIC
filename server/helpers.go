/*
This module contains all helper functions.
*/

package server

import (
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"

	"github.com/serverhorror/rog-go/reverse"
)

const (
	cpuDeploymentTime  = 18.0
	gpu1DeploymentTime = 46.0
	gpu2DeploymentTime = 65.0
)

/*
GetBandWidth solicits bandwidth of Pi Zero at Sedgwick Reserve.
*/
func GetBandWidth() float64 {
	url := "http://169.231.235.221/sedgtomayhem.txt"
	lines := "10"
	//fmt.Printf("Tailing http endpoint : %s\n", url)

	cmd := "curl " + url + " | tail -" + lines
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		println(err.Error())
		return 0
	}
	// Use reverse scanner to capture the last nonzero bandwidth
	scanner := reverse.NewScanner(strings.NewReader(string(out)))
	var (
		lastLine  string
		bandWidth float64
	)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fields := strings.Fields(lastLine)
		bandWidth, err = strconv.ParseFloat(fields[0], 64)
		if err != nil {
			fmt.Println(err.Error())
		}
		if bandWidth > 0.0 {
			break
		}
	}

	if bandWidth <= 0.0 {
		fmt.Println("The bandwidth zeroed out! Simulate on average 70.7916 !")
		bandWidth = 70.7916
	}
	return bandWidth
}

/*
HomeDir gets home directory of current user
*/
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		println(err.Error())
		return ""
	}
	return usr.HomeDir
}

/*
Extrapolate function inferences runtime by coefficient and intercept.
*/
func Extrapolate(mode string, x int) float64 {
	var coef float64
	var intercept float64
	switch mode {
	case "edge":
		coef = 2.39549861
		intercept = 13.600537473199736
	case "cpu":
		coef = 1.33380247
		intercept = 14.91093042617645
	case "gpu1":
		coef = 0.3271631
		intercept = 28.163551818338643
	case "gpu2":
		coef = 0.19928721
		intercept = 21.267003248222906
	}

	return float64(x)*coef + intercept
}

/*
GetTransferTime calculates the transfer time from Sedgwick reserve to Mayhem cloud to Nautilus
*/
func GetTransferTime(imageNum int) float64 {
	// Convert megabits to megabytes
	bandwidth := GetBandWidth() / 8.0
	// Average JPG image size of 1920 * 1080 = 0.212 MB
	JPGSize := 212 * 1e-3

	transferTime := float64(imageNum) * JPGSize / bandwidth
	return transferTime
}

/*
GetRunTime calculates the runtime of four scenarios: edge, cpu, gpu1, gpu2
*/
func GetRunTime(imageNum int) []float64 {
	edgeRuntime := Extrapolate("edge", imageNum)
	cpuRuntime := Extrapolate("cpu", imageNum)
	gpu1Runtime := Extrapolate("gpu1", imageNum)
	gpu2Runtime := Extrapolate("gpu2", imageNum)
	runtimes := []float64{edgeRuntime, cpuRuntime, gpu1Runtime, gpu2Runtime}
	return runtimes
}

/*
GetTotalTime calculate total time (Addition of transfer and run time) of four scenarios
*/
func GetTotalTime(imageNum int) map[float64]string {
	runtimes := GetRunTime(imageNum)
	totalTimes := make(map[float64]string)
	totalTimes[runtimes[0]] = "edge"
	totalTimes[runtimes[1]+GetAdditionTime("cpu", imageNum)] = "cpu"
	totalTimes[runtimes[2]+GetAdditionTime("gpu1", imageNum)] = "gpu1"
	totalTimes[runtimes[3]+GetAdditionTime("gpu2", imageNum)] = "gpu2"
	return totalTimes
}

/*
GetAdditionTime returns the sum of corresponding transfer and deployment time of runtime and image num
*/
func GetAdditionTime(runtime string, imageNum int) float64 {
	transferTime := GetTransferTime(imageNum)
	var additionTime float64
	switch runtime {
	case "edge":
		additionTime = 0.0
	case "cpu":
		additionTime = transferTime + cpuDeploymentTime
	case "gpu1":
		additionTime = transferTime + gpu1DeploymentTime
	case "gpu2":
		additionTime = transferTime + gpu2DeploymentTime
	}
	return additionTime
}

/*
parseElapsed capture time in output
*/
func parseElapsed(output []byte) float64 {
	re := regexp.MustCompile(`Time without model loading (\d*\.\d*)`)
	// []byte - elapsed time of task
	elapsed := re.FindSubmatch(output)
	if len(elapsed) == 0 {
		log.Println("No elapsed time is received in output ...")
		return 0.0
	}
	result, err := strconv.ParseFloat(string(elapsed[1]), 64)
	if err != nil {
		log.Println(err.Error())
		return 0.0
	}
	return result
}
