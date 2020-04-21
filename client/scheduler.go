/*
Scheduler decides where to run WTB inferencing job given Pi bandwidth, number of images
*/

package client

import (
	"errors"
	"fmt"
	"math"
	"os/exec"
	"regexp"
	"sort"

	retrygo "github.com/avast/retry-go"
)

/*
Schedule is the entry point of Scheduler.
Parameter: runtime is the preset runtime. If it's empty string, SelectRunTime will select one.
Return: runtime for appending processing time to corresponding table
*/
func Schedule(runtime string, imageNum int, zipPath string, app string, version string, all bool) []byte {
	var (
		output     []byte
		actTimeLog *TimeLog
		isDeployed bool
	)
	transferTimes := GetTransferTime(zipPath)

	// Redefine the selected runtime every selection
	selectedRuntimes := []string{}
	selectedRuntime, predTimeLog := SelectRunTime(imageNum, zipPath, app, version, runtime)
	fmt.Printf("The task is scheduled at %s \n", selectedRuntime)
	fmt.Printf("The bandwidth is %f megabits \n", GetBandWidth())
	fmt.Printf("The batch of %d images needs %f seconds to transfer to runtime %s\n",
		imageNum, transferTimes[selectedRuntime], selectedRuntime)

	if all {
		for runtime, isAvail := range runtimes {
			if isAvail {
				selectedRuntimes = append(selectedRuntimes, runtime)
			}
		}
	}
	selectedRuntimes = append(selectedRuntimes, selectedRuntime)

	for _, runtime := range selectedRuntimes {
		_, predTimeLog = SelectRunTime(imageNum, zipPath, app, version, runtime)

		retryErr := retrygo.Do(
			func() error {
				output, isDeployed, actTimeLog = Request(runtime, zipPath, imageNum, app, version, transferTimes[runtime])
				runtimes[runtime] = isDeployed
				if !isDeployed {
					return errors.New("request was not deployed")
				}
				return nil
			},
		)
		if retryErr != nil {
			fmt.Printf("Request failed: %v ...", retryErr.Error())
		}
		if actTimeLog != nil {
			actTimeLog.Transfer = transferTimes[runtime]
		}
		if actTimeLog != nil && actTimeLog.Processing != 0.0 {
			AppendRecordProcessing(dbName, runtime, imageNum, actTimeLog.Processing, app, version)
			//For setup regressions, the prediction is based on preset coef & intercept
			LogTimes(imageNum, app, version, runtime, predTimeLog, actTimeLog)
		}
	}

	return output
}

/*
Request is a wrap function both for executing jobs and setting up processing time table for regression
*/
func Request(runtime string, zipPath string, imageNum int, app string, version string, transferTime float64) ([]byte, bool, *TimeLog) {
	var (
		output     []byte
		isDeployed bool
		actTimeLog *TimeLog
	)
	switch runtime {
	case "edge":
		fmt.Println("Running on edge...")
		output, isDeployed, actTimeLog = RunOnEdge(zipPath, imageNum, app, version)
	default:
		fmt.Printf("Running on Nautilus...%s\n", runtime)
		output, isDeployed, actTimeLog = RunOnNautilus(runtime, imageNum, app, version, transferTime)
	}
	// The transfer time field is 0.0 in actTimeLog at this point
	return output, isDeployed, actTimeLog
}

/*
SelectRunTime select the runtime among four scenarios
*/
func SelectRunTime(imageNum int, zipPath string, app string, version string, runtime string) (string, *TimeLog) {
	var (
		selectedRuntime string = runtime
	)
	// If the runtime is manually set, the results only have preset runtime
	totalTimes, predTimeLogMap := GetTotalTime(zipPath, imageNum, app, version, runtime)
	fmt.Printf("totalTime: %v..\n", totalTimes)

	// Sort the totalTimes map by key
	keys := make([]float64, 0, len(totalTimes))
	for k := range totalTimes {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	if !math.IsNaN(keys[0]) {
		selectedRuntime = totalTimes[keys[0]]
	}

	// fmt.Printf("The task is scheduled at %s for %f seconds\n", selectedRuntime, keys[0])
	return selectedRuntime, predTimeLogMap[selectedRuntime]
}

/*
RunOnEdge runs the task on mini edge cloud with AVX support
*/
func RunOnEdge(zipPath string, imageNum int, app string, version string) ([]byte, bool, *TimeLog) {
	var (
		output     []byte
		err        error
		cmd        *exec.Cmd
		isDeployed bool
	)
	// repoPATH := HomeDir() + "/GPU_Serverless"

	// Run WTB image classification task
	FILE := "./apps/image_clf_inf-local.py "
	//cmdRun := "source venv/bin/activate && python " + FILE + strconv.Itoa(int(imageNum))
	cmdRun := "source venv/bin/activate && python " + FILE + zipPath
	cmd = exec.Command("bash", "-c", cmdRun)
	// cmd.Dir = repoPATH
	fmt.Printf("Start running task %s version %s on %d images \n", app, version, imageNum)
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("Error running task. msg: %s \n", err.Error())
		return output, isDeployed, nil
	}
	//fmt.Printf("Output of task %s\n", string(output))
	re := regexp.MustCompile(`Time with model.*`)
	lastline := re.Find(output)
	// fmt.Printf("lastline : %s..\n", lastline)
	return lastline, true, CreateTimeLog(0.0, 0.0, ParseElapsed(output))
}

/*
RunOnEdge runs the task on mini edge cloud with AVX support
*/
// func RunOnEdge(imageNum int, app string, version string) ([]byte, *TimeLog) {
// 	var (
// 		output []byte
// 		err    error
// 		cmd    *exec.Cmd
// 	)

// 	// Run WTB image classification task

// 	cmdRun := fmt.Sprintf("%s sh %s %d", minikubeConfig, invokeFile, imageNum)
// 	cmd = exec.Command("bash", "-c", cmdRun)
// 	fmt.Printf("Start running task %s version %s on %d images on Edge.. \n", app, version, imageNum)
// 	if output, err = cmd.Output(); err != nil {
// 		fmt.Printf("Error running task. msg: %s \n", err.Error())
// 		return output, nil
// 	}
// 	fmt.Printf("Output of task %s\n", string(output))

// 	return output, CreateTimeLog(0.0, 0.0, ParseElapsed(output))
// }
