/*
This module contains all helper functions.
*/

package client

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
func Extrapolate(runtime string, imageNum int, app string, version string) float64 {
	var (
		coef      float64
		intercept float64
	)
	coef, intercept = Regress(runtime, app, version, numDP)
	if coef == 0.0 && intercept == 0.0 {
		switch runtime {
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
	}

	return float64(imageNum)*coef + intercept
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
GetProcTime calculates the runtime of four scenarios: edge, cpu, gpu1, gpu2
*/
func GetProcTime(imageNum int, app string, version string) map[string]float64 {
	procTimes := make(map[string]float64)
	for i := 0; i < len(runtimes); i++ {
		procTimes[runtimes[i]] = Extrapolate(runtimes[i], imageNum, app, version)
	}
	return procTimes
}

/*
GetTotalTime calculate total time (Addition of transfer and run time) of four scenarios
*/
func GetTotalTime(imageNum int, app string, version string) map[float64]string {
	procTimes := GetProcTime(imageNum, app, version)
	totalTimes := make(map[float64]string)
	for i := 0; i < len(runtimes); i++ {
		runtime := runtimes[i]
		totalTimes[procTimes[runtime]+GetAdditionTime(runtime, imageNum)] = runtime
	}
	return totalTimes
}

/*
GetDeploymentTime returns the latest deployment time in the DeploymentTime table of specific runtime
*/
func GetDeploymentTime(runtime string) float64 {
	if runtime == "edge" {
		return 0.0
	}
	return QueryDeploymentTime(runtime)
}

/*
GetAdditionTime returns the sum of corresponding transfer and deployment time of runtime and image num
*/
func GetAdditionTime(runtime string, imageNum int) float64 {
	if runtime == "edge" {
		return 0.0
	}
	transferTime := GetTransferTime(imageNum)
	return transferTime + GetDeploymentTime(runtime)
}

/*
ParseElapsed capture time in output
*/
func ParseElapsed(output []byte) float64 {
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
