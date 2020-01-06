/*
Scheduler decides where to run WTB inferencing job given Pi bandwidth, number of images
*/

package client

import (
	"fmt"
	exec "os/exec"
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
		elapsed         float64
		transferTime    float64
		output          []byte
		selectedRuntime string
		predTimeLog     *TimeLog
		actTimeLog      *TimeLog
	)

	transferTime = GetTransferTime(imageNum)
	fmt.Printf("The bandwidth is %f megabits \n", GetBandWidth())
	fmt.Printf("The batch of %d images needs %f seconds to transfer\n", imageNum, transferTime)

	selectedRuntime, predTimeLog = SelectRunTime(imageNum, app, version, runtime)
	fmt.Println("Selected : " + selectedRuntime)
	fmt.Println("Runtime : " + runtime)
	// Update current runtime to accurately estimate deployment time
	currentRuntime = selectedRuntime

	output, elapsed, actTimeLog = Request(selectedRuntime, imageNum, app, version)
	if predTimeLog != nil {
		predTimeLog.Transfer = transferTime
	}
	if actTimeLog != nil {
		actTimeLog.Transfer = transferTime
	}

	if elapsed != 0.0 {
		AppendRecordProcessing(dbName, selectedRuntime, imageNum, elapsed, app, version)
		//For setup regressions, the prediction is based on preset coef & intercept
		LogTimes(predTimeLog, actTimeLog)
	}

	return output
}

/*
Request is a wrap function both for executing jobs and setting up processing time table for regression
*/
func Request(runtime string, imageNum int, app string, version string) ([]byte, float64, *TimeLog) {
	var (
		output     []byte
		elapsed    float64
		actTimeLog *TimeLog
	)
	switch runtime {
	case "edge":
		fmt.Println("Running on edge...")
		//output, elapsed, actTimeLog = RunOnEdge(imageNum, app, version)
	default:
		fmt.Println("Running on Nautilus...")
		output, elapsed, actTimeLog = RunOnNautilus(runtime, imageNum, app, version)
	}
	return output, elapsed, actTimeLog
}

/*
SelectRunTime select the runtime among four scenarios
*/
func SelectRunTime(imageNum int, app string, version string, runtime string) (string, *TimeLog) {

	// If the runtime is manually set, the results only have preset runtime
	totalTimes, predTimeLog := GetTotalTime(imageNum, app, version, runtime)
	fmt.Println(totalTimes)

	// Sort the totalTimes map by key
	keys := make([]float64, 0, len(totalTimes))
	for k := range totalTimes {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	selectedRuntime := totalTimes[keys[0]]
	fmt.Printf("The task is scheduled at %s for %f seconds\n", selectedRuntime, keys[0])
	return selectedRuntime, predTimeLog[selectedRuntime]
}

/*
RunOnEdge runs the task on mini edge cloud with AVX support
*/
func RunOnEdge(imageNum int, app string, version string) ([]byte, float64, *TimeLog) {
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
		return output, 0, nil
	}
	fmt.Printf("Output of task %s\n", string(output))
	procTime := ParseElapsed(output)

	return output, ParseElapsed(output), CreateTimeLog(0.0, 0.0, procTime)
}
