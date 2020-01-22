/*
This module contains all helper functions.
*/

package client

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"os/user"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

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
	coef, intercept = Regress(runtime, app, version, procTimeNumDP)
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
func GetTransferTime(imageNum int) map[string]float64 {
	transferTimes := make(map[string]float64)
	// Convert megabits to megabytes
	bandwidth := GetBandWidth() / 8.0
	// Average JPG image size of 1920 * 1080 = 0.212 MB
	JPGSize := 212 * 1e-3

	transferTime := float64(imageNum) * JPGSize / bandwidth
	for i := 0; i < len(runtimes); i++ {
		if runtimes[i] == "edge" {
			transferTimes[runtimes[i]] = 0.0
		} else {
			transferTimes[runtimes[i]] = transferTime
		}
	}
	return transferTimes
}

/*
GetDeploymentTime returns the latest deployment time in the DeploymentTime table of specific runtime
*/
func GetDeploymentTime(runtime string) map[string]float64 {
	var (
		selectedRuntimes []string
	)
	deploymentTimes := make(map[string]float64)
	if runtime != "" {
		selectedRuntimes = []string{runtime}
	} else {
		selectedRuntimes = runtimes
	}

	for _, runtime := range selectedRuntimes {
		var deploymentTime float64
		if runtime == "edge" {
			deploymentTime = 0.0
		} else {
			deploymentTime = QueryDeploymentTime(runtime)
		}
		deploymentTimes[runtime] = deploymentTime
	}

	return deploymentTimes
}

/*
GetProcTime calculates the runtime of four scenarios: edge, cpu, gpu1, gpu2
*/
func GetProcTime(imageNum int, app string, version string, runtime string) map[string]float64 {
	var (
		selectedRuntimes []string
	)
	procTimes := make(map[string]float64)
	if runtime != "" {
		selectedRuntimes = []string{runtime}
	} else {
		selectedRuntimes = runtimes
	}
	for i := 0; i < len(selectedRuntimes); i++ {
		procTimes[selectedRuntimes[i]] = Extrapolate(selectedRuntimes[i], imageNum, app, version)
	}
	return procTimes
}

/*
GetTotalTime calculate total time (Addition of transfer and run time) of four scenarios
*/
func GetTotalTime(imageNum int, app string, version string, runtime string) (map[float64]string, map[string]*TimeLog) {
	var (
		selectedRuntimes []string
	)
	transferTimes := GetTransferTime(imageNum)
	procTimes := GetProcTime(imageNum, app, version, runtime)
	deploymentTimes := GetDeploymentTime(runtime)

	totalTimes := make(map[float64]string)
	timeLogs := make(map[string]*TimeLog)

	// When runtime is manually set, only iterate over the chosen one
	if runtime != "" {
		selectedRuntimes = []string{runtime}
	} else {
		selectedRuntimes = runtimes
	}
	for i := 0; i < len(selectedRuntimes); i++ {
		currRuntime := selectedRuntimes[i]
		totalTimes[transferTimes[currRuntime]+deploymentTimes[currRuntime]+procTimes[currRuntime]] = currRuntime
		timeLogs[currRuntime] = CreateTimeLog(transferTimes[currRuntime], deploymentTimes[currRuntime], procTimes[currRuntime])
	}
	return totalTimes, timeLogs
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

/*
Average returns the mean of float numbers in an array
*/
func Average(arr []float64) float64 {
	var total float64
	for _, value := range arr {
		total += value
	}
	return total / float64(len(arr))
}

/*
Median returns the median of an array of float numbers
*/
func Median(arr []float64) float64 {
	sort.Float64s(arr)
	mid := len(arr) / 2
	if len(arr)%2 == 0 {
		return (arr[mid-1] + arr[mid]) / 2
	}
	return arr[mid]
}

/*
GetWindowSize calculates the optimal window size for median of deployment time
*/
func GetWindowSize(runtime string) int {
	deploymentTimes := QueryDeploymentTimeSeries(runtime)

	// The minimum of maxWinSize and length of deploymentTimes is the length of MAE array
	minWin := 0
	minMAE := math.MaxFloat64
	for winSize := 1; winSize < Min(maxWinSize, len(deploymentTimes)); winSize++ {
		// Guard deploymentTimes length is less than maxWinSize
		length := len(deploymentTimes) - winSize
		totalAbsErr := 0.0
		for idx := 0; idx < length; idx++ {
			pred := Median(deploymentTimes[idx : idx+winSize])
			totalAbsErr += math.Abs(pred - deploymentTimes[idx+winSize])

		}
		currMAE := totalAbsErr / float64(length)
		if currMAE < minMAE {
			minMAE = currMAE
			minWin = winSize
		}
	}
	return minWin
}

/*
CompareVersion compares two version strings
Return 1 if first > second
Return -1 if first < second
Return 0 if first = second
*/
func CompareVersion(version1 string, version2 string) int {
	arr1 := strings.Split(version1, ".")
	arr2 := strings.Split(version2, ".")
	for i, j := 0, 0; i < len(arr1) || j < len(arr2); i, j = i+1, j+1 {
		var (
			curr1 int
			curr2 int
			err   error
		)
		if i < len(arr1) {
			curr1, err = strconv.Atoi(arr1[i])
		} else {
			curr1 = 0
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		if j < len(arr2) {
			curr2, err = strconv.Atoi(arr2[j])
		} else {
			curr2 = 0
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		if curr1 > curr2 {
			return 1
		} else if curr1 < curr2 {
			return -1
		}
	}
	return 0
}

/*
CreateTimeLog creates Time Log by three components
*/
func CreateTimeLog(transfer float64, deployment float64, processing float64) *TimeLog {
	var (
		timeLog *TimeLog
		total   float64
	)
	total = transfer + deployment + processing
	timeLog = &TimeLog{
		Total:      total,
		Transfer:   transfer,
		Deployment: deployment,
		Processing: processing,
	}
	return timeLog
}

/*
LogTimes logs Predicted/Actual total response time, transfer time, deployment time, processing time
*/
func LogTimes(imageNum int, app string, version string, runtime string, predTimeLog *TimeLog, actTimeLog *TimeLog) {
	if predTimeLog != nil && actTimeLog != nil {
		AppendRecordLogTime(imageNum, app, version, runtime,
			predTimeLog.Total, predTimeLog.Transfer, predTimeLog.Deployment, predTimeLog.Processing,
			actTimeLog.Total, actTimeLog.Transfer, actTimeLog.Deployment, actTimeLog.Processing)
	}
}

/*
UpdateWindowSizes updates optimal window sizes
*/
func UpdateWindowSizes() {
	ts1 := time.Now()
	for _, runtime := range NautilusRuntimes {
		windowSizes[runtime] = GetWindowSize(runtime)
	}
	fmt.Println(time.Now().Sub(ts1))
}

/*
Min function returns the minimum between two integers
*/
func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
