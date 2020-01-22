/*
Scheduler decides where to run WTB inferencing job given Pi bandwidth, number of images
*/

package client

import (
	"fmt"
	"math"
	"os/exec"
	"sort"
	"strconv"
)

/*
Schedule is the entry point of Scheduler.
Parameter: runtime is the preset runtime. If it's empty string, SelectRunTime will select one.
Return: runtime for appending processing time to corresponding table
*/
func Schedule(runtime string, imageNum int, app string, version string) []byte {
	var (
		output          []byte
		selectedRuntime string
		predTimeLog     *TimeLog
		actTimeLog      *TimeLog
	)

	transferTimes := GetTransferTime(imageNum)

	selectedRuntime, predTimeLog = SelectRunTime(imageNum, app, version, runtime)
	fmt.Printf("The bandwidth is %f megabits \n", GetBandWidth())
	fmt.Printf("The batch of %d images needs %f seconds to transfer to runtime %s\n",
		imageNum, transferTimes[selectedRuntime], selectedRuntime)

	// Update current runtime to accurately estimate deployment time
	// if _, found := NautilusRuntimes[selectedRuntime]; found {
	// 	currentRuntime = selectedRuntime
	// }

	output, actTimeLog = Request(selectedRuntime, imageNum, app, version)
	if actTimeLog != nil {
		actTimeLog.Transfer = transferTimes[selectedRuntime]
	}
	// fmt.Printf("Selected Runtime: %s..\n", selectedRuntime)
	if actTimeLog != nil && actTimeLog.Processing != 0.0 {
		AppendRecordProcessing(dbName, selectedRuntime, imageNum, actTimeLog.Processing, app, version)
		//For setup regressions, the prediction is based on preset coef & intercept
		LogTimes(imageNum, app, version, selectedRuntime, predTimeLog, actTimeLog)
	}

	return output
}

/*
Request is a wrap function both for executing jobs and setting up processing time table for regression
*/
func Request(runtime string, imageNum int, app string, version string) ([]byte, *TimeLog) {
	var (
		output     []byte
		actTimeLog *TimeLog
	)
	switch runtime {
	case "edge":
		fmt.Println("Running on edge...")
		output, actTimeLog = RunOnEdge(imageNum, app, version)
	default:
		fmt.Println("Running on Nautilus...")
		output, actTimeLog = RunOnNautilus(runtime, imageNum, app, version)
	}
	// The transfer time field is 0.0 in actTimeLog at this point
	return output, actTimeLog
}

/*
SelectRunTime select the runtime among four scenarios
*/
func SelectRunTime(imageNum int, app string, version string, runtime string) (string, *TimeLog) {
	var (
		selectedRuntime string = runtime
	)
	// If the runtime is manually set, the results only have preset runtime
	totalTimes, predTimeLogMap := GetTotalTime(imageNum, app, version, runtime)
	fmt.Println(totalTimes)

	// Sort the totalTimes map by key
	keys := make([]float64, 0, len(totalTimes))
	for k := range totalTimes {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	if !math.IsNaN(keys[0]) {
		selectedRuntime = totalTimes[keys[0]]
	}

	fmt.Printf("The task is scheduled at %s for %f seconds\n", selectedRuntime, keys[0])
	return selectedRuntime, predTimeLogMap[selectedRuntime]
}

/*
RunOnEdge runs the task on mini edge cloud with AVX support
*/
func RunOnEdge(imageNum int, app string, version string) ([]byte, *TimeLog) {
	var (
		output []byte
		err    error
		cmd    *exec.Cmd
	)
	repoPATH := HomeDir() + "/GPU_Serverless"

	// Run WTB image classification task
	FILE := "./kubeless/image_clf/inference/local_version/image_clf_inf.py "
	cmdRun := "source venv/bin/activate && python " + FILE + strconv.Itoa(int(imageNum))
	cmd = exec.Command("bash", "-c", cmdRun)
	cmd.Dir = repoPATH
	fmt.Printf("Start running task %s version %s on %d images \n", app, version, imageNum)
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("Error running task. msg: %s \n", err.Error())
		return output, nil
	}
	fmt.Printf("Output of task %s\n", string(output))

	return output, CreateTimeLog(0.0, 0.0, ParseElapsed(output))
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
